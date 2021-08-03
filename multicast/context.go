package multicast

import (
	"context"
	"encoding/json"
	"math"
	"net"

	"github.com/gin-gonic/gin/binding"
	"github.com/pkg/errors"
	"github.com/xuzq3/glib/logx"
)

const abortIndex int8 = math.MaxInt8 / 2

type Handler func(ctx *MessageContext)

type Body struct {
	Cmd   string      `json:"cmd"`
	Seqno string      `json:"seqno,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type Message struct {
	N    int
	Src  *net.UDPAddr
	Byte []byte
}

type MessageBody struct {
	Cmd   string          `json:"cmd"`
	Seqno string          `json:"seqno,omitempty"`
	Data  json.RawMessage `json:"data,omitempty"`
}

type MessageContext struct {
	Message  *Message
	Body     *MessageBody
	Server   *Server
	Error    error
	handlers []Handler
	index    int8
	ctx      context.Context
}

func (c *MessageContext) reset() {
	c.Message = nil
	c.Body = nil
	c.Server = nil
	c.Error = nil
	c.handlers = nil
	c.index = -1
	c.ctx = nil
}

func (c *MessageContext) Next() {
	c.index++
	for c.index < int8(len(c.handlers)) {
		c.handlers[c.index](c)
		c.index++
	}
}

func (c *MessageContext) IsAborted() bool {
	return c.index >= abortIndex
}

func (c *MessageContext) Abort() {
	c.index = abortIndex
}

func (c *MessageContext) AbortWithError(err error) {
	c.index = abortIndex
	c.Error = err
}

func (c *MessageContext) Context() context.Context {
	if c.ctx == nil {
		return context.Background()
	}
	return c.ctx
}

func (c *MessageContext) WithContext(ctx context.Context) *MessageContext {
	c.ctx = ctx
	return c
}

func (c *MessageContext) ShouldBindJson(v interface{}) error {
	err := binding.JSON.BindBody(c.Body.Data, v)
	if err != nil {
		err = errors.Wrap(err, "invalid params")
		return err
	}
	return nil
}

func (c *MessageContext) LogError(err error) {
	logx.Error("handle message failed, cmd:%s, seq:%s, err:%s", c.Body.Cmd, c.Body.Seqno, err)
	c.Error = err
}
