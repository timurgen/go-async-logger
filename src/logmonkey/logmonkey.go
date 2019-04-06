//Package logmonkey provides functions for simple logging
package logmonkey

import (
	"fmt"
	"strings"
	"time"
)

//Logger message channel size
const LoggerBufferSize int = 1024

//time for logger graceful shutdown
const GracefulLoggerShutdownTimeMc = 100 * time.Millisecond

//Map of registered loggers
var loggers = make(map[string]*Logger)

//LogLevel numerical type
type LogLevel int

//Available log levels
const (
	TRACE LogLevel = iota
	DEBUG
	INFO
	WARNING
	ERROR
)

func (e LogLevel) String() string {
	nameMap := map[LogLevel]string{
		TRACE:   "TRACE",
		DEBUG:   "DEBUG",
		INFO:    "INFO",
		WARNING: "WARNING",
		ERROR:   "ERROR",
	}
	return nameMap[e]
}

//LogAppender interface. All appenders must implement it.
type LogAppender interface {
	ConsumeMessage(str string)
}

//LogFormatter interface. All formatters must implement it.
type LogFormatter interface {
	FormatMessage(message string, name string, level LogLevel, ts time.Time) string
}

//Default basic console appender
type ConsoleLogAppender struct {
}

//Default basic log formatter
type DefaultLogFormatter struct {
	Format string
}

//LogFormatter implementation for DefaultLogFormatter
func (lf *DefaultLogFormatter) FormatMessage(message string, name string, level LogLevel, ts time.Time) string {
	return fmt.Sprintf(lf.Format, ts.Format("2006-01-02T15:04:05.000000000"), name, level, message)
}

//LogAppender implementation for default ConsoleLogAppender
func (la *ConsoleLogAppender) ConsumeMessage(str string) {
	println(str)
}

//Logger structure
type Logger struct {
	name           string
	level          LogLevel
	appender       LogAppender
	formatter      LogFormatter
	messageChannel chan string
	closed         chan bool
}

//SetAppender - sets appender for logger
func (log *Logger) SetAppender(l LogAppender) {
	log.appender = l
}

//SetFormatter - sets formatter for logger
func (log *Logger) SetFormatter(f LogFormatter) {
	log.formatter = f
}

//SetLevel sets a level for logger
func (log *Logger) SetLevel(level LogLevel) {
	log.level = level
}

//GetLevel returns LogLevel for current Logger
func (log *Logger) GetLevel() LogLevel {
	return log.level
}

//listen starts listening logger message channel
func (log *Logger) listen() {
	for {
		select {
		case str := <-log.messageChannel:
			log.appender.ConsumeMessage(str)
		case closes := <-log.closed:
			if closes {
				//FIXME warning if messageChannel not empty
				return
			}
		}
	}
}

//Log logs a message with given level
func (log *Logger) Log(message string, level LogLevel, obj ...interface{}) {
	if log.level > level {
		return
	}
	ts := time.Now()
	formattedMessage := log.formatter.FormatMessage(message, log.name, level, ts)
	log.messageChannel <- formattedMessage
}

//Trace logs a message with TRACE level
func (log *Logger) Trace(message string, obj ...interface{}) {
	log.Log(message, TRACE, obj)
}

//Debug logs a message with DEBUG level
func (log *Logger) Debug(message string, obj ...interface{}) {
	log.Log(message, DEBUG, obj)
}

//Info logs a message with INFO level
func (log *Logger) Info(message string, obj ...interface{}) {
	log.Log(message, INFO, obj)
}

//Warning logs a message with WARNING level
func (log *Logger) Warning(message string, obj ...interface{}) {
	log.Log(message, WARNING, obj)
}

//Error logs a message with ERROR level
func (log *Logger) Error(message string, obj ...interface{}) {
	log.Log(message, ERROR, obj)
}

//GetLevelByName returns LogLevel by its name
func GetLevelByName(name string) LogLevel {
	name = strings.ToUpper(name)
	nameMap := map[string]LogLevel{
		"TRACE":   TRACE,
		"DEBUG":   DEBUG,
		"INFO":    INFO,
		"WARNING": WARNING,
		"ERROR":   ERROR,
	}
	return nameMap[name]
}

//GetLogger return logger instance associated with given name
func GetLogger(name string) *Logger {
	if _, ok := loggers[name]; !ok {
		logger := &Logger{
			name:           name,
			level:          INFO,
			appender:       &ConsoleLogAppender{},
			formatter:      &DefaultLogFormatter{Format: "%s - [%s] %s \t%s"},
			messageChannel: make(chan string, LoggerBufferSize),
			closed:         make(chan bool),
		}

		go logger.listen()
		loggers[name] = logger
	}

	return loggers[name]
}

//FlushAllLoggers wait until al loggers completes their queues or timeout is reached
//and terminates all loggers
func FlushAllLoggers() {
	flushStart := time.Now()
	timeToWait := GracefulLoggerShutdownTimeMc * time.Duration(len(loggers))

	for time.Now().Sub(flushStart)/time.Millisecond < timeToWait {
		for name, logger := range loggers {
			if len(logger.messageChannel) == 0 {
				logger.closed <- true
				delete(loggers, name)
			}
		}

		if len(loggers) == 0 {
			break
		}
	}
}
