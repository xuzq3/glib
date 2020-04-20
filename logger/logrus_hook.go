package logger

import (
	"github.com/sirupsen/logrus"
	"runtime"
	"strings"
	"sync"
)

const (
	maxLogrusCallerDepth = 20
)

var (
	minLogrusCallerDepth = 0
)

type logrusCallerHook struct {
	skip int
	once sync.Once
}

func newLogrusCallerHook(skip int) *logrusCallerHook {
	return &logrusCallerHook{
		skip: skip + 1,
	}
}

func (h *logrusCallerHook) Levels() []logrus.Level {
	return logrus.AllLevels
}

func (h *logrusCallerHook) Fire(entry *logrus.Entry) error {
	caller := h.getCaller()
	entry.Caller = caller
	return nil
}

func (h *logrusCallerHook) getCaller() *runtime.Frame {
	// jump into package github.com/sirupsen/logrus
	// and get minLogrusCallerDepth
	h.once.Do(func() {
		pcs := make([]uintptr, maxLogrusCallerDepth)
		_ = runtime.Callers(0, pcs)
		for i := 0; i < maxLogrusCallerDepth; i++ {
			f, _ := runtime.CallersFrames([]uintptr{pcs[i]}).Next()
			if strings.Contains(f.Function, "github.com/sirupsen/logrus") {
				minLogrusCallerDepth = i
				break
			}
		}
	})

	pcs := make([]uintptr, maxLogrusCallerDepth)
	depth := runtime.Callers(minLogrusCallerDepth, pcs)
	frames := runtime.CallersFrames(pcs[:depth])

	var frame *runtime.Frame = nil

	skip := -1
	for {
		f, again := frames.Next()
		if skip < 0 {
			// jump out package github.com/sirupsen/logrus
			// and start to skip
			if !strings.Contains(f.Function, "github.com/sirupsen/logrus") {
				skip = 0
			}
		} else {
			skip++
		}
		if skip == h.skip {
			frame = &f
			break
		}
		if !again {
			break
		}
	}
	return frame
}
