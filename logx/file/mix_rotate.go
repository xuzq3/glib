package file

import (
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"sync"
	"time"
)

type MixRotateFile struct {
	log        *lumberjack.Logger
	lastTime   time.Time
	rotateTime time.Duration
	locker     sync.Mutex
}

func NewMixRotateFile(filename string, maxSize int, rotateTime time.Duration, maxAge int, maxBackups int) *MixRotateFile {
	log := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxAge:     maxAge,
		MaxBackups: maxBackups,
		LocalTime:  true,
	}
	mtime := getFileModTime(filename)

	logger := &MixRotateFile{
		log:        log,
		lastTime:   mtime,
		rotateTime: rotateTime,
		locker:     sync.Mutex{},
	}
	return logger
}

func (l *MixRotateFile) Write(p []byte) (n int, err error) {
	if !l.lastTime.IsZero() && l.isOverTime() {
		l.locker.Lock()
		defer l.locker.Unlock()
		if l.isOverTime() {
			err := l.Rotate()
			if err != nil {
				return 0, err
			}
		}
	}

	l.lastTime = time.Now()
	return l.log.Write(p)
}

func (l *MixRotateFile) isOverTime() bool {
	if l.rotateTime <= 0 {
		return false
	}
	nt := time.Now().Truncate(l.rotateTime)
	lt := l.lastTime.Truncate(l.rotateTime)
	if nt != lt {
		return true
	}
	return false
}

func (l *MixRotateFile) Rotate() error {
	l.lastTime = time.Now()
	return l.log.Rotate()
}

func (l *MixRotateFile) Close() error {
	return l.log.Close()
}

func getFileModTime(filename string) time.Time {
	f, err := os.Stat(filename)
	if err == nil {
		return f.ModTime()
	}
	return time.Unix(0, 0)
}
