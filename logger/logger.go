package logger

import (
	"fmt"
)

/*
 * LogLevel states
 *   - 6 states
 *   - ToString()
 * Logger interface - allows to log traces
 *   - Log(ll, ctx, data)
 * LoggerImp - implementation of Logger interface - currently logs to cli
 *   - NewLoggerImp()
 *   - Log(ll, ctx, data)
 *   - Enable(bool)
 *   - Disable(bool)
 *   - SetMinLogLevel(ll)
 * LogWrapper struct - keeps ctx and forwards data to logger
 *   - Log(ll, data)
 *   - NewLogWrapper(l, ctx)
 */

//logging level
type LogLevel int

//logging level states
const (
	VERBOSE LogLevel = iota
	DEBUG
	INFO
	WARN
	ERROR
	FATAL
)

//log level converter to string
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

//wrapper for logger, allows to keeps ctx inside
type LogWrapper struct {
	logger Logger //logger
	ctx    string //log context
}

//forwards data to logger and add context
func (lw *LogWrapper) Log(ll LogLevel, data ...interface{}) {
	lw.logger.Log(ll, lw.ctx, data...)
}

//LogWrapper constructor
func NewLogWrapper(l Logger, c string) *LogWrapper {
	return &LogWrapper{
		logger: l,
		ctx:    c,
	}
}

//Logger interface - allows to log traces
type Logger interface {
	Log(LogLevel, string, ...interface{})
}

//Logger interface implementation
type LoggerImp struct {
	isEnabled   bool
	minLogLevel LogLevel
}

//LoggerImp constructor
func NewLoggerImp() *LoggerImp {
	return &LoggerImp{
		isEnabled:   true,
		minLogLevel: DEBUG,
	}
}

//Logger interface log method implementation
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

//Enables Logger
func (loggerImp *LoggerImp) Enable() {
	loggerImp.isEnabled = true
}

//Disable Logger
func (loggerImp *LoggerImp) Disable() {
	loggerImp.isEnabled = false
}

//Set minimum log level
func (loggerImp *LoggerImp) SetMinLogLevel(level LogLevel) {
	loggerImp.minLogLevel = level
}
