package klog

import (
	"io"
	"os"

	"github.com/valyala/bytebufferpool"
	"go.uber.org/atomic"
)

var _ io.Writer = (*AsyncWriter)(nil)

type AsyncWriter struct {
	w          io.Writer
	c          chan *bytebufferpool.ByteBuffer
	size       int
	closed     *atomic.Bool
	bufferPool bytebufferpool.Pool
}

func NewAsyncWriter(writer io.Writer, size int) *AsyncWriter {
	w := &AsyncWriter{
		w:      writer,
		c:      make(chan *bytebufferpool.ByteBuffer, size),
		size:   size,
		closed: atomic.NewBool(false),
	}
	go w.loop()
	return w
}

func (w *AsyncWriter) Write(p []byte) (n int, err error) {
	if w.closed.Load() {
		return 0, os.ErrClosed
	}

	// 缓冲池已满，主动丢弃
	if len(w.c) >= w.size {
		return 0, nil
	}

	// 异步情况下，不能直接使用传入的切片，需要拷贝，防止被篡改
	buf := w.bufferPool.Get()
	n, err = buf.Write(p)
	if err != nil {
		return 0, err
	}

	w.c <- buf
	return n, nil
}

func (w *AsyncWriter) loop() {
	for buf := range w.c {
		_, _ = w.w.Write(buf.Bytes())
		w.bufferPool.Put(buf)
	}
}

func (w *AsyncWriter) Close() {
	if w.closed.CompareAndSwap(false, true) {
		close(w.c)
	}
}
