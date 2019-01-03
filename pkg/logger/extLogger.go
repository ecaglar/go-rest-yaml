//Package logger performs logging extending default log package
//Usage is logger  := logger.CreateLogger() logger.LogInfo(...) logger.LogWarning(...) logger.LogError(...)
package logger

import (
	"log"
	"os"
)

//Defines three level of logging instances which is info, warning and error
type ExtLogger struct {
	info    *log.Logger
	warning *log.Logger
	error   *log.Logger
}

//CreateLogger creates a Logger which support Info, Warning and Error level logging.
//Out channels is currently fixed but can be parameterized
//Log format is : date-time-short file
//Default channel for each level is;
//Info = 		os.Stdout
//Warning = 	os.Stdout
//Error = 	os.Stderr
func CreateLogger() *ExtLogger {

	return &ExtLogger{
		info:    log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		warning: log.New(os.Stdout, "WARNING: ", log.Ldate|log.Ltime|log.Lshortfile),
		error:   log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}
}

//LogInfo performs Info level logging to stdout
func (l *ExtLogger) LogInfo(msg ...string) {
	l.info.Println(msg)
}
//LogWarning performs warning level logging to stdout
func (l *ExtLogger) LogWarning(msg ...string) {
	l.warning.Println(msg)
}
//LogError performs warning level logging to stderr
func (l *ExtLogger) LogError(msg ...string) {
	l.error.Println(msg)
}
