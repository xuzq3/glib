package compressx

import (
	"time"
)

const defaultBlockSize = 4 * 1024

type Options struct {
	BlockSize int64
	SleepTime time.Duration
}

func NewOptions() *Options {
	return &Options{
		BlockSize: defaultBlockSize,
		SleepTime: 0,
	}
}

func (o *Options) SetBlockSize(size int64) *Options {
	o.BlockSize = size
	return o
}

func (o *Options) SetSleepTime(d time.Duration) *Options {
	o.SleepTime = d
	return o
}
