package logs

import (
	"fmt"
	"strings"
	"time"
)

const (
	LogLevelDebug = "DEBUG"
	LogLevelInfo  = "INFO"
	LogLevelWarn  = "WARN"
	LogLevelErr   = "ERROR"
	timeLayout    = "02.01.2006 15:04:05"
)

func LogDebug(message string, args ...interface{}) {
	log(LogLevelDebug, message, args...)
}

func LogInfo(message string, args ...interface{}) {
	log(LogLevelInfo, message, args...)
}

func LogWarn(message string, args ...interface{}) {
	log(LogLevelWarn, message, args...)
}

func LogErrCustom(message string, args ...interface{}) {
	log(LogLevelErr, message, args...)
}

func LogErr(message string, e error) {
	log(LogLevelErr, fmt.Sprintf("%s: %s", message, e.Error()))
}

func log(level string, message string, args ...interface{}) {
	t := time.Now()

	line := t.Format(timeLayout) + " " + level + ": "

	if !strings.HasSuffix(message, "\n") && !strings.HasSuffix(message, "\r") {
		message += "\n"
	}
	fmt.Printf(line+message, args...)
}
