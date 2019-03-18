package logx

type Logger struct {
	lb *loggerBase
}

func NewLogger() *Logger {
	return &Logger{
		lb: newLoggerBase(),
	}
}

func (l *Logger) SetLevel(level Level) {
	l.lb.SetLevel(level)
}

func (l *Logger) GetLevel() Level {
	return l.lb.GetLevel()
}

func (l *Logger) AddHandler(h Handler) {
	l.lb.AddHandler(h)
}

func (l *Logger) Trace(v ...interface{}) {
	l.lb.Trace(v...)
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	l.lb.Tracef(format, v...)
}

func (l *Logger) Debug(v ...interface{}) {
	l.lb.Debug(v...)
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	l.lb.Debugf(format, v...)
}

func (l *Logger) Info(v ...interface{}) {
	l.lb.Info(v...)
}

func (l *Logger) Infof(format string, v ...interface{}) {
	l.lb.Infof(format, v...)
}

func (l *Logger) Warn(v ...interface{}) {
	l.lb.Warn(v...)
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	l.lb.Warnf(format, v...)
}

func (l *Logger) Error(v ...interface{}) {
	l.lb.Error(v...)
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	l.lb.Errorf(format, v...)
}

func (l *Logger) Fatal(v ...interface{}) {
	l.lb.Fatal(v...)
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	l.lb.Fatalf(format, v...)
}
