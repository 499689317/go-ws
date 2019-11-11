package ws

import (
	"testing"
	"time"

	"github.com/499689317/go-log"
	zlog "github.com/rs/zerolog/log"
)

type conf struct {
	WsListenAddr string
	WsConnNum    int
	WsBufLen     int
	WsMsgLen     uint32
	WsTimeout    time.Duration
}

func (c *conf) WSListenAddr() string {
	return c.WsListenAddr
}
func (c *conf) WSConnNum() int {
	return c.WsConnNum
}
func (c *conf) WSBufLen() int {
	return c.WsBufLen
}
func (c *conf) WSMsgLen() uint32 {
	return c.WsMsgLen
}
func (c *conf) WSTimeout() time.Duration {
	return c.WsTimeout
}

func TestNewWebSocketServer(t *testing.T) {

	c := &conf{
		WsListenAddr: ":8008",
		WsConnNum:    1000,
		WsBufLen:     300,
		WsMsgLen:     300,
		WsTimeout:    30,
	}

	s := NewWebSocket(c)
	s.Run()
}

func TestWebSocketServer(t *testing.T) {

	zlog.Logger = zlog.With().Caller().Logger()
	log.Init()
	log.SetLogLevel(1)

	s := new(WebSocket)
	s.Addr = ":8888"
	s.ConnNum = 1000
	s.BufLen = 100
	s.MsgLen = 1000
	s.Timeout = 10 * time.Second
	s.Run()
}
