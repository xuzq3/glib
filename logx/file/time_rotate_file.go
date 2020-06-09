package file

import (
	"fmt"
	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"path/filepath"
	"time"
)

func getRotateFormat(rotateTime time.Duration) string {
	var format string
	if rotateTime.Hours() >= 24 {
		format = "%Y-%m-%d"
	} else if rotateTime.Hours() >= 1 {
		format = "%Y-%m-%dT%H"
	} else if rotateTime.Minutes() >= 1 {
		format = "%Y-%m-%dT%H-%M"
	} else {
		format = "%Y-%m-%dT%H-%M-%S"
	}
	return format
}

func getRotateTime(rotateTime time.Duration) time.Duration {
	var d time.Duration
	if rotateTime.Hours() >= 24 {
		d = rotateTime.Round(time.Hour * 24)
	} else if rotateTime.Hours() >= 1 {
		d = rotateTime.Round(time.Hour)
	} else if rotateTime.Minutes() >= 1 {
		d = rotateTime.Round(time.Minute)
	} else {
		d = rotateTime.Round(time.Second)
	}
	return d
}

func genFilename(name string, rotateFormat string) string {
	dir := filepath.Dir(name)
	filename := filepath.Base(name)
	ext := filepath.Ext(filename)
	prefix := filename[:len(filename)-len(ext)]
	return filepath.Join(dir, fmt.Sprintf("%s-%s%s", prefix, rotateFormat, ext))
}

type TimeRotateFile struct {
	log *rotatelogs.RotateLogs
}

// rotateTime: the number of seconds between rotation
func NewTimeRotateFile(filename string, rotateTime time.Duration, rotateCount int, rotateMaxAge int) (*TimeRotateFile, error) {
	format := getRotateFormat(rotateTime)
	realtime := getRotateTime(rotateTime)
	realname := genFilename(filename, format)
	options := []rotatelogs.Option{
		rotatelogs.WithRotationTime(realtime),
		rotatelogs.WithMaxAge(time.Hour * 24 * time.Duration(rotateMaxAge)),
		rotatelogs.WithRotationCount(uint(rotateCount)),
	}
	log, err := rotatelogs.New(realname, options...)
	if err != nil {
		return nil, err
	}

	logger := &TimeRotateFile{
		log: log,
	}
	return logger, nil
}

func (l *TimeRotateFile) Write(p []byte) (n int, err error) {
	return l.log.Write(p)
}

func (l *TimeRotateFile) Rotate() error {
	return l.log.Rotate()
}

func (l *TimeRotateFile) Close() error {
	return l.log.Close()
}
