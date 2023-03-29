package ws

type Group struct {
	server   *Server
	handlers []MessageHandler
}

func (g *Group) OnCmd(cmd string, handlers ...MessageHandler) {
	g.server.messageHandlers[cmd] = append(g.handlers, handlers...)
}

func (g *Group) Group(handlers ...MessageHandler) *Group {
	return &Group{
		server:   g.server,
		handlers: append(g.handlers, handlers...),
	}
}
