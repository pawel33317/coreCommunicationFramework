package logger

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/pawel33317/coreCommunicationFramework/db_handler"
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
	dbHandler   db_handler.DbLogger
	mu          sync.Mutex
}

//LoggerImp constructor
func NewLoggerImp(dbLogger db_handler.DbLogger) *LoggerImp {
	return &LoggerImp{
		isEnabled:   true,
		minLogLevel: DEBUG,
		dbHandler:   dbLogger,
	}
}

func (loggerImp *LoggerImp) logToDB(level LogLevel, ctx string, i ...interface{}) {
	var slice []string

	for _, v := range i {
		slice = append(slice, fmt.Sprint(v))
	}

	loggerImp.dbHandler.Log(time.Now().Unix(), int(level), ctx, strings.Join(slice, " "))
}

//Logger interface log method implementation
func (loggerImp *LoggerImp) Log(level LogLevel, ctx string, i ...interface{}) {
	if !loggerImp.isEnabled || loggerImp.minLogLevel > level {
		return
	}
	loggerImp.mu.Lock()
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
	loggerImp.mu.Unlock()
	if loggerImp.dbHandler != nil {
		loggerImp.logToDB(level, ctx, i...)
	}
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
