# go-ws


+ install go-ws

`go get -u github.com/499689317/go-ws`

+ how to use it

````

import (
	"time"

	"github.com/499689317/go-log"
	"github.com/499689317/go-ws"
	zlog "github.com/rs/zerolog/log"
)

type conf struct {
	WsListenAddr   string
	WsConnNum      int
	WsBufLen       int
	WsMsgLen       uint32
	WsTimeout      time.Duration
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

func main() {

	zlog.Logger = zlog.With().Caller().Logger()
	log.Init()
	log.SetLogLevel(1)

	// conf
	c := &conf{
		WsListenAddr:   ":8070",
		WsConnNum:      3000,
		WsBufLen:       300,
		WsMsgLen:       300,
		WsTimeout:      30,
	}

	// ws
	s := ws.NewWebSocket(c)
	s.Run()
}

````