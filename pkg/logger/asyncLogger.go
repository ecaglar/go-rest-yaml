/*
Package logger defines a async logger using channels.
Three different log level has been defined which are INFO,WARNING ,ERROR and FATAL
INFO and WARNING logs are sent to Stdout and ERROR FATAL logs to StdErr by default
*/
package logger

import (
	"log"
	"os"
)

type LogLevel uint8

//LogLevel defines log levels that can be used.
const (
	INFO    LogLevel = iota
	WARNING LogLevel = iota
	ERROR   LogLevel = iota
	FATAL   LogLevel = iota
)

//LogLevelStr defines log levels that can be used.
var LogLevelStr = [...]string{

	"INFO",
	"WARNING",
	"ERROR",
	"FATAL",
}

func (level LogLevel) string() string {

	return LogLevelStr[level]
}

/*
AsyncLogger objects defines logger and channel for each log level.
Since loggers work in async manner, it is possible to stop logger as well.
for that purpose, stop channel has been defined. when logger reads a value
from stop channel, it returns.
*/
type AsyncLogger struct {
	info           *log.Logger
	warning        *log.Logger
	error          *log.Logger
	fatal          *log.Logger
	infoLogChan    chan AsyncLogMsg
	warningLogChan chan AsyncLogMsg
	errorLogChan   chan AsyncLogMsg
	fatalLogChan   chan AsyncLogMsg

	//stop signal
	stop chan bool
}

/*
Defines the structure of log messages.
level is a constant to identify log level
and logMsg is the actual log message.
*/
type AsyncLogMsg struct {
	level  LogLevel
	logMsg []string
}

//CreateAsyncLogger creates a async logger instance and initialize the channels.
//The purpose for using different channels for different levels is to provide flexibility
//so that we can implement different logic for different levels if needed.
//channels also start being listened as soon as we create logger instance.
func CreateAsyncLogger() *AsyncLogger {

	var asyncLogger = AsyncLogger{
		info:           log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warning:        log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		error:          log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		fatal:          log.New(os.Stderr, "FATAL: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogChan:    make(chan AsyncLogMsg),
		warningLogChan: make(chan AsyncLogMsg),
		errorLogChan:   make(chan AsyncLogMsg),
		fatalLogChan:   make(chan AsyncLogMsg),
		stop:           make(chan bool),
	}
	asyncLogger.startLogger()
	return &asyncLogger
}

//startLogger initiates a go route that waits receiving data from logger channels.
func (l *AsyncLogger) startLogger() {
	go l.listen()
}

//listen initiates a infinite loop to listen log channels until stop message received.
func (l *AsyncLogger) listen() {
	for {
		select {
		case logMsg := <-l.infoLogChan:
			l.info.Println(logMsg.level.string(), " : ", logMsg.logMsg)
		case logMsg := <-l.warningLogChan:
			l.info.Println(logMsg.level.string(), " : ", logMsg.logMsg)
		case logMsg := <-l.errorLogChan:
			l.info.Println(logMsg.level.string(), " : ", logMsg.logMsg)
		case logMsg := <-l.fatalLogChan:
			l.fatal.Println(logMsg.level.string(), " : ", logMsg.logMsg)
			os.Exit(1)
		case <-l.stop:
			return
		}
	}
}

//Log function performs actual logging by passing log message into related channel.
//Gets log level and log message as arguments.
func (l *AsyncLogger) Log(level LogLevel, msg ...string) {
	switch level {
	case INFO:
		go func() { l.infoLogChan <- AsyncLogMsg{level: level, logMsg: msg} }()

	case WARNING:
		go func() { l.warningLogChan <- AsyncLogMsg{level: level, logMsg: msg} }()

	case ERROR:
		go func() { l.errorLogChan <- AsyncLogMsg{level: level, logMsg: msg} }()

	case FATAL:
		go func() { l.fatalLogChan <- AsyncLogMsg{level: level, logMsg: msg} }()

	}
}

//Stop function is responsible for ending logging loop.
func (l *AsyncLogger) Stop() {
	go func() { l.stop <- true }()
}
