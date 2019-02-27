package logx

type Logger struct {
	level    Level
	handlers []Handler
}

func NewLogger() *Logger {
	return &Logger{
		level:    DefaultLevel,
		handlers: make([]Handler, 0),
	}
}

func (l *Logger) SetLevel(level Level) {
	l.level = level
}

func (l *Logger) GetLevel() Level {
	return l.level
}

func (l *Logger) AddHandler(h Handler) {
	l.handlers = append(l.handlers, h)
}

func (l *Logger) Trace(v ...interface{}) {
	if TRACE < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Trace(v...)
	}
}

func (l *Logger) Tracef(format string, v ...interface{}) {
	if TRACE < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Tracef(format, v...)
	}
}

func (l *Logger) Debug(v ...interface{}) {
	if DEBUG < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Debug(v...)
	}
}

func (l *Logger) Debugf(format string, v ...interface{}) {
	if DEBUG < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Debugf(format, v...)
	}
}

func (l *Logger) Info(v ...interface{}) {
	if INFO < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Info(v...)
	}
}

func (l *Logger) Infof(format string, v ...interface{}) {
	if INFO < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Infof(format, v...)
	}
}

func (l *Logger) Warn(v ...interface{}) {
	if WARN < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Warn(v...)
	}
}

func (l *Logger) Warnf(format string, v ...interface{}) {
	if WARN < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Warnf(format, v...)
	}
}

func (l *Logger) Error(v ...interface{}) {
	if ERROR < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Error(v...)
	}
}

func (l *Logger) Errorf(format string, v ...interface{}) {
	if ERROR < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Errorf(format, v...)
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	if FATAL < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Fatal(v...)
	}
}

func (l *Logger) Fatalf(format string, v ...interface{}) {
	if FATAL < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Fatalf(format, v...)
	}
}
