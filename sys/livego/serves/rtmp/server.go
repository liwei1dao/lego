package rtmp

import (
	"fmt"
	"net"
	"reflect"

	"github.com/liwei1dao/lego/sys/livego/container/flv"
	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/log"
)

func NewServer(sys core.ISys, log log.ILogger, getter core.GetWriter) (server *Server, err error) {
	server = &Server{
		sys:    sys,
		log:    log,
		getter: getter,
	}
	err = server.init()
	return
}

type Server struct {
	sys    core.ISys
	log    log.ILogger
	getter core.GetWriter
}

func (this *Server) init() (err error) {
	var (
		rtmpListen net.Listener
	)
	if rtmpListen, err = net.Listen("tcp", this.sys.GetRTMPAddr()); err != nil {
		this.log.Errorf("init err:%v", err)
		return
	}
	this.log.Infof("RTMP Listen On:%s", this.sys.GetRTMPAddr())
	go this.Serve(rtmpListen)
	return
}

func (this *Server) Serve(listener net.Listener) (err error) {
	for {
		var netconn net.Conn
		if netconn, err = listener.Accept(); err != nil {
			this.log.Errorf("Serve err: ", err)
			return
		}
		conn := core.NewConn(netconn, this.sys, this.log)
		this.log.Debugf("new client, connect remote: ", conn.RemoteAddr().String(),
			"local:", conn.LocalAddr().String())
		go this.handleConn(conn)
	}
}

func (this *Server) handleConn(conn *core.Conn) (err error) {
	if err := conn.HandshakeServer(); err != nil {
		conn.Close()
		this.log.Errorf("handleConn HandshakeServer err: ", err)
		return err
	}
	connServer := core.NewConnServer(conn)
	if err := connServer.ReadMsg(); err != nil {
		conn.Close()
		this.log.Errorf("handleConn read msg err:%v", err)
		return err
	}
	appname, name, _ := connServer.GetInfo()
	if appname != this.sys.GetAppname() {
		err := fmt.Errorf("application name=%s is not configured", appname)
		conn.Close()
		this.log.Errorf("CheckAppName err: ", err)
	}
	this.log.Debugf("handleConn: IsPublisher=%v", connServer.IsPublisher())
	if connServer.IsPublisher() {
		if this.sys.GetRTMPNoAuth() {
			key, err := this.sys.GetKey(name)
			if err != nil {
				err := fmt.Errorf("Cannot create err=%s", err.Error())
				conn.Close()
				this.log.Errorf("GetKey err:%v", err)
				return err
			}
			name = key
		}
		channel, err := this.sys.GetChannel(name)
		if err != nil {
			err := fmt.Errorf("invalid key err=%s", err.Error())
			conn.Close()
			this.log.Errorf("CheckKey err:%v", err)
			return err
		}
		connServer.PublishInfo.Name = channel
		if this.sys.GetStaticPush() != nil && len(this.sys.GetStaticPush()) > 0 {
			this.log.Debugf("GetStaticPushUrlList: %v", this.sys.GetStaticPush())
		}
		reader := NewReader(connServer)
		this.sys.HandleReader(reader)
		this.log.Debugf("new publisher: %+v", reader.Info())
		if this.getter != nil {
			writeType := reflect.TypeOf(this.getter)
			this.log.Debugf("handleConn:writeType=%v", writeType)
			writer := this.getter.GetWriter(reader.Info())
			this.sys.HandleWriter(writer)
		}
		if this.sys.GetFLVArchive() {
			flvWriter := flv.NewFlvDvr(this.sys, this.log)
			this.sys.HandleWriter(flvWriter.GetWriter(reader.Info()))
		}
	} else {
		writer := NewWriter(connServer)
		this.log.Debugf("new player: %+v", writer.Info())
		this.sys.HandleWriter(writer)
	}
	return
}
