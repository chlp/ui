package logger

import (
	"fmt"
	"log"
	"os"
)

var logger *log.Logger

func InitLogger(logFile string) {
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

func Printf(format string, args ...interface{}) {
	fmt.Printf(format, args...)
	if logger != nil {
		logger.Printf(format, args...)
	}
}

func Fatalf(format string, args ...interface{}) {
	Printf("Fatal: "+format, args...)
	if logger != nil {
		logger.Fatalf(format, args...)
	}
	os.Exit(1)
}
