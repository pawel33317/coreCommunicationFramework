package logger

import (
	"fmt"
)

type LogLevel int

const (
	VERBOSE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

func (level LogLevel) ToString() string {
	switch level {
	case VERBOSE:
		return "VERBOSE"
	case DEBUG:
		return "DEBUG"
	case INFO:
		return "INFO"
	case WARN:
		return "WARN"
	case ERROR:
		return "ERROR"
	case FATAL:
		return "FATAL"
	default:
		return "<UNKNOWN>"
	}
}

type LogWrapper struct {
	logger Logger
	ctx    string
}

func (lw *LogWrapper) Log(ll LogLevel, data ...interface{}) {
	lw.logger.Log(ll, lw.ctx, data)
}

type Logger interface {
	Log(LogLevel, string, ...interface{})
}

type LoggerImp struct {
	isEnabled   bool
	minLogLevel LogLevel
}

func NewLoggerImp() *LoggerImp {
	return &LoggerImp{
		isEnabled:   true,
		minLogLevel: DEBUG,
	}
}

func (loggerImp *LoggerImp) Log(level LogLevel, ctx string, i ...interface{}) {
	if !loggerImp.isEnabled || loggerImp.minLogLevel > level {
		return
	}
	fmt.Print("[", ctx, "]")
	iter := 4 - len(ctx)
	for iter > 0 {
		fmt.Print(" ")
		iter--
	}
	fmt.Print("[", level.ToString(), "]")
	for _, v := range i {
		fmt.Print(" ", v)
	}
	fmt.Println()
}

func (loggerImp *LoggerImp) Enable() {
	loggerImp.isEnabled = true
}

func (loggerImp *LoggerImp) Disable() {
	loggerImp.isEnabled = false
}

func (loggerImp *LoggerImp) SetMinLogLevel(level LogLevel) {
	loggerImp.minLogLevel = level
}
