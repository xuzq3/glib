package logx

var defaultLogger *loggerBase = newLoggerBase()

func AddHandler(h Handler) {
	defaultLogger.AddHandler(h)
}

func SetLevel(level Level) {
	defaultLogger.SetLevel(level)
}

func Trace(v ...interface{}) {
	defaultLogger.Trace(v...)
}

func Tracef(format string, v ...interface{}) {
	defaultLogger.Tracef(format, v...)
}

func Debug(v ...interface{}) {
	defaultLogger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	defaultLogger.Debugf(format, v...)
}

func Info(v ...interface{}) {
	defaultLogger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	defaultLogger.Infof(format, v...)
}

func Warn(v ...interface{}) {
	defaultLogger.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	defaultLogger.Warnf(format, v...)
}

func Error(v ...interface{}) {
	defaultLogger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	defaultLogger.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	defaultLogger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	defaultLogger.Fatalf(format, v...)
}
