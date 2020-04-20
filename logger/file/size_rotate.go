package file

import (
	"gopkg.in/natefinch/lumberjack.v2"
)

type SizeRotateFile struct {
	log *lumberjack.Logger
}

// maxSize: the maximum size in megabytes of the log file before it gets rotated
func NewSizeRotateFile(filename string, maxSize int, maxBackups int, maxAge int) *SizeRotateFile {
	log := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    maxSize,
		MaxBackups: maxBackups,
		MaxAge:     maxAge,
		LocalTime:  true,
	}
	logger := &SizeRotateFile{
		log: log,
	}
	return logger
}

func (l *SizeRotateFile) Write(p []byte) (n int, err error) {
	return l.log.Write(p)
}

func (l *SizeRotateFile) Rotate() error {
	return l.log.Rotate()
}

func (l *SizeRotateFile) Close() error {
	return l.log.Close()
}
