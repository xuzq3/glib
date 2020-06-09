package logx

import (
	"io"
)

const (
	levelKey  = "level"
	timeKey   = "time"
	callerKey = "caller"
	msgKey    = "msg"
	dataKey   = "data"
)

type Option struct {
	level string
	outs  []io.Writer
	//disableConsole bool
	jsonFormat bool
	callerSkip int
}

func NewOption() *Option {
	return &Option{
		level: DebugLevel,
	}
}

//func (o *Option) DisableConsole() *Option {
//	o.disableConsole = true
//	return o
//}
//
//func (o *Option) EnableConsole() *Option {
//	o.disableConsole = false
//	return o
//}

func (o *Option) AddOutput(out io.Writer) *Option {
	o.outs = append(o.outs, out)
	return o
}

func (o *Option) AddOutputs(out ...io.Writer) *Option {
	o.outs = append(o.outs, out...)
	return o
}

func (o *Option) SetJsonFormat() *Option {
	o.jsonFormat = true
	return o
}

func (o *Option) SetTextFormat() *Option {
	o.jsonFormat = false
	return o
}

func (o *Option) SetLevel(level string) *Option {
	o.level = level
	return o
}

func (o *Option) AddCallerSkip(skip int) *Option {
	o.callerSkip += skip
	return o
}
