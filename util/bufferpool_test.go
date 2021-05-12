package util

import (
	"bytes"
	"testing"
)

func TestBufferPool_GET_PUT(t *testing.T) {
	pool := NewBufferPool()
	pool2 := NewBufferPool()

	for i := 0; i < 100; i++ {
		b := pool.Get()
		if b == nil {
			t.Fatal("BufferPool Get result is nil")
		}

		buf, ok := b.(*bytes.Buffer)
		if !ok {
			t.Fatal("BufferPool Get result is not buffer")
		}

		buf.WriteString("test")
		pool2.Put(buf)
	}

	for i := 0; i < 100; i++ {
		buf := pool2.GetBuffer()
		if buf == nil {
			t.Fatal("BufferPool GetBuffer result is nil")
		}

		if buf.String() != "test" {
			t.Fatal("BufferPool not use cache")
		}
	}
}
