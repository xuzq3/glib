package ws

import (
	"context"
	"github.com/gorilla/websocket"
)

type ConnHandler func(c *ConnContext)

type ConnContext struct {
	Conn     *websocket.Conn
	Server   *Server
	Error    error
	handlers []ConnHandler
	index    int8
	ctx      context.Context
}

func (c *ConnContext) reset() {
	c.Conn = nil
	c.Server = nil
	c.Error = nil
	c.handlers = nil
	c.index = -1
	c.ctx = nil
}

func (c *ConnContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *ConnContext) IsAborted() bool {
	return c.index >= abortIndex
}

func (c *ConnContext) Abort() {
	c.index = abortIndex
}

func (c *ConnContext) AbortWithError(err error) {
	c.index = abortIndex
	c.Error = err
}

func (c *ConnContext) Context() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

func (c *ConnContext) WithContext(ctx context.Context) *ConnContext {
	c.ctx = ctx
	return c
}
