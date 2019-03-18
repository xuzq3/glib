package logx

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"time"
)

const dateFormat = "2006-01-02"

type dateRotateFileHandler struct {
	level      Level
	fmt        *formatter
	logdir     string
	logfname   string
	out        *os.File
	lock       *sync.Mutex
	callDepth  int
	time       time.Time
	maxFileCnt int
}

func NewDateRotateFileHandler(dir string, fname string, maxFileCnt int) (*dateRotateFileHandler, error) {
	err := ensureDirExist(dir)
	if err != nil {
		return nil, err
	}

	file := filepath.Join(dir, fname)
	f, err := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	fi, err := os.Stat(file)
	if err != nil {
		return nil, err
	}
	mtime := fi.ModTime()

	if maxFileCnt <= 0 {
		maxFileCnt = maxInt
	}

	fh := &dateRotateFileHandler{
		level:      DefaultLevel,
		fmt:        defaultFormatter,
		logdir:     dir,
		logfname:   fname,
		out:        f,
		lock:       new(sync.Mutex),
		callDepth:  defaultCallDepth,
		time:       mtime,
		maxFileCnt: maxFileCnt,
	}
	return fh, nil
}

func (h *dateRotateFileHandler) isNeedToRenameLogFile() bool {
	if h.time.Format(dateFormat) != time.Now().Format(dateFormat) {
		return true
	}
	return false
}

func (h *dateRotateFileHandler) renameLogFile() {
	h.out.Close()

	file := filepath.Join(h.logdir, h.logfname)
	bkfile := filepath.Join(h.logdir, fmt.Sprintf("%s.%s", h.logfname, h.time.Format(dateFormat)))
	os.Remove(bkfile)
	os.Rename(file, bkfile)

	f, _ := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	h.out = f
	h.time = time.Now()
}

func (h *dateRotateFileHandler) removeOverLogFile() {
	logFileCnt := 0
	minMtime := time.Now().Unix()
	rmfname := ""
	fis, err := ioutil.ReadDir(h.logdir)
	if err != nil {
		return
	}
	for _, fi := range fis {
		if fi.IsDir() {
			continue
		}
		if strings.HasPrefix(fi.Name(), h.logfname+".") {
			suffix := strings.TrimPrefix(fi.Name(), h.logfname+".")
			if len(suffix) == len(dateFormat) {
				if _, err := time.Parse(dateFormat, suffix); err == nil {
					logFileCnt = logFileCnt + 1
					mtime := fi.ModTime().Unix()
					if mtime < minMtime {
						minMtime = mtime
						rmfname = fi.Name()
					}
				}
			}
		}
	}
	if logFileCnt > h.maxFileCnt {
		os.Remove(filepath.Join(h.logdir, rmfname))
	}
}

func (h *dateRotateFileHandler) SetLevel(level Level) {
	h.level = level
}

func (h *dateRotateFileHandler) GetLevel() Level {
	return h.level
}

func (h *dateRotateFileHandler) SetColored(colored bool) {
	h.fmt.setColored(colored)
}

func (h *dateRotateFileHandler) GetColored() bool {
	return h.fmt.getColored()
}

func (h *dateRotateFileHandler) SetTimeFormat(timeFormat string) {
	h.fmt.setTimeFormat(timeFormat)
}

func (h *dateRotateFileHandler) GetTimeFormat() string {
	return h.fmt.getTimeFormat()
}

func (h *dateRotateFileHandler) SetCallDepth(depth int) {
	h.callDepth = depth
}

func (h *dateRotateFileHandler) GetCallDepth() int {
	return h.callDepth
}

func (h *dateRotateFileHandler) write(level Level, v ...interface{}) {
	if level < h.level {
		return
	}
	msg := fmt.Sprint(v...)
	h.output(h.callDepth, level, msg)
}

func (h *dateRotateFileHandler) writef(level Level, format string, v ...interface{}) {
	if level < h.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	h.output(h.callDepth, level, msg)
}

func (h *dateRotateFileHandler) output(callDepth int, level Level, msg string) {
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
	if h.isNeedToRenameLogFile() {
		h.renameLogFile()
		h.removeOverLogFile()
	}
	h.out.Write(buf.Bytes())
}

func (h *dateRotateFileHandler) Trace(v ...interface{}) {
	h.write(TRACE, v...)
}

func (h *dateRotateFileHandler) Tracef(format string, v ...interface{}) {
	h.writef(TRACE, format, v...)
}

func (h *dateRotateFileHandler) Debug(v ...interface{}) {
	h.write(DEBUG, v...)
}

func (h *dateRotateFileHandler) Debugf(format string, v ...interface{}) {
	h.writef(DEBUG, format, v...)
}

func (h *dateRotateFileHandler) Info(v ...interface{}) {
	h.write(INFO, v...)
}

func (h *dateRotateFileHandler) Infof(format string, v ...interface{}) {
	h.writef(INFO, format, v...)
}

func (h *dateRotateFileHandler) Warn(v ...interface{}) {
	h.write(WARN, v...)
}

func (h *dateRotateFileHandler) Warnf(format string, v ...interface{}) {
	h.writef(WARN, format, v...)
}

func (h *dateRotateFileHandler) Error(v ...interface{}) {
	h.write(ERROR, v...)
}

func (h *dateRotateFileHandler) Errorf(format string, v ...interface{}) {
	h.writef(ERROR, format, v...)
}

func (h *dateRotateFileHandler) Fatal(v ...interface{}) {
	h.write(FATAL, v...)
}

func (h *dateRotateFileHandler) Fatalf(format string, v ...interface{}) {
	h.writef(FATAL, format, v...)
}
