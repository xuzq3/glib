package util

import (
	"bytes"
	"sync"
)

type BufferPool struct {
	pool *sync.Pool
}

func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
	}
}

func (p *BufferPool) Get() interface{} {
	return p.pool.Get()
}

func (p *BufferPool) Put(buf interface{}) {
	p.pool.Put(buf)
}

func (p *BufferPool) GetBuffer() *bytes.Buffer {
	buf, ok := p.pool.Get().(*bytes.Buffer)
	if !ok {
		buf = new(bytes.Buffer)
	}
	return buf
}

func (p *BufferPool) PutBuffer(buf *bytes.Buffer) {
	p.pool.Put(buf)
}
