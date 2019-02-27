package logx

var DefaultLogger *Logger = NewLogger()

func AddHandler(h Handler) {
	DefaultLogger.AddHandler(h)
}

func SetLevel(level Level) {
	DefaultLogger.SetLevel(level)
}

func Trace(v ...interface{}) {
	DefaultLogger.Trace(v...)
}

func Tracef(format string, v ...interface{}) {
	DefaultLogger.Tracef(format, v...)
}

func Debug(v ...interface{}) {
	DefaultLogger.Debug(v...)
}

func Debugf(format string, v ...interface{}) {
	DefaultLogger.Debugf(format, v...)
}

func Info(v ...interface{}) {
	DefaultLogger.Info(v...)
}

func Infof(format string, v ...interface{}) {
	DefaultLogger.Infof(format, v...)
}

func Warn(v ...interface{}) {
	DefaultLogger.Warn(v...)
}

func Warnf(format string, v ...interface{}) {
	DefaultLogger.Warnf(format, v...)
}

func Error(v ...interface{}) {
	DefaultLogger.Error(v...)
}

func Errorf(format string, v ...interface{}) {
	DefaultLogger.Errorf(format, v...)
}

func Fatal(v ...interface{}) {
	DefaultLogger.Fatal(v...)
}

func Fatalf(format string, v ...interface{}) {
	DefaultLogger.Fatalf(format, v...)
}
