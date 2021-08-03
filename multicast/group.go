package multicast

type Group struct {
	server   *Server
	handlers []Handler
}

func (g *Group) OnCmd(cmd string, handlers ...Handler) {
	g.server.routes[cmd] = append(g.handlers, handlers...)
}

func (g *Group) Group(handlers ...Handler) *Group {
	return &Group{
		server:   g.server,
		handlers: append(g.handlers, handlers...),
	}
}
