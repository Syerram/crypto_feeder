package main

import (
	"io"
	"io/ioutil"
	"log"
	"os"
)

const (
	TRACE = 1 << iota
	WARNING
	INFO
	ERROR
)

type Logger struct {
	Trace   *log.Logger
	Info    *log.Logger
	Warning *log.Logger
	Error   *log.Logger
	level   int
}

func NewLogger(logLevel int, baseLogPath string) *Logger {
	//setup loggers with files
	logger := Logger{}
	//TODO: externalize log file path
	if len(baseLogPath) == 0 {
		baseLogPath = "/tmp/"
	}
	trace_w, warning_w, info_w, error_w := GetLogFile(baseLogPath+"ct_trace.log"),
		GetLogFile(baseLogPath+"ct_warning.log"),
		GetLogFile(baseLogPath+"ct_info.log"),
		GetLogFile(baseLogPath+"ct_error.log")

	switch logLevel {
	case WARNING:
		trace_w = ioutil.Discard
	case INFO:
		trace_w = ioutil.Discard
		warning_w = ioutil.Discard
	case ERROR:
		trace_w = ioutil.Discard
		warning_w = ioutil.Discard
		info_w = ioutil.Discard
	}

	logger.level = logLevel
	logger.Trace = log.New(trace_w, "TRACE: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Info = log.New(info_w, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Warning = log.New(warning_w, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile)
	logger.Error = log.New(error_w, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile)
	return &logger
}

func GetLogFile(path string) io.Writer {
	file, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file:", err)
		os.Exit(1)
	}
	return file
}
