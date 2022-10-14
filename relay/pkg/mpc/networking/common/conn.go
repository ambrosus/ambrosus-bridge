package common

import (
	"sync"

	"github.com/gorilla/websocket"
)

type Conn struct {
	Conn *websocket.Conn
	sync.Mutex
}

func (c *Conn) Write(msg []byte) error {
	c.Lock()
	defer c.Unlock()
	return c.Conn.WriteMessage(websocket.BinaryMessage, msg)
}

func (c *Conn) Read() ([]byte, error) {
	_, msg, err := c.Conn.ReadMessage()
	return msg, err
}

func (c *Conn) NormalClose() error {
	return c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

func (c *Conn) ErrorClose(protocolError int, err string) error {
	return c.Conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(protocolError, err))
}

func (c *Conn) Close(err error) error {
	if err != nil {
		return c.ErrorClose(websocket.CloseProtocolError, err.Error())
	}
	return c.NormalClose()
}
