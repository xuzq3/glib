package compress

import (
	"time"
)

const defaultBlockSize = 4 * 1024

type Options struct {
	BlockSize int64
	DelayTime time.Duration
}

func NewOptions() *Options {
	return &Options{
		BlockSize: defaultBlockSize,
		DelayTime: 0,
	}
}

func (o *Options) SetBlockSize(size int64) *Options {
	o.BlockSize = size
	return o
}

func (o *Options) SetDelayTime(d time.Duration) *Options {
	o.DelayTime = d
	return o
}
