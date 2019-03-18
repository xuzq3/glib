package logx

import (
	"bytes"
	"runtime"
	"strconv"
	"time"
)

type formatter struct {
	colored    bool
	timeFormat string
}

func newFormater() *formatter {
	return &formatter{
		colored:    defaultColored,
		timeFormat: defaultTimeFormat,
	}
}

func (f *formatter) setColored(colored bool) {
	f.colored = colored
}

func (f *formatter) getColored() bool {
	return f.colored
}

func (f *formatter) setTimeFormat(timeFormat string) {
	f.timeFormat = timeFormat
}

func (f *formatter) getTimeFormat() string {
	return f.timeFormat
}

func (f *formatter) format(t time.Time, lv Level, file string, line int, msg string) *bytes.Buffer {
	buf := new(bytes.Buffer)
	if f.colored {
		buf.WriteString("\033[")
		buf.WriteString(strconv.Itoa(int(lv.Color())))
		buf.WriteByte('m')
	}
	buf.WriteByte('[')
	buf.WriteString(lv.String())
	buf.WriteByte(']')
	buf.WriteByte(':')
	buf.WriteByte('[')
	buf.WriteString(t.Format(f.timeFormat))
	buf.WriteByte(']')
	buf.WriteByte(':')
	buf.WriteByte('[')
	buf.WriteString(file)
	buf.WriteByte(':')
	buf.WriteString(strconv.Itoa(line))
	buf.WriteByte(']')
	buf.WriteByte(' ')
	buf.WriteString(msg)
	if f.colored {
		buf.WriteString("\033[0m")
	}
	buf.WriteByte('\n')
	// if runtime.GOOS == "windows" {
	// 	buf.WriteByte('\n')
	// } else {
	// 	buf.WriteString("\r\n")
	// }
	return buf
}

var defaultColored bool = func() bool {
	if runtime.GOOS == "windows" {
		return false
	} else {
		return true
	}
}()

var (
	defaultTimeFormat string = "2006-01-02 15:04:05"
	defaultCallDepth  int    = 4
)

var defaultFormatter *formatter = newFormater()
