package logx

import (
	"fmt"
	"github.com/sirupsen/logrus"
	"io"
	"path/filepath"
	"runtime"
	"strings"
)

type LogrusLogger struct {
	opt    *Option
	logger *logrus.Logger
}

func NewLogrusLogger(opt *Option) *LogrusLogger {
	l := &LogrusLogger{
		opt: opt,
	}
	l.initLogger()
	return l
}

func (l *LogrusLogger) initLogger() {
	logger := logrus.New()

	// set level
	var level logrus.Level
	switch strings.ToLower(l.opt.level) {
	case DebugLevel:
		level = logrus.DebugLevel
	case InfoLevel:
		level = logrus.InfoLevel
	case WarnLevel:
		level = logrus.WarnLevel
	case ErrorLevel:
		level = logrus.ErrorLevel
	default:
		level = logrus.DebugLevel
	}
	logger.SetLevel(level)

	// set output
	outs := make([]io.Writer, 0)
	//if !l.opt.disableConsole {
	//	outs = append(outs, os.Stdout)
	//}
	if len(l.opt.outs) > 0 {
		for _, out := range l.opt.outs {
			outs = append(outs, out)
		}
	}
	logger.SetOutput(io.MultiWriter(outs...))

	// set caller
	//logger.SetReportCaller(true)
	logger.AddHook(newLogrusCallerHook(l.opt.callerSkip))

	// set formatter
	if l.opt.jsonFormat {
		logger.SetFormatter(&logrusJsonFormatter{})
		//logger.SetFormatter(&logrus.JSONFormatter{
		//	TimestampFormat:  "2006-01-02T15:04:05.000",
		//	CallerPrettyfier: logrusCallerPrettyfier,
		//	DataKey:          "data",
		//})
	} else {
		logger.SetFormatter(&logrus.TextFormatter{
			TimestampFormat: "2006-01-02T15:04:05.000",
			//CallerPrettyfier: logrusCallerPrettyfier,
			DisableSorting: true,
		})
	}
	l.logger = logger
}

func logrusCallerPrettyfier(frame *runtime.Frame) (string, string) {
	funcVal := ""
	fileVal := ""
	if frame != nil {
		if frame.Function != "" {
			funcVal = filepath.Base(frame.Function)
		}
		if frame.File != "" {
			fileVal = fmt.Sprintf("%s:%d", filepath.Base(frame.File), frame.Line)
		}
	}
	return funcVal, fileVal
}

func (l *LogrusLogger) Debug(args ...interface{}) {
	l.logger.Debug(args...)
}

func (l *LogrusLogger) Debugf(format string, args ...interface{}) {
	l.logger.Debugf(format, args...)
}

func (l *LogrusLogger) Info(args ...interface{}) {
	l.logger.Info(args...)
}

func (l *LogrusLogger) Infof(format string, args ...interface{}) {
	l.logger.Infof(format, args...)
}

func (l *LogrusLogger) Warn(args ...interface{}) {
	l.logger.Warn(args...)
}

func (l *LogrusLogger) Warnf(format string, args ...interface{}) {
	l.logger.Warnf(format, args...)
}

func (l *LogrusLogger) Error(args ...interface{}) {
	l.logger.Error(args...)
}

func (l *LogrusLogger) Errorf(format string, args ...interface{}) {
	l.logger.Errorf(format, args...)
}

func (l *LogrusLogger) Fatal(args ...interface{}) {
	l.logger.Fatal(args...)
}

func (l *LogrusLogger) Fatalf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *LogrusLogger) Panic(args ...interface{}) {
	l.logger.Panic(args...)
}

func (l *LogrusLogger) Panicf(format string, args ...interface{}) {
	l.logger.Fatalf(format, args...)
}

func (l *LogrusLogger) WithFields(fields Fields) ILogger {
	return &logrusLogEntry{
		logger: l,
		entry:  l.logger.WithFields(convertToLogrusFields(fields)),
	}
}

func (l *LogrusLogger) WithError(err error) ILogger {
	return &logrusLogEntry{
		logger: l,
		entry:  l.logger.WithError(err),
	}
}

func (l *LogrusLogger) Output() io.Writer {
	return io.MultiWriter(l.opt.outs...)
}

type logrusLogEntry struct {
	logger *LogrusLogger
	entry  *logrus.Entry
}

func (e *logrusLogEntry) Debug(args ...interface{}) {
	e.entry.Debug(args...)
}

func (e *logrusLogEntry) Debugf(format string, args ...interface{}) {
	e.entry.Debugf(format, args...)
}

func (e *logrusLogEntry) Info(args ...interface{}) {
	e.entry.Info(args...)
}

func (e *logrusLogEntry) Infof(format string, args ...interface{}) {
	e.entry.Infof(format, args...)
}

func (e *logrusLogEntry) Warn(args ...interface{}) {
	e.entry.Warn(args...)
}

func (e *logrusLogEntry) Warnf(format string, args ...interface{}) {
	e.entry.Warnf(format, args...)
}

func (e *logrusLogEntry) Error(args ...interface{}) {
	e.entry.Error(args...)
}

func (e *logrusLogEntry) Errorf(format string, args ...interface{}) {
	e.entry.Errorf(format, args...)
}

func (e *logrusLogEntry) Fatal(args ...interface{}) {
	e.entry.Fatal(args...)
}

func (e *logrusLogEntry) Fatalf(format string, args ...interface{}) {
	e.entry.Fatalf(format, args...)
}

func (e *logrusLogEntry) Panic(args ...interface{}) {
	e.entry.Panic(args...)
}

func (e *logrusLogEntry) Panicf(format string, args ...interface{}) {
	e.entry.Fatalf(format, args...)
}

func (e *logrusLogEntry) WithFields(fields Fields) ILogger {
	return &logrusLogEntry{
		logger: e.logger,
		entry:  e.entry.WithFields(convertToLogrusFields(fields)),
	}
}

func (e *logrusLogEntry) WithError(err error) ILogger {
	return &logrusLogEntry{
		logger: e.logger,
		entry:  e.entry.WithError(err),
	}
}

func (e *logrusLogEntry) Output() io.Writer {
	return e.logger.Output()
}

func convertToLogrusFields(fields Fields) logrus.Fields {
	logrusFields := logrus.Fields{}
	for index, val := range fields {
		logrusFields[index] = val
	}
	return logrusFields
}
