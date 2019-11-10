package ws

import (
	"testing"
	"time"

	"github.com/499689317/go-log"
	zlog "github.com/rs/zerolog/log"
)

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
