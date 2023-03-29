package ws

import (
	"context"
	"encoding/json"
	"github.com/gin-gonic/gin/binding"
	"github.com/gorilla/websocket"
	"github.com/pkg/errors"
)

type MessageHandler func(ctx *MessageContext)

type MessageBody struct {
	Cmd   string          `json:"cmd"`
	Seqno string          `json:"seqno,omitempty"`
	Data  json.RawMessage `json:"data,omitempty"`
}

type SendBody struct {
	Cmd   string      `json:"cmd"`
	Seqno string      `json:"seqno,omitempty"`
	Data  interface{} `json:"data,omitempty"`
}

type RespBody struct {
	Cmd   string      `json:"cmd"`
	Seqno string      `json:"seqno,omitempty"`
	Code  int         `json:"code"`
	Msg   string      `json:"msg"`
	Data  interface{} `json:"data,omitempty"`
}

type MessageContext struct {
	Conn        *websocket.Conn
	MessageType int
	Message     []byte
	JsonBody    *MessageBody
	Server      *Server
	Error       error
	handlers    []MessageHandler
	index       int8
	ctx         context.Context
}

func (c *MessageContext) reset() {
	c.Conn = nil
	c.MessageType = 0
	c.Message = nil
	c.JsonBody = nil
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
	err := binding.JSON.BindBody(c.JsonBody.Data, v)
	if err != nil {
		err = errors.Wrap(err, "invalid params")
		return err
	}
	return nil
}

func (c *MessageContext) WriteMessage(data []byte) error {
	return c.Conn.WriteMessage(websocket.TextMessage, data)
}
