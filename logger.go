package otplugins

import (
	"fmt"
	"log"
	"os"
)

var globalLogger Logger = log.New(os.Stderr, "", log.Flags())
var globalLoggerLevel = LogLevelDebug

func SetGlobalLogger(logger Logger, level LogLevel) {
	globalLogger = logger
	globalLoggerLevel = level
}

func Log(level LogLevel, args ...interface{}) error {
	if level < globalLoggerLevel {
		return nil
	}

	return globalLogger.Output(2, fmt.Sprint(args...))
}

func Logf(level LogLevel, format string, args ...interface{}) error {
	if level < globalLoggerLevel {
		return nil
	}

	return globalLogger.Output(2, fmt.Sprintf(format, args...))
}

type Logger interface {
	Output(callDepth int, s string) error
}

type LogLevel int

const (
	LogLevelDebug LogLevel = iota
	LogLevelInfo
	LogLevelWarning
	LogLevelError
)

func (lvl LogLevel) String() string {
	switch lvl {
	case LogLevelInfo:
		return "INFO"
	case LogLevelWarning:
		return "WARN"
	case LogLevelError:
		return "ERROR"
	}
	return "DEBUG"
}
