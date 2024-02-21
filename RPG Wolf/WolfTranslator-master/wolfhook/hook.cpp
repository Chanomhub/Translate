#define _CRT_SECURE_NO_WARNINGS
#define WIN32_LEAN_AND_MEAN
#include <windows.h>
#include "httplib.h"

#include <iostream>
#include <functional>
#include <assert.h>
#include <detours/detours.h>
#include <unicode/unistr.h>
#include <cstdio>
#include <memory>
#include <stdexcept>
#include <string>
#include <array>
#include <sstream>
#include <list>
#include <stack>
#include <set>
#include <map>

#include "filelogger.h"
#include "pool.h"
#include "hexdump.h"
#include <capstone/capstone.h>

std::map<std::string, std::string> translated_strs;
std::set< std::string> complete_translated_strs;
std::thread* translation_thread;

std::mutex thread_mutex;
std::list<std::string> translation_queue;


ThreadPool* pool;
FileLogger* fileLogger;

constexpr int BUFSIZE = 128;
constexpr int POOL_SIZE = 1;

HANDLE hJob;

void* DrawString_target;

typedef char* (__fastcall *get_cstr_func)(int thi, int notUsed, int a);
get_cstr_func get_data;

typedef char* (__fastcall *copy_cstr_func)(int thi, int notUsed, const char* a);
copy_cstr_func copy_cstr;

void Log(const char* format, ...) {
	char linebuf[1024];
	va_list v;
	va_start(v, format);
	vsprintf(linebuf, format, v);
	va_end(v);
	*fileLogger << linebuf;
}

// Disable inline functoin optimization
// otherwise it might not work
// seperated from transverse_stack to fight stack check
// that corrupts ebp
int get_ebp() {
	int stack_base;
	__asm {
		mov stack_base, ebp
	}
	return stack_base;
}

bool DetourTransaction(std::function<bool()> callback) {
	LONG status = DetourTransactionBegin();
	if (status != NO_ERROR) {
		Log("DetourTransactionBegin failed with %08x\n", status);
		return status;
	}

	if (callback()) {
		status = DetourTransactionCommit();
		if (status != NO_ERROR) {
			Log("DetourTransactionCommit failed with %08x\n", status);
		}
	}
	else {
		status = DetourTransactionAbort();
		if (status == NO_ERROR) {
			Log("Aborted transaction.\n");
		}
		else {
			Log("DetourTransactionAbort failed with %08x\n", status);
		}
	}
	return status == NO_ERROR;
}

std::string sjisToUtf8(const std::string& value)
{
	try {
		icu::UnicodeString src(value.c_str(), "shift_jis");
		int length = src.extract(0, src.length(), NULL, "utf8");
		Log("Legnth: %d\n", length);
		std::vector<char> result(length + 1);
		src.extract(0, src.length(), &result[0], "utf8");
		return std::string(result.begin(), result.end() - 1);
	}
	catch (std::exception e) {
		return "";
	}

}

void create_process(LPWSTR cmd) {
	SECURITY_ATTRIBUTES saAttr;

	hJob = CreateJobObject(NULL, NULL);

	JOBOBJECT_EXTENDED_LIMIT_INFORMATION jeli = {
		0
	};
	jeli.BasicLimitInformation.LimitFlags = JOB_OBJECT_LIMIT_KILL_ON_JOB_CLOSE;
	SetInformationJobObject(hJob, JobObjectExtendedLimitInformation, &jeli, sizeof(jeli));

	AssignProcessToJobObject(hJob, GetCurrentProcess());

	HANDLE g_hChildStd_IN_Rd = NULL;
	HANDLE g_hChildStd_IN_Wr = NULL;
	HANDLE g_hChildStd_OUT_Rd = NULL;
	HANDLE g_hChildStd_OUT_Wr = NULL;

	saAttr.nLength = sizeof(SECURITY_ATTRIBUTES);
	saAttr.bInheritHandle = TRUE;
	saAttr.lpSecurityDescriptor = NULL;

	if (!CreatePipe(&g_hChildStd_OUT_Rd, &g_hChildStd_OUT_Wr, &saAttr, 0))
		assert(false);


	if (!SetHandleInformation(g_hChildStd_OUT_Rd, HANDLE_FLAG_INHERIT, 0))
		assert(false);


	if (!CreatePipe(&g_hChildStd_IN_Rd, &g_hChildStd_IN_Wr, &saAttr, 0))
		assert(false);


	if (!SetHandleInformation(g_hChildStd_IN_Wr, HANDLE_FLAG_INHERIT, 0))
		assert(false);


	PROCESS_INFORMATION piProcInfo;
	STARTUPINFO siStartInfo;
	BOOL bSuccess = FALSE;


	ZeroMemory(&piProcInfo, sizeof(PROCESS_INFORMATION));

	ZeroMemory(&siStartInfo, sizeof(STARTUPINFO));
	siStartInfo.cb = sizeof(STARTUPINFO);

	bSuccess = CreateProcess(NULL,
		cmd,
		NULL,
		NULL,
		TRUE,
		CREATE_NO_WINDOW,
		NULL,
		NULL,
		&siStartInfo,
		&piProcInfo);

	if (!bSuccess)
		assert(false);

	AssignProcessToJobObject(hJob, piProcInfo.hProcess);
	CloseHandle(piProcInfo.hProcess);
	CloseHandle(piProcInfo.hThread);

	CloseHandle(g_hChildStd_OUT_Wr);
	CloseHandle(g_hChildStd_IN_Rd);
}

std::string translate_impl(std::string str) {
	httplib::Client cli("localhost", 5731);
	auto resp = cli.Post("/", str, "text/plain");
	cli.set_read_timeout(600, 0);

	if (resp) {
		return resp->body;
	}
	
	return "";
}

std::string translate_string(std::string str) {
	if (str.size() <= 1) {
		return str;
	}
	std::vector<bool> types;
	std::vector<std::string> cutted;
	bool exit = false;
	int lasti = 0;
	int i = 0;
	bool state = str[0] == 0x1;
	for (; i < str.size(); ++i) {
		if (str[i] == 0x1) {
			if (!state) {
				cutted.push_back(str.substr(lasti, i-lasti));
				lasti = i;
				types.push_back(false);
				state = true;
			}
			auto k = str[i + 2] - 1;
			Log("K: %d\n", k);
			switch (k) {
			case 0:
				break;
			case 1:
				i += 4;
				break;
			case 2:
				i += 3;
				break;
			case 3:
				i += 3;
				break;
			case 4:
			case 5:
				break;
			case 6:
				i++;
				break;
			case 8:
				// TODO more investigation
				break;
			case 10:
				break;
			case 11:
				i += 4;
				break;
			case 13:
			case 14:
			case 15:
			case 16:
				i += 2;
				break;
			case 17:
				i += 2;
				break;
			case 18:
			case 19:
			case 20:
				break;
			case 21:
			case 22:
			case 30:
			case 33:
				Log("TEMPLATE\n");
				break;
			case 23:
				i += 3;
				break;
			case 24:
				i += 3;
				break;
			case 25:
			case 26:
				break;
			case 27:
			case 28:
				i += 3;
				break;
			case 29:
				i += 3;
				break;
			case 31:
				break;
			case 32:
				i += 3;
				break;
			case 34:
				i += 3;
				break;
			case 35:
				i += 3;
				break;
			case 36:
				i += 3;
				break;
			default:
				Log("Invalid\n");
			}
			i += 2;
		}
		else {
			if (state) {
				cutted.push_back(str.substr(lasti, i - lasti));
				lasti = i;
				types.push_back(true);
				state = false;
			}
			if ((str[i] >= 0x81 && str[i] <= 0x9F) || (str[i] >= 0xE0 && str[i] <= 0xFB)) {
				++i;
			}
		}
	}
	if (lasti < str.size()) {
		cutted.push_back(str.substr(lasti, str.size() - lasti));
		types.push_back(state);
	}
	
	std::string out;
	for (int i = 0; i < cutted.size(); i++) {
		std::string toadd;
		if (!types[i]) {
			Log("Cutted hex: \n");
			auto decoded = sjisToUtf8(cutted[i]);
			hex_dump((void*)cutted[i].c_str(), min(cutted[i].size(), 400), fileLogger->myFile);
			toadd = translate_impl(decoded);
			if (toadd == "") {
				return "";
			}
		}
		else
			toadd = cutted[i];
		out += toadd;
	}

	return out;
}

void translation_handler()
{
	std::set<std::string> queued;
	while (true) {
		std::string item;
		{
			auto lock = std::unique_lock<std::mutex>(thread_mutex);
			if (translation_queue.empty()) {
				lock.unlock();
				auto x = std::chrono::milliseconds(10);
				std::this_thread::sleep_for(x);
				continue;
			}
			item = translation_queue.front();
			translation_queue.pop_front();
			if (queued.find(item) != queued.end())
				continue;
			queued.insert(item);
		}
	
		Log("Original hex: \n");
		hex_dump((void*)item.c_str(), min(item.size(), 800), fileLogger->myFile);

		pool->enqueue([&] (std::string source) {
			auto translated = translate_string(source);
			auto lock = std::unique_lock<std::mutex>(thread_mutex);
			translated_strs.emplace(source, translated);
			complete_translated_strs.insert(translated);
			queued.erase(source);
		}, item);
	}
}

void request_translation(std::string str) {
	translation_queue.remove(str);
	translation_queue.push_front(str);
}

void(__fastcall *DrawString_trampoline1) (int thi, void* notUsed, char a2, double a3, double a4);
void(__fastcall *DrawString_trampoline) (int thi, void* notUsed, char a2, double a3, double a4, double a5);

void DrawString_impl(int thi) {
	if (!pool || !translation_thread) {
		wchar_t cmd[] = L".\\translator\\translator.exe";
		create_process(cmd);
		pool = new ThreadPool(POOL_SIZE);
		translation_thread = new std::thread(translation_handler);
	}
	auto x = get_data(thi, thi, 0);
	if (x && *x) {
		auto lock = std::unique_lock<std::mutex>(thread_mutex);
		if (complete_translated_strs.find(x) == complete_translated_strs.end() && translated_strs.find(x) == translated_strs.end()) {
			request_translation(x);
		}
		else if (translated_strs.find(x) != translated_strs.end()) {
			copy_cstr(thi, thi, translated_strs[x].c_str());
		}
	}
}

void __fastcall DrawString1(int thi, void* notUsed, char a2, double a3, double a4)
{
	DrawString_impl(thi);
	DrawString_trampoline1(thi,notUsed, a2, a3, a4);
}

void __fastcall DrawString(int thi, void* notUsed, char a2, double a3, double a4, double a5)
{
	DrawString_impl(thi);
	DrawString_trampoline(thi, notUsed, a2, a3, a4, a5);
}

void load_wolfpatchinfo() {
	std::ifstream f("wolfpatchinfo", std::ios::in);
	int draw_string_addr, get_cstr_addr, copy_buf_addr;
	f >> draw_string_addr;
	f>> get_cstr_addr;
	f >> copy_buf_addr;
	DrawString_target = reinterpret_cast<void*>(draw_string_addr);
	get_data = reinterpret_cast<get_cstr_func>(get_cstr_addr);
	copy_cstr = reinterpret_cast<copy_cstr_func>(copy_buf_addr);
}

bool is_file_exist(const char *fileName)
{
	std::ifstream infile(fileName);
	return infile.good();
}

BOOL WINAPI DllMain(
	HINSTANCE hinstDLL,
	DWORD fdwReason,
	LPVOID lpReserved)
{
	switch (fdwReason)
	{
	case DLL_PROCESS_ATTACH:
		fileLogger = new FileLogger("1.0", "log.txt");
		Log("Attaching");
		load_wolfpatchinfo();
		DetourTransaction([&]() {
			void *target = nullptr,
				*detour = nullptr;
			auto is_ver1 = is_file_exist("wolfdrawstring1");
			if (is_ver1) {
				DetourAttachEx(&DrawString_target,
					DrawString1,
					reinterpret_cast<PDETOUR_TRAMPOLINE*>(&DrawString_trampoline1),
					&target,
					&detour);
			}
			else {
				DetourAttachEx(&DrawString_target,
					DrawString,
					reinterpret_cast<PDETOUR_TRAMPOLINE*>(&DrawString_trampoline),
					&target,
					&detour);
			}


			return true;
		});

		break;
	case DLL_THREAD_ATTACH:
		break;

	case DLL_THREAD_DETACH:
		break;

	case DLL_PROCESS_DETACH:
		CloseHandle(hJob);

		break;
	}
	return TRUE;
}

extern "C" __declspec(dllexport)VOID NullExport(VOID)
{
}
