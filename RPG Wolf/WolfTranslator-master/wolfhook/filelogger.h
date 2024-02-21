#pragma once
#include <iostream>
#include <fstream>

class FileLogger {

public:
	explicit FileLogger(const char *engine_version, const char *fname = "ige_log.txt")
	{
		myFile.open(fname, std::ofstream::out | std::ofstream::app);

		myFile << "Log file created" << std::endl << std::endl;

	}

	~FileLogger() {

		myFile << std::endl << std::endl;

		myFile.close();

	}

	friend FileLogger &operator << (FileLogger &logger, const char *text) {
		logger.myFile << text;
		return logger;

	}

	FileLogger(const FileLogger &) = delete;
	FileLogger &operator= (const FileLogger &) = delete;
	std::ofstream           myFile;
};