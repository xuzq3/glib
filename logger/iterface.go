package logger

import "io"

type Fields map[string]interface{}

type ILogger interface {
	Debug(args ...interface{})
	Debugf(format string, args ...interface{})
	Info(args ...interface{})
	Infof(format string, args ...interface{})
	Warn(args ...interface{})
	Warnf(format string, args ...interface{})
	Error(args ...interface{})
	Errorf(format string, args ...interface{})
	Panic(args ...interface{})
	Panicf(format string, args ...interface{})
	Fatal(args ...interface{})
	Fatalf(format string, args ...interface{})
	WithFields(fields Fields) ILogger
	WithError(err error) ILogger
	Output() io.Writer
}

const (
	//DebugLevel has verbose message
	DebugLevel = "debug"
	//InfoLevel is default log level
	InfoLevel = "info"
	//WarnLevel is for logging messages about possible issues
	WarnLevel = "warn"
	//ErrorLevel is for logging errors
	ErrorLevel = "error"
	////FatalLevel is for logging fatal messages. The system shutdowns after logging the message.
	//FatalLevel = "fatal"
	////PanicLevel is for logging panic messages. The system panics after logging the message.
	//PanicLevel = "panic"
)
