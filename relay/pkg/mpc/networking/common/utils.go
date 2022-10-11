package common

import "github.com/gorilla/websocket"

func NormalClose(conn *websocket.Conn) error {
	return conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
}

var (
	KeygenOperation    = []byte("keygen")
	HeaderTssID        = "X-TSS-ID"
	HeaderTssOperation = "X-TSS-Operation"
)
