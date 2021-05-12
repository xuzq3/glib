package internal

import (
	"testing"
)

func TestBufferPool_Get(t *testing.T) {
	pool := newBufferPool()
	buf := pool.Get()
	if buf == nil {
		t.Errorf("BufferPool Get result is nil")
	}
}
