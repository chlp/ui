package logger

import (
	"log"
	"os"
)

var logger *log.Logger

func InitLogger(logFile string) {
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	logger = log.New(file, "", log.LstdFlags)
}

func Printf(format string, args ...interface{}) {
	logger.Printf(format, args...)
}

func Fatalf(format string, args ...interface{}) {
	logger.Fatalf(format, args...)
}
