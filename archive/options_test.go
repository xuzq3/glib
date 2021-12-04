package archive

import (
	"testing"
	"time"
)

func TestOptions(t *testing.T) {
	var blockSize int64 = 1024
	var delayTime time.Duration = time.Millisecond

	opts := NewOptions().SetBlockSize(1024).SetDelayTime(delayTime)

	if opts.BlockSize != blockSize {
		t.Error("SetBlockSize failed")
	}
	if opts.DelayTime != delayTime {
		t.Error("SetSleepTime failed")
	}
}
