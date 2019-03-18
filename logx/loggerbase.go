package logx

type loggerBase struct {
	level    Level
	handlers []Handler
}

func newLoggerBase() *loggerBase {
	return &loggerBase{
		handlers: make([]Handler, 0),
		level:    DefaultLevel,
	}
}

func (l *loggerBase) SetLevel(level Level) {
	l.level = level
}

func (l *loggerBase) GetLevel() Level {
	return l.level
}

func (l *loggerBase) AddHandler(h Handler) {
	l.handlers = append(l.handlers, h)
}

func (l *loggerBase) Trace(v ...interface{}) {
	if TRACE < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Trace(v...)
	}
}

func (l *loggerBase) Tracef(format string, v ...interface{}) {
	if TRACE < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Tracef(format, v...)
	}
}

func (l *loggerBase) Debug(v ...interface{}) {
	if DEBUG < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Debug(v...)
	}
}

func (l *loggerBase) Debugf(format string, v ...interface{}) {
	if DEBUG < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Debugf(format, v...)
	}
}

func (l *loggerBase) Info(v ...interface{}) {
	if INFO < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Info(v...)
	}
}

func (l *loggerBase) Infof(format string, v ...interface{}) {
	if INFO < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Infof(format, v...)
	}
}

func (l *loggerBase) Warn(v ...interface{}) {
	if WARN < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Warn(v...)
	}
}

func (l *loggerBase) Warnf(format string, v ...interface{}) {
	if WARN < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Warnf(format, v...)
	}
}

func (l *loggerBase) Error(v ...interface{}) {
	if ERROR < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Error(v...)
	}
}

func (l *loggerBase) Errorf(format string, v ...interface{}) {
	if ERROR < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Errorf(format, v...)
	}
}

func (l *loggerBase) Fatal(v ...interface{}) {
	if FATAL < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Fatal(v...)
	}
}

func (l *loggerBase) Fatalf(format string, v ...interface{}) {
	if FATAL < l.level {
		return
	}
	for _, h := range l.handlers {
		h.Fatalf(format, v...)
	}
}
