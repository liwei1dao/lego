package rtmp

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/livego/container/flv"
	"github.com/liwei1dao/lego/sys/livego/core"
)

func NewServer(sys core.ISys) (server *Server, err error) {
	server = &Server{
		sys:     sys,
		streams: new(sync.Map),
	}
	err = server.init()
	return
}

type Server struct {
	sys core.ISys
	///流管理
	streams *sync.Map
}

func (this *Server) init() (err error) {
	var (
		rtmpListen net.Listener
	)
	if rtmpListen, err = net.Listen("tcp", this.sys.GetRTMPAddr()); err != nil {
		this.sys.Errorf("init err:%v", err)
		return
	}
	this.sys.Infof("RTMP Listen On:%s", this.sys.GetRTMPAddr())
	go this.Serve(rtmpListen)
	return
}

func (this *Server) Serve(listener net.Listener) (err error) {
	for {
		var netconn net.Conn
		if netconn, err = listener.Accept(); err != nil {
			this.sys.Errorf("Serve err: ", err)
			return
		}
		conn := core.NewConn(netconn, this.sys)
		this.sys.Debugf("new client, connect remote: ", conn.RemoteAddr().String(),
			"local:", conn.LocalAddr().String())
		go this.handleConn(conn)
	}
}

func (this *Server) handleConn(conn *core.Conn) (err error) {
	if err := conn.HandshakeServer(); err != nil {
		conn.Close()
		this.sys.Errorf("handleConn HandshakeServer err: ", err)
		return err
	}
	connServer := core.NewConnServer(conn)
	if err := connServer.ReadMsg(); err != nil {
		conn.Close()
		this.sys.Errorf("handleConn read msg err:%v", err)
		return err
	}
	appname, name, _ := connServer.GetInfo()
	if appname != this.sys.GetAppname() {
		err := fmt.Errorf("application name=%s is not configured", appname)
		conn.Close()
		this.sys.Errorf("CheckAppName err: ", err)
	}
	this.sys.Debugf("handleConn: IsPublisher=%v", connServer.IsPublisher())
	if connServer.IsPublisher() {
		if this.sys.GetRTMPNoAuth() {
			key, err := this.sys.GetKey(name)
			if err != nil {
				err := fmt.Errorf("Cannot create err=%s", err.Error())
				conn.Close()
				this.sys.Errorf("GetKey err:%v", err)
				return err
			}
			name = key
		}
		channel, err := this.sys.GetChannel(name)
		if err != nil {
			err := fmt.Errorf("invalid key err=%s", err.Error())
			conn.Close()
			this.sys.Errorf("CheckKey err:%v", err)
			return err
		}
		connServer.PublishInfo.Name = channel
		if this.sys.GetStaticPush() != nil && len(this.sys.GetStaticPush()) > 0 {
			this.sys.Debugf("GetStaticPushUrlList: %v", this.sys.GetStaticPush())
		}
		reader := NewReader(connServer)
		this.HandleReader(reader)
		this.sys.Debugf("new publisher: %+v", reader.Info())
		if this.sys.GetFLVArchive() {
			flvWriter := flv.NewFlvDvr(this.sys)
			this.HandleWriter(flvWriter.GetWriter(reader.Info()))
		}
	} else {
		writer := NewWriter(connServer)
		this.sys.Debugf("new player: %+v", writer.Info())
		this.HandleWriter(writer)
	}
	return
}

///流管理***********************************************************************
func (this *Server) GetStreams() *sync.Map {
	return this.streams
}

//监测存活
func (this *Server) CheckAlive() {
	for {
		<-time.After(5 * time.Second)
		this.streams.Range(func(key, val interface{}) bool {
			v := val.(*core.Stream)
			if v.CheckAlive() == 0 {
				this.streams.Delete(key)
			}
			return true
		})
	}
}

//读处理
func (this *Server) HandleReader(r core.ReadCloser) {
	info := r.Info()
	this.sys.Debugf("HandleReader: info[%v]", info)
	var stm *core.Stream
	i, ok := this.streams.Load(info.Key)
	if stm, ok = i.(*core.Stream); ok {
		stm.TransStop()
		id := stm.ID()
		if id != core.EmptyID && id != info.UID {
			ns := core.NewStream()
			stm.Copy(ns)
			stm = ns
			this.streams.Store(info.Key, ns)
		}
	} else {
		stm = core.NewStream()
		this.streams.Store(info.Key, stm)
		stm.SetInfo(info)
	}
	stm.AddReader(r)
}

//写处理
func (this *Server) HandleWriter(w core.WriteCloser) {
	info := w.Info()
	this.sys.Debugf("HandleWriter: info[%v]", info)

	var s *core.Stream
	item, ok := this.streams.Load(info.Key)
	if !ok {
		this.sys.Debugf("HandleWriter: not found create new info[%v]", info)
		s = core.NewStream()
		this.streams.Store(info.Key, s)
		s.SetInfo(info)
	} else {
		s = item.(*core.Stream)
		s.AddWriter(w)
	}
}
