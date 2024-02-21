#include <iostream>
#include <algorithm>
#include <parser-library/parse.h>
#include <capstone/capstone.h>
#include <set>
#include <fstream>
#include <iterator>
#include "pool.h"

using namespace peparse;


int min_dis = 2147483640;
int found_addr = 0;
std::mutex mutex;

struct userdata {
	VA entryPoint;
	uint8_t *code_buf = nullptr;
	uint32_t code_len;
	VA base;
};

constexpr int line_cost = 4;

size_t diff_size(size_t a, size_t b) {
	if (a > b)
		return a - b;
	else return b - a;
}


int assembly_edit_distance(const std::vector<uint8_t>& src, const std::set<uint32_t>& src_ends, const std::vector<uint8_t>& dest, const std::set<uint32_t>& dest_ends) {
	if ( diff_size(src.size(), dest.size()) > dest.size()) {
		return diff_size(src.size() ,dest.size());
	}
	std::vector<int> current_row(dest.size() + 1, 0);
	std::vector<int> first_column(src.size() + 1, 0);

	for (int i = 0; i < dest.size() + 1; ++i) {
		current_row[i] = i;
	}

	for (int i = 0; i < src.size() + 1; ++i) {
		first_column[i] = i;
	}

	for (int i = 1; i < src.size() + 1; ++i) {
		std::vector<int> next_row(dest.size() + 1, 0);
		int last_col = first_column[i];
		for (int j = 1; j < dest.size() + 1; ++j) {
			next_row[j] = std::min(current_row[j] +1, last_col + 1);
			if (src[i - 1] == dest[j - 1]) {
				next_row[j] = std::min(next_row[j], current_row[j-1]);
			}
			else {
				next_row[j] = std::min(next_row[j], current_row[j - 1] + 1);
			}
			last_col = next_row[j];
		}
		current_row = next_row;
	}
	return current_row[dest.size()];
}

void worker(int func, const std::vector<uint8_t>& src, const std::vector<uint8_t>& dest) {
	auto dis = assembly_edit_distance(src, std::set<uint32_t>(), dest, std::set<uint32_t>());
	std::cout << "looked:" << std::hex << func << std::endl;
	auto lock = std::unique_lock<std::mutex>(mutex);
	if (min_dis > dis) {
		min_dis = dis;
		found_addr = func;
	}
}

bool is_number(std::string str) {
	for (int i = 0; i < str.size(); ++i) {
		if (!((str[i] >= '0' && str[i] <= '9') || str[i] == ' ' || (str[i] >= 'a' && str[i] <= 'f')))
			return false;
	}
	return true;
}

int main()
{
	auto pe = ParsePEFromFile("Game.exe");
	VA entryPoint;
	GetEntryPoint(pe, entryPoint);
	userdata ud = {
		entryPoint,
		nullptr,
		0
	};
	IterSec(pe, [](void *N,
		const VA &secBase,
		const std::string &secName,
		const image_section_header &s,
		const bounded_buffer *data) {
		auto state = (userdata*)N;
		if (secBase <= state->entryPoint && state->entryPoint < secBase + data->bufLen) {
			auto buf = new uint8_t[data->bufLen];
			memcpy(buf, data->buf, data->bufLen);
			state->code_buf = buf;
			state->code_len = data->bufLen;
			state->base = secBase;
		}
		return 0;
	}, &ud);

	if (ud.code_buf == nullptr) {
		std::cout << "Couldn't find appropriate text segment" << std::endl;
		return 1;
	}

	std::vector<std::vector<uint8_t> > code;
	std::map<uint32_t, int> addr_lookup;
	std::set<uint32_t> funcs;

	csh handle;
	cs_insn *insn;

	auto error = cs_open(CS_ARCH_X86, CS_MODE_32, &handle);
	if (error != CS_ERR_OK) {
		std::cout << "WTF" << std::endl;
		return 1;
	}
	cs_option(handle, CS_OPT_SKIPDATA, CS_OPT_ON);
	std::cout << "opened code size:" << ud.code_len << std::endl;
	auto count = cs_disasm(handle, ud.code_buf, ud.code_len, ud.base, 0, &insn);
	if (count > 0) {
		for (int i = 0; i < count; ++i) {
			addr_lookup.emplace(insn[i].address, code.size());

			if (strcmp(insn[i].mnemonic, "ret") == 0)
			{
				code.push_back(std::vector<uint8_t>());
			}
			else {
				code.push_back(std::vector<uint8_t>(insn[i].bytes, insn[i].bytes + insn[i].size));
				
			}

			if (strcmp(insn[i].mnemonic, "call") == 0) {
				auto str = std::string(&insn[i].op_str[2]);
				if (is_number(str)) {
					funcs.insert(std::strtol(str.c_str(), nullptr, 16));
				}
			}
		}
	}
	cs_close(&handle);

	std::ifstream f("drawstring.dat", std::ios::binary | std::ios::ate);
	int size = f.tellg();
	f.seekg(0, std::ios::beg);
	std::vector<uint8_t> draw_string_code;
	draw_string_code.resize(size);
	f.read((char*)draw_string_code.data(), size);

	error = cs_open(CS_ARCH_X86, CS_MODE_32, &handle);
	if (error != CS_ERR_OK) {
		std::cout << "WTF" << std::endl;
		return 1;
	}
	cs_option(handle, CS_OPT_SKIPDATA, CS_OPT_ON);

	std::vector<uint8_t> draw_string;
	std::set<uint32_t> draw_string_ends;
	count = cs_disasm(handle, draw_string_code.data(), size, 0, 0, &insn);
	if (count > 0) {
		for (int i = 0; i < count; ++i) {
			draw_string.insert(draw_string.end(), insn[i].bytes, insn[i].bytes + insn[i].size);
			draw_string_ends.insert(insn[i].address + insn[i].size);
		}
	}
	cs_close(&handle);
	
	std::map<uint32_t, uint32_t> size_map;
	
	{
		ThreadPool pool(4);
		for (auto func : funcs) {
			if (addr_lookup.find(func) == addr_lookup.end()) {
				func = func - 1;
			}
			if (addr_lookup.find(func) == addr_lookup.end()) {
				func = func - 1;
			}
			if (addr_lookup.find(func) == addr_lookup.end()) {
				func = func + 3;
			}
			if (addr_lookup.find(func) == addr_lookup.end()) {
				func = func + 1;
			}
			if (addr_lookup.find(func) != addr_lookup.end()) {
				int i = addr_lookup.at(func);
				std::vector<uint8_t> payload;
				std::set<uint32_t> ends;
				int size = 0;
				while (true) {
					if (code[i].size() == 0) {
						break;
					}
					payload.insert(payload.end(), code[i].begin(), code[i].end());
					ends.insert(payload.size());
					++i;
					++size;
				}
				pool.enqueue(worker, func, payload, draw_string);
				size_map.emplace(func, size + 1);
			}
		}
	}
	
	error = cs_open(CS_ARCH_X86, CS_MODE_32, &handle);
	if (error != CS_ERR_OK) {
		std::cout << "WTF" << std::endl;
		return 1;
	}
	cs_option(handle, CS_OPT_SKIPDATA, CS_OPT_ON);
	
	auto draw_string_size = size_map.at(found_addr);

	if (ud.code_buf[found_addr - ud.base] != 0x55) {
		found_addr -= 1;
	}
	if (ud.code_buf[found_addr - ud.base] != 0x55) {
		found_addr -= 1;
	}
	if (ud.code_buf[found_addr - ud.base] != 0x55) {
		found_addr += 3;
	}
	if (ud.code_buf[found_addr - ud.base] != 0x55) {
		found_addr += 1;
	}
	if (ud.code_buf[found_addr - ud.base] != 0x55) {
		std::cout << "couldn't find push ebp" << std::endl;
		return 1;
	}

	count = cs_disasm(handle, &ud.code_buf[found_addr - ud.base], ud.code_len - found_addr, found_addr, draw_string_size + 2, &insn);
	
	int get_cstr_addr = 0;
	int copy_buf_addr = 0;

	if (count > 0) {
		for (int i = 0; i < count; ++i) {
			if (strcmp(insn[i].mnemonic, "call") == 0)
			{
				if (!get_cstr_addr) {
					if (strcmp(insn[i - 1].mnemonic, "mov") == 0 && strcmp(insn[i - 2].mnemonic, "push") == 0) {
						auto mov_detail = std::string(insn[i - 1].op_str);
						auto push_detail = std::string(insn[i - 2].op_str);
						if (mov_detail.rfind("ecx") == 0 && push_detail == "0") {
							auto str = std::string(&insn[i].op_str[2]);
							get_cstr_addr = std::strtol(str.c_str(), nullptr, 16);
						}
					}
				}
				if (!copy_buf_addr) {
					if (strcmp(insn[i - 1].mnemonic, "lea") == 0 && strcmp(insn[i - 2].mnemonic, "push") == 0) {
						auto lea_detail = std::string(insn[i - 1].op_str);
						auto push_detail = std::string(insn[i - 2].op_str);
						bool is_push_slash = false;
						if (push_detail.size() > 2) {
 							if (is_number(push_detail.substr(2))) {
								auto slash_addr = std::strtol(push_detail.substr(2).c_str(), nullptr, 16);
								uint8_t slash_byte;
								if (ReadByteAtVA(pe, slash_addr, slash_byte)) {
									is_push_slash = slash_byte == '\\';
								}
							}
							if (lea_detail.rfind("ecx") == 0 && is_push_slash) {
								auto str = std::string(&insn[i].op_str[2]);
								copy_buf_addr = std::strtol(str.c_str(), nullptr, 16);
							}
						}
					}
				}
			}

			if (strcmp(insn[i].mnemonic, "ret") == 0)
			{
				break;
			}
		}
	}

	if (!get_cstr_addr || !copy_buf_addr) {
		std::cout << "could not find standard library string functions" << std::endl;
		return 1;
	}

	std::cout << "successfully found all required functions!" << std::endl;

	std::cout << "draw_string:" << std::hex << found_addr << std::endl;
	std::cout << "get_cstr:" << std::hex << get_cstr_addr << std::endl;
	std::cout << "copy_buf:" << std::hex << copy_buf_addr << std::endl;

	std::ofstream fout("wolfpatchinfo", std::ios::out);
	fout << found_addr << std::endl;
	fout << get_cstr_addr << std::endl;
	fout << copy_buf_addr << std::endl;
	fout.close();

	/*
	int i = addr_lookup.at(0x00494480);
	std::cout << "func:" << i << std::endl;
	std::vector<uint8_t> payload;
	std::set<uint32_t> ends;
	while (true) {
		if (code[i].size() == 0) {
			break;
		}
		payload.insert(payload.end(), code[i].begin(), code[i].end());
		ends.insert(payload.size());
		++i;
	}
	std::ofstream fout("drawstring.dat", std::ios::out | std::ios::binary);
	fout.write((char*)payload.data(), payload.size());
	fout.close();
	*/
	//std::cout << "payload:" << i << std::endl;
	//int dis = assembly_edit_distance(payload, ends, draw_string_code);
	//std::cout << dis << std::endl;

	return 0;
}