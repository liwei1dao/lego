package livego

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/livego/stream"
	"github.com/liwei1dao/lego/sys/log"
)

type LiveGo struct {
	options Options
	log     log.ILog
	streams *sync.Map //流管理
}

func (this *LiveGo) Serve(listener net.Listener) (err error) {
	for {
		var netconn net.Conn
		if netconn, err = listener.Accept(); err != nil {
			this.Errorf("Serve err: ", err)
			return
		}
		conn := core.NewConn(netconn, this)
		this.Debugf("new client, connect remote: ", conn.RemoteAddr().String(),
			"local:", conn.LocalAddr().String())
		go this.handleConn(conn)
	}
}

func (this *LiveGo) handleConn(conn *core.Conn) (err error) {
	if err := conn.HandshakeServer(); err != nil {
		conn.Close()
		this.Errorf("handleConn HandshakeServer err: ", err)
		return err
	}
	connServer := core.NewConnServer(conn)
	if err := connServer.ReadMsg(); err != nil {
		conn.Close()
		this.Errorf("handleConn read msg err:%v", err)
		return err
	}
	appname, name, _ := connServer.GetInfo()
	if appname != this.options.Appname {
		err := fmt.Errorf("application name=%s is not configured", appname)
		conn.Close()
		this.Errorf("CheckAppName err: ", err)
	}
	this.Debugf("handleConn: IsPublisher=%v", connServer.IsPublisher())
	if connServer.IsPublisher() {
		if this.options.RTMPNoAuth {
			key, err := this.GetKey(name)
			if err != nil {
				err := fmt.Errorf("Cannot create err=%s", err.Error())
				conn.Close()
				this.Errorf("GetKey err:%v", err)
				return err
			}
			name = key
		}
		channel, err := this.GetChannel(name)
		if err != nil {
			err := fmt.Errorf("invalid key err=%s", err.Error())
			conn.Close()
			this.Errorf("CheckKey err:%v", err)
			return err
		}
		connServer.PublishInfo.Name = channel
		if this.options.StaticPush != nil && len(this.options.StaticPush) > 0 {
			this.Debugf("GetStaticPushUrlList: %v", this.options.StaticPush)
		}
		reader := stream.NewReader(connServer)
		this.handler.HandleReader(reader)
		this.Debugf("new publisher: %+v", reader.Info())
	} else {
		writer := stream.NewWriter(connServer)
		this.Debugf("new player: %+v", writer.Info())
		this.handler.HandleWriter(writer)
	}
	return
}

///流管理***********************************************************************
//监测存活
func (this *LiveGo) CheckAlive() {
	for {
		<-time.After(5 * time.Second)
		this.streams.Range(func(key, val interface{}) bool {
			v := val.(*Stream)
			if v.CheckAlive() == 0 {
				this.streams.Delete(key)
			}
			return true
		})
	}
}

//读处理
func (this *LiveGo) HandleReader(r core.ReadCloser) {
	info := r.Info()
	this.Debugf("HandleReader: info[%v]", info)
	var stm *Stream
	i, ok := this.streams.Load(info.Key)
	if stm, ok = i.(*Stream); ok {
		stm.TransStop()
		id := stm.ID()
		if id != EmptyID && id != info.UID {
			ns := NewStream()
			stm.Copy(ns)
			stm = ns
			this.streams.Store(info.Key, ns)
		}
	} else {
		stm = NewStream()
		this.streams.Store(info.Key, stm)
		stm.info = info
	}
	stm.AddReader(r)
}

///参数列表***********************************************************************
func (this *LiveGo) GetAppname() string {
	return this.options.Appname
}
func (this *LiveGo) GetHls() bool {
	return this.options.Hls
}
func (this *LiveGo) GetFlv() bool {
	return this.options.Flv
}
func (this *LiveGo) GetApi() bool {
	return this.options.Api
}
func (this *LiveGo) GetStaticPush() []string {
	return this.options.StaticPush
}
func (this *LiveGo) GetRTMPNoAuth() bool {
	return this.options.RTMPNoAuth
}
func (this *LiveGo) GetEnableTLSVerify() bool {
	return this.options.EnableTLSVerify
}
func (this *LiveGo) GetTimeout() int {
	return this.options.Timeout
}
func (this *LiveGo) GetConnBuffSzie() int {
	return this.options.ConnBuffSzie
}
func (this *LiveGo) GetDebug() bool {
	return this.options.Debug
}

///日志***********************************************************************
func (this *LiveGo) Debugf(format string, a ...interface{}) {
	if this.options.Debug {
		this.log.Debugf("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Infof(format string, a ...interface{}) {
	if this.options.Debug {
		this.log.Infof("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Warnf(format string, a ...interface{}) {
	if this.options.Debug {
		this.log.Warnf("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Errorf(format string, a ...interface{}) {
	if this.options.Debug {
		this.log.Errorf("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Panicf(format string, a ...interface{}) {
	if this.options.Debug {
		this.log.Panicf("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Fatalf(format string, a ...interface{}) {
	if this.options.Debug {
		this.log.Fatalf("[SYS LiveGo] "+format, a)
	}
}

///房间***********************************************************************
func (this *LiveGo) SetKey(channel string) (key string, err error) {
	return
}
func (this *LiveGo) GetKey(channel string) (newKey string, err error) {
	return
}
func (this *LiveGo) GetChannel(key string) (channel string, err error) {
	return
}
func (this *LiveGo) DeleteChannel(channel string) (ok bool) {
	return
}
func (this *LiveGo) DeleteKey(key string) (ok bool) {
	return
}

///静态推送管理***********************************************************************
func (this *LiveGo) GetAndCreateStaticPushObject(rtmpurl string) (pushobj *core.StaticPush) {
	return
}

func (this *LiveGo) GetStaticPushObject(rtmpurl string) (pushobj *core.StaticPush, err error) {
	return
}

func (this *LiveGo) ReleaseStaticPushObject(rtmpurl string) {
	return
}
