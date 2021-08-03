package multicast

import (
	"context"
	"encoding/json"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/xuzq3/glib/logx"
	"github.com/xuzq3/glib/util"
	"golang.org/x/net/ipv4"
)

type Server struct {
	wildcardIP string
	groupIP    string
	port       int
	isLoopback bool
	buffSize   int

	ctx        context.Context
	cancel     context.CancelFunc
	ticker     *time.Ticker
	once       chan struct{}
	bytePool   *BytePool
	msgCtxPool *MessageContextPool
	groupAddr  *net.UDPAddr
	conn       *net.UDPConn
	pconn      *ipv4.PacketConn
	iface      *net.Interface
	joinIfaces map[string]struct{}
	routes     map[string][]Handler
}

func NewServer(wildcardIP string, groupIP string, port int, buffSize int) *Server {
	return &Server{
		wildcardIP: wildcardIP,
		groupIP:    groupIP,
		port:       port,
		buffSize:   buffSize,
		once:       make(chan struct{}, 1),
		bytePool:   NewBytePool(buffSize),
		msgCtxPool: NewMessageContextPool(),
		joinIfaces: make(map[string]struct{}),
		routes:     make(map[string][]Handler),
	}
}

func (s *Server) SetLoopback(en bool) {
	s.isLoopback = en
}

func (s *Server) Start() error {
	err := s.init()
	if err != nil {
		return err
	}
	//err = s.reloadInterface()
	//if err != nil {
	//	return err
	//}
	s.ctx, s.cancel = context.WithCancel(context.Background())
	s.ticker = time.NewTicker(time.Second * 30)
	s.once <- struct{}{}
	go s.loop()
	go s.recv()
	return nil
}

func (s *Server) Stop() {
	s.ticker.Stop()
	s.cancel()
}

func (s *Server) init() error {
	groupIP := net.ParseIP(s.groupIP)
	if isMulticastIp := net.IP.IsMulticast(groupIP); !isMulticastIp {
		err := fmt.Errorf("ip is not a multicast address: %s", s.groupIP)
		logx.Error("init multicast failed: %s", err.Error())
		return err
	}

	groupAddr := &net.UDPAddr{
		IP:   groupIP,
		Port: s.port,
	}
	s.groupAddr = groupAddr
	logx.Info("multicast group address: %s", groupAddr.String())

	wildcardIP := net.ParseIP(s.wildcardIP)
	listenAddress := &net.UDPAddr{
		IP:   wildcardIP,
		Port: s.port,
	}
	logx.Info("multicast start listen: %s", listenAddress.String())

	conn, err := net.ListenUDP("udp4", listenAddress)
	if err != nil {
		logx.Error("multicast listen udp failed: %s", err.Error())
		return err
	}
	s.conn = conn
	s.pconn = ipv4.NewPacketConn(conn)

	logx.Info("multicast listen success: %s", listenAddress.String())

	if s.isLoopback {
		err = s.pconn.SetMulticastLoopback(true)
		if err != nil {
			logx.Error("multicast SetMulticastLoopback failed: %s", err.Error())
			return err
		}
	}

	//err = s.pconn.SetControlMessage(ipv4.FlagDst, true)
	//if err != nil {
	//	logx.Error("multicast SetControlMessage failed: %s", err.Error())
	//	return err
	//}
	return nil
}

func (s *Server) reloadInterface() error {
	logx.Debug("reload multicast interface")
	ifaces := s.getValidInterfaces()

	// 检查原有的网卡能否继续使用
	if s.iface != nil {
		for _, iface := range ifaces {
			if s.iface.Name == iface.Name {
				return nil
			}
		}
	}

	// 寻找有效网卡
	for _, iface := range ifaces {
		err := s.setInterface(iface)
		if err != nil {
			continue
		}
		return nil
	}

	// 环回网卡
	_ = func() error {
		iface, err := util.GetLoopbackInterface()
		if err != nil {
			logx.Error("multicast GetLoopbackInterface failed: %s", err.Error())
			return err
		}
		if iface == nil {
			return nil
		}
		if s.iface == nil || s.iface.Name != iface.Name {
			err = s.setInterface(iface)
			if err != nil {
				return err
			}
		}
		return nil
	}()
	return nil
}

func (s *Server) setInterface(iface *net.Interface) error {
	var err error
	group := &net.UDPAddr{
		IP: net.ParseIP(s.groupIP),
	}

	if _, ok := s.joinIfaces[iface.Name]; ok {
		err = s.pconn.LeaveGroup(iface, group)
		if err != nil {
			logx.Error("multicast interface %s leave group failed: %s", iface.Name, err.Error())
			return err
		}
	}
	err = s.pconn.JoinGroup(iface, group)
	if err != nil {
		logx.Error("multicast interface %s join group failed: %s", iface.Name, err.Error())
		return err
	}
	s.joinIfaces[iface.Name] = struct{}{}

	err = s.pconn.SetMulticastInterface(iface)
	if err != nil {
		logx.Error("multicast set interface %s failed: %s", iface.Name, err.Error())
		return err
	}

	logx.Info("multicast set interface %s", iface.Name)
	s.iface = iface
	return nil
}

func (s *Server) getValidInterfaces() []*net.Interface {
	ifaces, err := net.Interfaces()
	if err != nil {
		logx.Error("multicast net.Interfaces failed: %s", err.Error())
		return nil
	}

	list := make([]*net.Interface, 0)
	for i := range ifaces {
		iface := ifaces[i]

		// 过滤无效地址
		if iface.Flags&net.FlagUp == 0 {
			continue
		}
		// 过滤环回地址
		if iface.Flags&net.FlagLoopback != 0 {
			continue
		}
		// 过滤非多播地址
		if iface.Flags&net.FlagMulticast == 0 {
			continue
		}
		// 过滤虚拟地址
		name := iface.Name
		name = strings.ToLower(name)
		if strings.Contains(name, "vmnet") {
			continue
		}
		// 过滤获取不了IP或者IP为保留地址169.254.x.x的网卡
		ip := util.GetInterfaceIPv4(&iface)
		if ip == "" || strings.HasPrefix(ip, "169.254.") {
			continue
		}
		list = append(list, &iface)
	}
	return list
}

func (s *Server) loop() {
	for {
		select {
		case <-s.ctx.Done():
			return
		case <-s.ticker.C:
			_ = s.reloadInterface()
		case <-s.once:
			_ = s.reloadInterface()
		}
	}
}

func (s *Server) Send(b []byte) error {
	logx.Debug("multicast Send msg:%s", string(b))
	if s.pconn == nil {
		err := fmt.Errorf("udp was unconnected")
		logx.Error("multicast Send failed: %s", err.Error())
		return err
	}
	_, err := s.pconn.WriteTo(b, nil, s.groupAddr)
	if err != nil {
		logx.Error("multicast Send failed: %s", err.Error())
		return err
	}
	return nil
}

func (s *Server) SendJson(v interface{}) error {
	js, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return s.Send(js)
}

func (s *Server) recv() {
	for {
		select {
		case <-s.ctx.Done():
			return
		default:
		}

		b := s.bytePool.Get()
		n, src, err := s.conn.ReadFromUDP(b)
		if err != nil {
			logx.Error("multicast ReadFrom failed: %s", err.Error())
			continue
		}
		//if !cm.Dst.IsMulticast() || !cm.Dst.Equal(s.groupAddr.IP) {
		//	logger.LogDebug("multicast read not multicast msg or unknown group msg")
		//	continue
		//}

		go s.handle(n, src, b)
	}
}

func (s *Server) handle(n int, src *net.UDPAddr, b []byte) {
	defer func() {
		s.bytePool.Put(b)
		if r := recover(); r != nil {
			logx.Error("multicast handle recover from %v", r)
		}
	}()

	msg := b[:n]
	//logx.Debug("multicast recv %s %s", src.String(), string(msg))
	ctx, err := s.parseMessage(n, src, msg)
	if err != nil {
		logx.Error("multicast parse message failed, msg:%s, err:%s", string(msg), err.Error())
		return
	}
	ctx.Next()
}

func (s *Server) parseMessage(n int, src *net.UDPAddr, b []byte) (*MessageContext, error) {
	var body MessageBody
	err := json.Unmarshal(b, &body)
	if err != nil {
		return nil, err
	}
	handlers, ok := s.routes[body.Cmd]
	if !ok {
		return nil, ErrUnknownCmd
	}

	msg := &Message{
		N:    n,
		Src:  src,
		Byte: b,
	}

	ctx := s.msgCtxPool.Get()
	ctx.reset()
	ctx.Message = msg
	ctx.Body = &body
	ctx.Server = s
	ctx.handlers = handlers
	return ctx, nil
}

func (s *Server) OnCmd(cmd string, handlers ...Handler) {
	s.routes[cmd] = handlers
}

func (s *Server) Group(handlers ...Handler) *Group {
	return &Group{
		server:   s,
		handlers: handlers,
	}
}
