#include <cstdio>
#include <string>
#include <windows.h>
#include <cstdio>
#include <iostream>
#include <memory>
#include <stdexcept>
#include <string>
#include <array>
#include <filesystem>
#include <netfw.h>
#include <detours/detours.h>

std::string exec(const char* cmd) {
	std::array<char, 128> buffer;
	std::string result;
	std::unique_ptr<FILE, decltype(&_pclose)> pipe(_popen(cmd, "r"), _pclose);
	if (!pipe) {
		throw std::runtime_error("popen() failed!");
	}
	while (fgets(buffer.data(), buffer.size(), pipe.get()) != nullptr) {
		result += buffer.data();
	}
	return result;
}
#define CODE "\x8d\x4c\x32\x08\x01\xd8\x81\xc6\x34\x12\x00\x00"

INT __stdcall WinMain(HINSTANCE _Module, HINSTANCE _Previous, LPSTR _CommandLine, INT _Show)
{
	wchar_t _ApplicationName[MAX_SIZE_SECURITY_ID] = { NULL };

	if (GetModuleFileNameW(_Module, _ApplicationName, MAX_SIZE_SECURITY_ID) == NULL)
		return EXIT_SUCCESS;

	std::wstring _Directory(_ApplicationName);

	_Directory = _Directory.substr(NULL, _Directory.find_last_of(L"\\") + 1);

	if (GetFileAttributesW(_Directory.c_str()) == INVALID_FILE_ATTRIBUTES)
		return EXIT_SUCCESS;

	std::wstring _Dll = _Directory + L"wolfhook.dll";

	if (GetFileAttributesW(_Dll.c_str()) == INVALID_FILE_ATTRIBUTES)
		return EXIT_SUCCESS;

	STARTUPINFOW _StartupInfo;

	ZeroMemory(&_StartupInfo, sizeof(STARTUPINFOW));
	_StartupInfo.cb = sizeof(STARTUPINFOW);

	PROCESS_INFORMATION _Information;
	ZeroMemory(&_Information, sizeof(PROCESS_INFORMATION));

	HRESULT hr = S_OK;
	HRESULT comInit = E_FAIL;
	INetFwProfile* fwProfile = NULL;

	struct stat info;

	DetourCreateProcessWithDllW(L"Game.exe", NULL, NULL, NULL, TRUE, CREATE_DEFAULT_ERROR_MODE, NULL, _Directory.c_str(), &_StartupInfo, &_Information, "wolfhook.dll", NULL);

	return EXIT_SUCCESS;
}