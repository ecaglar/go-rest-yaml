package logger

import (
	"log"
	"os"
)

type LogLevel uint8

const (
	INFO LogLevel= iota
	WARNING LogLevel= iota
	ERROR LogLevel = iota
)

var LogLevelStr  = [...]string {

	"INFO",
	"WARNING",
	"ERROR",
}
func (level LogLevel) string() string{

	return LogLevelStr[level]
}

//Defines three level of logging instances which is info, warning and error
type AsyncLogger struct {
	info    *log.Logger
	warning   *log.Logger
	error     *log.Logger
	infoLogChan chan AsyncLogMsg
	warningLogChan chan AsyncLogMsg
	errorLogChan chan AsyncLogMsg

}
type AsyncLogMsg struct {

	level LogLevel
	logMsg []string

}

func CreateAsyncLogger() *AsyncLogger {

	var asyncLogger = AsyncLogger{
		info:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warning: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		error:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
		infoLogChan: make(chan AsyncLogMsg),
		warningLogChan: make(chan AsyncLogMsg),
		errorLogChan: make(chan AsyncLogMsg),
	}
	asyncLogger.startLogger()
	return &asyncLogger
}

func(l *AsyncLogger) startLogger(){
	go l.listenInfoChan()
	go l.listenWarningChan()
	go l.listenErrorChan()

}

func (l *AsyncLogger) listenInfoChan(){
	for logObj := range l.infoLogChan{
		l.info.Println(logObj.level.string(), " : ", logObj.logMsg)
	}
}
func (l *AsyncLogger) listenWarningChan(){
	for logObj := range l.warningLogChan{
		l.warning.Println(logObj.level.string(), " : ", logObj.logMsg)
	}
}
func (l *AsyncLogger) listenErrorChan(){
	for logObj := range l.errorLogChan{
		l.error.Println(logObj.level.string(), " : ", logObj.logMsg)
	}
}

func (l *AsyncLogger) Log(level LogLevel, msg ...string) {
	switch level {
	case INFO:
	default:
		go func(){l.infoLogChan <- AsyncLogMsg{level:level, logMsg:msg}}()

	case WARNING:
		go func(){l.warningLogChan <- AsyncLogMsg{level:level, logMsg:msg}}()

	case ERROR:
		go func(){l.errorLogChan <- AsyncLogMsg{level:level, logMsg:msg}}()

	}
}



