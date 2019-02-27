package logx

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type consoleHandler struct {
	level     Level
	fmt       *formatter
	out       io.Writer
	lock      *sync.Mutex
	callDepth int
}

func NewConsoleHandler() *consoleHandler {
	return &consoleHandler{
		level:     DefaultLevel,
		fmt:       defaultFormatter,
		out:       os.Stdout,
		lock:      new(sync.Mutex),
		callDepth: 4,
	}
}

func (h *consoleHandler) SetLevel(level Level) {
	h.level = level
}

func (h *consoleHandler) GetLevel() Level {
	return h.level
}

func (h *consoleHandler) SetColored(colored bool) {
	h.fmt.setColored(colored)
}

func (h *consoleHandler) GetColored() bool {
	return h.fmt.getColored()
}

func (h *consoleHandler) SetTimeFormat(timeFormat string) {
	h.fmt.setTimeFormat(timeFormat)
}

func (h *consoleHandler) GetTimeFormat() string {
	return h.fmt.getTimeFormat()
}

func (h *consoleHandler) write(level Level, v ...interface{}) {
	if level < h.level {
		return
	}
	msg := fmt.Sprint(v...)
	h.output(h.callDepth, level, msg)
}

func (h *consoleHandler) writef(level Level, format string, v ...interface{}) {
	if level < h.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	h.output(h.callDepth, level, msg)
}

func (h *consoleHandler) output(callDepth int, level Level, msg string) {
	_, file, line, ok := runtime.Caller(callDepth)
	if !ok {
		file = "???"
		line = 0
	}
	_, file = filepath.Split(file)
	now := time.Now()
	buf := h.fmt.format(now, level, file, line, msg)

	h.lock.Lock()
	defer h.lock.Unlock()
	h.out.Write(buf.Bytes())
}

func (h *consoleHandler) Trace(v ...interface{}) {
	h.write(TRACE, v...)
}

func (h *consoleHandler) Tracef(format string, v ...interface{}) {
	h.writef(TRACE, format, v...)
}

func (h *consoleHandler) Debug(v ...interface{}) {
	h.write(DEBUG, v...)
}

func (h *consoleHandler) Debugf(format string, v ...interface{}) {
	h.writef(DEBUG, format, v...)
}

func (h *consoleHandler) Info(v ...interface{}) {
	h.write(INFO, v...)
}

func (h *consoleHandler) Infof(format string, v ...interface{}) {
	h.writef(INFO, format, v...)
}

func (h *consoleHandler) Warn(v ...interface{}) {
	h.write(WARN, v...)
}

func (h *consoleHandler) Warnf(format string, v ...interface{}) {
	h.writef(WARN, format, v...)
}

func (h *consoleHandler) Error(v ...interface{}) {
	h.write(ERROR, v...)
}

func (h *consoleHandler) Errorf(format string, v ...interface{}) {
	h.writef(ERROR, format, v...)
}

func (h *consoleHandler) Fatal(v ...interface{}) {
	h.write(FATAL, v...)
}

func (h *consoleHandler) Fatalf(format string, v ...interface{}) {
	h.writef(FATAL, format, v...)
}
