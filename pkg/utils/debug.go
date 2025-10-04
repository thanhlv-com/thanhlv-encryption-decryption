package utils

import (
	"fmt"
	"log"
	"os"
	"time"
)

var debugLogger *log.Logger
var isDebugEnabledFunc func() bool

func init() {
	debugLogger = log.New(os.Stderr, "[DEBUG] ", log.Ldate|log.Ltime|log.Lshortfile)
}

func SetDebugEnabledFunc(fn func() bool) {
	isDebugEnabledFunc = fn
}

func isDebugEnabled() bool {
	if isDebugEnabledFunc != nil {
		return isDebugEnabledFunc()
	}
	return false
}

func DebugLog(format string, args ...interface{}) {
	if !isDebugEnabled() {
		return
	}
	message := fmt.Sprintf(format, args...)
	debugLogger.Printf("%s", message)
}

func DebugLogf(format string, args ...interface{}) {
	if !isDebugEnabled() {
		return
	}
	debugLogger.Printf(format, args...)
}

func DebugLogWithTimestamp(message string) {
	if !isDebugEnabled() {
		return
	}
	timestamp := time.Now().Format("2006-01-02 15:04:05")
	debugLogger.Printf("[%s] %s", timestamp, message)
}