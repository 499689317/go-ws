package ws

import (
	"errors"
	"net"
	"sync"
	//"fmt"
	"github.com/gorilla/websocket"
	"github.com/499689317/go-log"
	"github.com/499689317/go-coding"
)

type WConn struct {
	mu        sync.Mutex
	conn      *websocket.Conn
	writeChan chan []byte
	msgLen    uint32
}

func newWConn(conn *websocket.Conn, bufLen int, msgLen uint32) *WConn {
	wConn := new(WConn)
	wConn.conn = conn
	wConn.writeChan = make(chan []byte, bufLen)
	wConn.msgLen = msgLen
	go func() {
		// TODO 使用range来读取channel时，range可以感知channel的关闭，当channel关闭时，range就会结束并退出for循环
		for x := range wConn.writeChan {
			if x == nil {
				log.Info().Msg("exit conn")
				break
			}
			e := conn.WriteMessage(websocket.BinaryMessage, x)
			if e != nil {
				break
			}
		}
		// TODO 退出连接，清空writeChan后关闭连接
		wConn.doDestroy()

	}()

	return wConn
}

func (c *WConn) Run() {
	ecode := coding.NewCoding()
	for {
		x, e := c.Read()
		if e != nil {
			log.Error().Err(e).Msg("read byte error")
			break
		}
		l, e := ecode.Redecode(x)
		for _, v := range l {
			m := v.(*coding.Message)
			log.Info().Uint16("Size", m.Size).Uint16("Version", m.Version).Uint16("Server", m.Server).Uint32("TimeStamp", m.TimeStamp).Msg("byte decode")
		}
	}
}

func (c *WConn) LocalAddr() net.Addr {
	return c.conn.LocalAddr()
}
func (c *WConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}

// 判断缓冲区是否已满----被阻塞了
func (c *WConn) isWriteChanFull() bool {
	return len(c.writeChan) == cap(c.writeChan)
}

// 读取数据  TODO 是否需要线程安全？？？
func (c *WConn) Read() ([]byte, error) {
	_, x, e := c.conn.ReadMessage()
	return x, e
}

// 写入数据  TODO 写入数据一定要保证线程安全，写入过程确保其它线程不改变x
func (c *WConn) doWrite(x []byte) {

	if c.isWriteChanFull() {
		c.doDestroy()
		return
	}
	c.writeChan <- x
}
func (c *WConn) Write(x ...[]byte) error {

	c.mu.Lock()
	defer c.mu.Unlock()

	// TODO 判断连接是否关闭or是否正在关闭

	var l uint32
	for i := 0; i < len(x); i++ {
		l += uint32(len(x[i]))
	}

	// 消息包大小检测
	if l > c.msgLen {
		return errors.New("l > c.msgLen")
	} else if l <= 0 {
		return errors.New("l <= 0")
	}

	y := make([]byte, l)
	n := 0
	for i := 0; i < len(x); i++ {
		copy(y[n:], x[i])
		n += len(x[i])
	}

	c.doWrite(y)

	return nil
}

// close connection TODO this is soft quite
func (c *WConn) Close() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.doWrite(nil)
}

// destroy connection TODO this is hard quite
func (c *WConn) doDestroy() {
	c.conn.UnderlyingConn().(*net.TCPConn).SetLinger(0)
	c.conn.Close()
	// 关闭writeChan
	// TODO channel不需要通过close来释放资源，只要没有goroutine持有channel,相关资源会自动释放
	close(c.writeChan)
}
func (c *WConn) Destroy() {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.doDestroy()
}
