package logging

import (
	"fmt"
	"log"
	"os"
	"strings"
)

const (
	DebugLevel = "DEBUG"
	InfoLevel  = "INFO"
)

var logLevel string

func init() {
	ll := strings.ToLower(os.Getenv("LOG_LEVEL"))
	switch ll {
	case "debug":
		logLevel = DebugLevel
	case "info":
		logLevel = InfoLevel
	default:
		fmt.Println("LOG_LEVEL not set - defaulting to INFO")
		logLevel = InfoLevel
	}
}

// Debugf logs a statement using fmt printing semantics if the log
// level is set to debug.
func Debugf(format string, v interface{}) {
	if logLevel == DebugLevel {
		msg := fmt.Sprintf(format, v)
		printLog(DebugLevel, msg)
	}
}

// Debugln logs a single line if log level is set to debug.
func Debugln(msg string) {
	if logLevel == DebugLevel {
		printLog(DebugLevel, msg)
	}
}

// Infof logs a statement using fmt printing semantics no matter the log level.
func Infof(format string, v interface{}) {
	msg := fmt.Sprintf(format, v)
	printLog(InfoLevel, msg)
}

// Infoln logs a single line no matter the log level.
func Infoln(msg string) {
	printLog(InfoLevel, msg)
}

func printLog(level string, message string) {
	log.Printf("%s: %s", level, message)
}
