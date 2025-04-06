package logger

import (
	"fmt"
	"log"
	"os"
)

var (
	appName   string
	logger    *log.Logger
	withDebug bool
)

func InitLogger(name, logFile string, debug bool) {
	appName = name
	withDebug = debug

	if logFile == "" {
		// will use only stdout
		return
	}
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logger = log.New(file, "", log.LstdFlags)
}

func Debugf(format string, args ...interface{}) {
	if !withDebug {
		return
	}
	format = "[debug] " + format
	if appName != "" {
		format = appName + " | " + format
	}
	fmt.Printf(format+"\n", args...)
	if logger != nil {
		logger.Printf(format, args...)
	}
}

func Printf(format string, args ...interface{}) {
	if appName != "" {
		format = appName + " | " + format
	}
	fmt.Printf(format+"\n", args...)
	if logger != nil {
		logger.Printf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	format = "[fatal] " + format
	if appName != "" {
		format = appName + " | " + format
	}
	Printf(format, args...)
	if logger != nil {
		logger.Fatalf(format, args...)
	}
	os.Exit(1)
}
