package ws

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/gorilla/websocket"
	"github.com/xuzq3/glib/logx"
)

type Server struct {
	upgrader        *websocket.Upgrader
	msgCtxPool      *messageContextPool
	connCtxPool     *connContextPool
	connectHandlers []ConnHandler
	messageHandlers map[string][]MessageHandler
}

func NewServer() *Server {
	upgrader := &websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { return true },
	}
	return &Server{
		upgrader:        upgrader,
		msgCtxPool:      newMessageContextPool(),
		connCtxPool:     newConnectContextPool(),
		messageHandlers: make(map[string][]MessageHandler),
	}
}

func (s *Server) Upgrade(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		logx.WithError(err).Error("websocket upgrade failed")
		return
	}
	defer conn.Close()

	s.handleConnect(conn)
}

func (s *Server) handleConnect(conn *websocket.Conn) {
	ctx := s.connCtxPool.Get()
	ctx.reset()
	ctx.Conn = conn
	ctx.Server = s
	ctx.handlers = s.connectHandlers

	defer s.connCtxPool.Put(ctx)

	ctx.Next()

	for {
		messageType, message, err := conn.ReadMessage()
		if err != nil {
			logx.WithError(err).Error("websocket ReadMessage failed")
			return
		}
		switch messageType {
		case websocket.TextMessage, websocket.BinaryMessage:
			s.handleMessage(conn, messageType, message)
		case websocket.CloseMessage:
			return
		default:
		}
	}
}

func (s *Server) handleMessage(conn *websocket.Conn, messageType int, message []byte) {
	ctx, err := s.parseMessage(conn, messageType, message)
	if err != nil {
		logx.WithError(err).Error("websocket parseMessage failed")
		return
	}
	defer s.msgCtxPool.Put(ctx)

	ctx.Next()
}

func (s *Server) parseMessage(conn *websocket.Conn, messageType int, message []byte) (*MessageContext, error) {
	var body MessageBody
	err := json.Unmarshal(message, &body)
	if err != nil {
		return nil, err
	}
	handlers, ok := s.messageHandlers[body.Cmd]
	if !ok {
		return nil, errors.New("unknown cmd")
	}

	ctx := s.msgCtxPool.Get()
	ctx.reset()
	ctx.Conn = conn
	ctx.MessageType = messageType
	ctx.Message = message
	ctx.JsonBody = &body
	ctx.Server = s
	ctx.handlers = handlers
	return ctx, nil
}

func (s *Server) OnConnect(handlers ...ConnHandler) {
	s.connectHandlers = append(s.connectHandlers, handlers...)
}

func (s *Server) OnCmd(cmd string, handlers ...MessageHandler) {
	s.messageHandlers[cmd] = handlers
}

func (s *Server) Group(handlers ...MessageHandler) *Group {
	return &Group{
		server:   s,
		handlers: handlers,
	}
}
