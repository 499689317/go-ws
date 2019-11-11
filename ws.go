package ws

import (
	"errors"
	"time"

	"github.com/499689317/go-log"
)

type Configurable interface {
	WSListenAddr() string
	WSConnNum() int
	WSBufLen() int
	WSMsgLen() uint32
	WSTimeout() time.Duration
}

type WebSocket struct {
	Addr     string
	ConnNum  int
	BufLen   int
	MsgLen   uint32
	Timeout  time.Duration
	killChan chan bool
	config   Configurable
}

func NewWebSocket(c Configurable) *WebSocket {

	w := &WebSocket{
		Addr:    c.WSListenAddr(),
		ConnNum: c.WSConnNum(),
		BufLen:  c.WSBufLen(),
		MsgLen:  c.WSMsgLen(),
		Timeout: c.WSTimeout(),
	}

	w.config = c

	return w
}

func (w *WebSocket) Run() error {

	var s *WServer
	if w.Addr == "" {
		return errors.New("websocket addr error")
	}

	w.killChan = make(chan bool)

	s = new(WServer)
	s.Addr = w.Addr
	s.ConnNum = w.ConnNum
	s.BufLen = w.BufLen
	s.MsgLen = w.MsgLen
	s.Timeout = w.Timeout

	s.Run()
	log.Info().Msg("start ws server ok")
	<-w.killChan

	s.Close()
	return nil
}

func (w *WebSocket) Close() {
	w.killChan <- true
	close(w.killChan)
}

// // 对外代理层
// type proxy struct {
// 	conn *WConn
// }

// func newProxy(conn *WConn) {
// 	return &proxy{
// 		conn: conn
// 	}
// }

// func (p *proxy) Run() {
// 	for {

// 		m, e := p.conn.Read()
// 		if e != nil {
// 			log.Debug("e: %v", e)
// 			break
// 		}

// 		/**
// 		 * TODO 对消息编解码后派发到业务层
// 		 */

// 	}
// }
