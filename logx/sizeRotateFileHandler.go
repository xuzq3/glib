package logx

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"sync"
	"time"
)

type sizeRotateFileHandler struct {
	level       Level
	fmt         *formatter
	logdir      string
	logfname    string
	out         *os.File
	lock        *sync.Mutex
	callDepth   int
	maxFileCnt  int
	maxFileSize int64
	curFileSize int64
	nextFileIdx int
}

func NewSizeRotateFileHandler(dir string, fname string, maxFileCnt int, maxFileSize int64) (*sizeRotateFileHandler, error) {
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
	curFileSize := fi.Size()

	nextFileIdx := getNextFileIdx(dir, fname, maxFileCnt)

	if maxFileCnt <= 0 {
		maxFileCnt = maxInt
	}

	fh := &sizeRotateFileHandler{
		level:       DefaultLevel,
		fmt:         defaultFormatter,
		logdir:      dir,
		logfname:    fname,
		out:         f,
		lock:        new(sync.Mutex),
		callDepth:   defaultCallDepth,
		maxFileCnt:  maxFileCnt,
		maxFileSize: maxFileSize,
		curFileSize: curFileSize,
		nextFileIdx: nextFileIdx,
	}
	return fh, nil
}

func getNextFileIdx(dir string, fname string, maxFileCnt int) int {
	minMtime := time.Now().Unix()
	minIdx := 1
	for i := 1; i < maxFileCnt; i++ {
		file := filepath.Join(dir, fmt.Sprintf("%s.%d", fname, i))
		fi, err := os.Stat(file)
		if err != nil {
			minIdx = i
			break
		}
		mtime := fi.ModTime().Unix()
		if mtime < minMtime {
			minMtime = mtime
			minIdx = i
		}
	}
	return minIdx
}

func (h *sizeRotateFileHandler) isNeedToRenameLogFile() bool {
	if h.curFileSize > h.maxFileSize {
		return true
	}
	return false
}

func (h *sizeRotateFileHandler) renameLogFile() {
	h.out.Close()

	file := filepath.Join(h.logdir, h.logfname)
	bkfile := filepath.Join(h.logdir, fmt.Sprintf("%s.%d", h.logfname, h.nextFileIdx))
	os.Remove(bkfile)
	os.Rename(file, bkfile)

	f, _ := os.OpenFile(file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0666)
	h.out = f

	h.nextFileIdx = h.nextFileIdx%h.maxFileCnt + 1
	h.curFileSize = 0
}

func (h *sizeRotateFileHandler) SetLevel(level Level) {
	h.level = level
}

func (h *sizeRotateFileHandler) GetLevel() Level {
	return h.level
}

func (h *sizeRotateFileHandler) SetColored(colored bool) {
	h.fmt.setColored(colored)
}

func (h *sizeRotateFileHandler) GetColored() bool {
	return h.fmt.getColored()
}

func (h *sizeRotateFileHandler) SetTimeFormat(timeFormat string) {
	h.fmt.setTimeFormat(timeFormat)
}

func (h *sizeRotateFileHandler) GetTimeFormat() string {
	return h.fmt.getTimeFormat()
}

func (h *sizeRotateFileHandler) SetCallDepth(depth int) {
	h.callDepth = depth
}

func (h *sizeRotateFileHandler) GetCallDepth() int {
	return h.callDepth
}

func (h *sizeRotateFileHandler) write(level Level, v ...interface{}) {
	if level < h.level {
		return
	}
	msg := fmt.Sprint(v...)
	h.output(h.callDepth, level, msg)
}

func (h *sizeRotateFileHandler) writef(level Level, format string, v ...interface{}) {
	if level < h.level {
		return
	}
	msg := fmt.Sprintf(format, v...)
	h.output(h.callDepth, level, msg)
}

func (h *sizeRotateFileHandler) output(callDepth int, level Level, msg string) {
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
	}
	h.out.Write(buf.Bytes())
	h.curFileSize = h.curFileSize + int64(buf.Len())
}

func (h *sizeRotateFileHandler) Trace(v ...interface{}) {
	h.write(TRACE, v...)
}

func (h *sizeRotateFileHandler) Tracef(format string, v ...interface{}) {
	h.writef(TRACE, format, v...)
}

func (h *sizeRotateFileHandler) Debug(v ...interface{}) {
	h.write(DEBUG, v...)
}

func (h *sizeRotateFileHandler) Debugf(format string, v ...interface{}) {
	h.writef(DEBUG, format, v...)
}

func (h *sizeRotateFileHandler) Info(v ...interface{}) {
	h.write(INFO, v...)
}

func (h *sizeRotateFileHandler) Infof(format string, v ...interface{}) {
	h.writef(INFO, format, v...)
}

func (h *sizeRotateFileHandler) Warn(v ...interface{}) {
	h.write(WARN, v...)
}

func (h *sizeRotateFileHandler) Warnf(format string, v ...interface{}) {
	h.writef(WARN, format, v...)
}

func (h *sizeRotateFileHandler) Error(v ...interface{}) {
	h.write(ERROR, v...)
}

func (h *sizeRotateFileHandler) Errorf(format string, v ...interface{}) {
	h.writef(ERROR, format, v...)
}

func (h *sizeRotateFileHandler) Fatal(v ...interface{}) {
	h.write(FATAL, v...)
}

func (h *sizeRotateFileHandler) Fatalf(format string, v ...interface{}) {
	h.writef(FATAL, format, v...)
}
