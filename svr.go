package ws

import (
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"github.com/499689317/go-log"
)

type WHandler struct {
	mu       sync.Mutex
	wg       sync.WaitGroup
	connNum  int
	bufLen   int
	msgLen   uint32
	upgrader websocket.Upgrader
	conns    map[*websocket.Conn]struct{}
}

// WHandler implement http.Handler ServeHTTP
func (h *WHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	if r.Method != "GET" {
		// return errors.New("request error")
		return
	}

	// http upgrader to socket
	c, e := h.upgrader.Upgrade(w, r, nil)
	if e != nil {
		return
	}

	// 设置最大读取字节
	c.SetReadLimit(int64(h.msgLen))

	// 需要一个并发等待机制，确保连接退出时能清理干净
	h.wg.Add(1)
	defer h.wg.Done()

	h.mu.Lock()
	if h.conns == nil {
		h.mu.Unlock()
		c.Close()
		// return errors.New("h conns nil error")
		return
	}
	// 判断是否超出最大连接限制
	if len(h.conns) >= h.connNum {
		h.mu.Unlock()
		c.Close()
		// return errors.New("limit connection")
		return
	}
	h.conns[c] = struct{}{}
	h.mu.Unlock()

	// 封装conn，方便后期使用
	o := newWConn(c, h.bufLen, h.msgLen)
	log.Info().Msg("new client connection")
	// TODO 读取客户端消息-------

	o.Run()

	// close connection
	log.Info().Msg("track conn destroy //// close client connection")
	o.Close()
	h.mu.Lock()
	delete(h.conns, c)
	h.mu.Unlock()
}

type WServer struct {
	Addr    string
	ConnNum int
	BufLen  int
	MsgLen  uint32
	Timeout time.Duration
	ln      net.Listener
	h       *WHandler
}

func (s *WServer) Run() {

	ln, e := net.Listen("tcp", s.Addr)
	if e != nil {
		// 终止进程
		log.Fatal().Err(e).Msg("net listener error")
	}

	// 默认3000最大连接数
	if s.ConnNum <= 0 {
		s.ConnNum = 3000
	}
	// 默认写入500字节缓冲区
	if s.BufLen <= 0 {
		s.BufLen = 500
	}
	// 默认最大写入消息大小4kb
	if s.MsgLen <= 0 {
		s.MsgLen = 4096
	}
	// 默认连接超时时间30s
	if s.Timeout <= 0 {
		s.Timeout = 30 * time.Second
	}

	// WHandler
	s.ln = ln
	s.h = &WHandler{
		conns: make(map[*websocket.Conn]struct{}),
		connNum: s.ConnNum,
		bufLen:  s.BufLen,
		msgLen:  s.MsgLen,
		upgrader: websocket.Upgrader{
			HandshakeTimeout: s.Timeout,
			CheckOrigin: func(_ *http.Request) bool {
				return true
			},
		},
	}

	// WebSocket Server
	httpserver := &http.Server{
		Addr:           s.Addr,
		Handler:        s.h,
		ReadTimeout:    s.Timeout,
		WriteTimeout:   s.Timeout,
		MaxHeaderBytes: 1024,
	}

	// start server
	go httpserver.Serve(ln)
}
func (s *WServer) Close() {
	// 关闭监听
	s.ln.Close()
	// 断开所有连接，要等待正在连接的客户端连接完再断开
	// TODO 创建连接是在另一个线程中进行的
	s.h.mu.Lock()
	for c := range s.h.conns {
		c.Close()
	}
	s.h.conns = nil
	s.h.mu.Unlock()

	// 等清理完再退出
	s.h.wg.Wait()
}
