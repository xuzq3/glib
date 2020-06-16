package compressx

import (
	"testing"
	"time"
)

func TestOptions(t *testing.T) {
	var blockSize int64 = 1024
	var sleepTime time.Duration = time.Millisecond

	opts := NewOptions().SetBlockSize(1024).SetSleepTime(sleepTime)

	if opts.BlockSize != blockSize {
		t.Error("SetBlockSize failed")
	}
	if opts.SleepTime != sleepTime {
		t.Error("SetSleepTime failed")
	}
}
