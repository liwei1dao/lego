package gate

import (
	"bytes"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

func NewWsConn(conn *websocket.Conn, heartbeat time.Duration) IConn {
	c := &WsConn{}
	c.Init(conn, heartbeat)
	return c
}

type WsConn struct {
	conn      *websocket.Conn
	heartbeat time.Duration
	closeFlag bool
	buffer    bytes.Buffer
}

func (this *WsConn) Init(conn *websocket.Conn, heartbeat time.Duration) {
	this.conn = conn
	this.heartbeat = heartbeat
	this.closeFlag = false
	this.conn.SetReadDeadline(time.Now().Add(this.heartbeat))
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
	this.conn.SetReadDeadline(time.Now().Add(this.heartbeat))
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
