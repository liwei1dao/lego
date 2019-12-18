package gate

import (
	"bytes"
	"github.com/gorilla/websocket"
	"net"
	"time"
)

func NewWsConn(conn *websocket.Conn) IConn {
	c := &WsConn{}
	c.Init(conn)
	return c
}

type WsConn struct {
	conn      *websocket.Conn
	closeFlag bool
	buffer    bytes.Buffer
}

func (this *WsConn) Init(conn *websocket.Conn) {
	this.conn = conn
	this.closeFlag = false
}
func (this *WsConn) Write(b []byte) (n int, err error) {
	err = this.conn.WriteMessage(websocket.BinaryMessage, b)
	if err != nil {
		return 0, err
	}
	return len(b), nil
}
func (this *WsConn) Read(p []byte) (int, error) {
	_, b, err := this.conn.ReadMessage()
	if err != nil {
		return 0, err
	} else {
		this.buffer.Write(b)
	}
	return this.buffer.Read(p)
}
func (this *WsConn) Close() {
	this.closeFlag = true
	this.conn.Close()
}

func (this *WsConn) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}
func (this *WsConn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}
func (this *WsConn) SetDeadline(t time.Time) error {
	err := this.conn.SetWriteDeadline(t)
	if err != nil {
		return err
	}
	return this.conn.SetWriteDeadline(t)
}
func (this *WsConn) SetReadDeadline(t time.Time) error {
	return this.conn.SetReadDeadline(t)
}
func (this *WsConn) SetWriteDeadline(t time.Time) error {
	return this.conn.SetWriteDeadline(t)
}
