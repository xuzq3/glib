package logx

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type fileHandler struct {
	level     Level
	fmt       *formatter
	logfile   string
	out       *os.File
	lock      *sync.Mutex
	callDepth int
}

func NewFileHandler(file string) (*fileHandler, error) {
	f, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}
	fh := &fileHandler{
		level:     DefaultLevel,
		fmt:       defaultFormatter,
		logfile:   file,
		out:       f,
		lock:      new(sync.Mutex),
		callDepth: defaultCallDepth,
	}
	return fh, nil
}

func (h *fileHandler) SetLevel(level Level) {
	h.level = level
}

func (h *fileHandler) GetLevel() Level {
	return h.level
}

func (h *fileHandler) SetColored(colored bool) {
	h.fmt.setColored(colored)
}

func (h *fileHandler) GetColored() bool {
	return h.fmt.getColored()
}

func (h *fileHandler) SetTimeFormat(timeFormat string) {
	h.fmt.setTimeFormat(timeFormat)
}

func (h *fileHandler) GetTimeFormat() string {
	return h.fmt.getTimeFormat()
}

func (h *fileHandler) SetCallDepth(depth int) {
	h.callDepth = depth
}

func (h *fileHandler) GetCallDepth() int {
	return h.callDepth
}

func (h *fileHandler) write(level Level, v ...interface{}) {
	if level < h.level {
		return
	}
	msg := fmt.Sprint(v...)
	h.output(h.callDepth, level, msg)
}

func (h *fileHandler) writef(level Level, format string, v ...interface{}) {
	if level < h.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	h.output(h.callDepth, level, msg)
}

func (h *fileHandler) output(callDepth int, level Level, msg string) {
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

func (h *fileHandler) Trace(v ...interface{}) {
	h.write(TRACE, v...)
}

func (h *fileHandler) Tracef(format string, v ...interface{}) {
	h.writef(TRACE, format, v...)
}

func (h *fileHandler) Debug(v ...interface{}) {
	h.write(DEBUG, v...)
}

func (h *fileHandler) Debugf(format string, v ...interface{}) {
	h.writef(DEBUG, format, v...)
}

func (h *fileHandler) Info(v ...interface{}) {
	h.write(INFO, v...)
}

func (h *fileHandler) Infof(format string, v ...interface{}) {
	h.writef(INFO, format, v...)
}

func (h *fileHandler) Warn(v ...interface{}) {
	h.write(WARN, v...)
}

func (h *fileHandler) Warnf(format string, v ...interface{}) {
	h.writef(WARN, format, v...)
}

func (h *fileHandler) Error(v ...interface{}) {
	h.write(ERROR, v...)
}

func (h *fileHandler) Errorf(format string, v ...interface{}) {
	h.writef(ERROR, format, v...)
}

func (h *fileHandler) Fatal(v ...interface{}) {
	h.write(FATAL, v...)
}

func (h *fileHandler) Fatalf(format string, v ...interface{}) {
	h.writef(FATAL, format, v...)
}
