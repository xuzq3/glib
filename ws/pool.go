package ws

import "sync"

type bytePool struct {
	n    int
	pool *sync.Pool
}

func newBytePool(n int) *bytePool {
	return &bytePool{
		n: n,
		pool: &sync.Pool{
			New: func() interface{} {
				return make([]byte, n, n)
			},
		},
	}
}

func (p *bytePool) Get() []byte {
	b, ok := p.pool.Get().([]byte)
	if !ok {
		b = make([]byte, p.n, p.n)
	}
	return b
}

func (p *bytePool) Put(b []byte) {
	p.pool.Put(b)
}

type connContextPool struct {
	pool *sync.Pool
}

func newConnectContextPool() *connContextPool {
	return &connContextPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(ConnContext)
			},
		},
	}
}

func (p *connContextPool) Get() *ConnContext {
	c, ok := p.pool.Get().(*ConnContext)
	if !ok {
		c = new(ConnContext)
	}
	return c
}

func (p *connContextPool) Put(c *ConnContext) {
	p.pool.Put(c)
}

type messageContextPool struct {
	pool *sync.Pool
}

func newMessageContextPool() *messageContextPool {
	return &messageContextPool{
		pool: &sync.Pool{
			New: func() interface{} {
				return new(MessageContext)
			},
		},
	}
}

func (p *messageContextPool) Get() *MessageContext {
	c, ok := p.pool.Get().(*MessageContext)
	if !ok {
		c = new(MessageContext)
	}
	return c
}

func (p *messageContextPool) Put(c *MessageContext) {
	p.pool.Put(c)
}
