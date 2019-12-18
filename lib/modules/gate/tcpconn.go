package gate

import (
	"net"
	"sync"
	"time"
)

func NewTcpConn(conn net.Conn) IConn {
	c := &TcpConn{}
	c.Init(conn)
	return c
}

type TcpConn struct {
	sync.Mutex
	conn      net.Conn
	closeFlag bool
}

func (this *TcpConn) Init(conn net.Conn) {
	this.conn = conn
	this.closeFlag = false
}
func (this *TcpConn) Write(b []byte) (n int, err error) {
	this.Lock()
	defer this.Unlock()
	if this.closeFlag || b == nil {
		return
	}
	return this.conn.Write(b)
}
func (this *TcpConn) Read(b []byte) (int, error) {
	return this.conn.Read(b)
}
func (this *TcpConn) Close() {
	this.closeFlag = true
	this.conn.(*net.TCPConn).SetLinger(0)
	this.conn.Close()
}

func (this *TcpConn) LocalAddr() net.Addr {
	return this.conn.LocalAddr()
}
func (this *TcpConn) RemoteAddr() net.Addr {
	return this.conn.RemoteAddr()
}
func (this *TcpConn) SetDeadline(t time.Time) error {
	return this.conn.SetDeadline(t)
}
func (this *TcpConn) SetReadDeadline(t time.Time) error {
	return this.conn.SetReadDeadline(t)
}
func (this *TcpConn) SetWriteDeadline(t time.Time) error {
	return this.conn.SetWriteDeadline(t)
}
