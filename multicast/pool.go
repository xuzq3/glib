package multicast

import (
	"sync"
)

type BytePool struct {
	n    int
	pool *sync.Pool
}

func NewBytePool(n int) *BytePool {
	return &BytePool{
		n: n,
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, n, n)
			},
		},
	}
}

func (p *BytePool) Get() []byte {
	b, ok := p.pool.Get().([]byte)
	if !ok {
		b = make([]byte, p.n, p.n)
	}
	return b
}

func (p *BytePool) Put(b []byte) {
	p.pool.Put(b)
}

type MessageContextPool struct {
	pool *sync.Pool
}

func NewMessageContextPool() *MessageContextPool {
	return &MessageContextPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(MessageContext)
			},
		},
	}
}

func (p *MessageContextPool) Get() *MessageContext {
	c, ok := p.pool.Get().(*MessageContext)
	if !ok {
		c = new(MessageContext)
	}
	return c
}

func (p *MessageContextPool) Put(c *MessageContext) {
	p.pool.Put(c)
}
