package rtmp

import (
	"fmt"
	"net"
	"reflect"

	"github.com/liwei1dao/lego"
	"github.com/liwei1dao/lego/sys/livego/container/flv"
	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/livego/packet"
	"github.com/liwei1dao/lego/sys/log"
)

const (
	maxQueueNum           = 1024
	SAVE_STATICS_INTERVAL = 5000
)

func NewRtmpServer(h packet.Handler, getter packet.GetWriter) (s *Rtmp_Server, err error) {
	s = &Rtmp_Server{
		handler: h,
		getter:  getter,
	}
	return
}

type Rtmp_Server struct {
	engine  core.IEngine
	handler packet.Handler
	getter  packet.GetWriter
}

func (this *Rtmp_Server) Serve(listener net.Listener) (err error) {
	defer lego.Recover("[SYS LiveGo] Serve:")
	for {
		var netconn net.Conn
		if netconn, err = listener.Accept(); err != nil {
			return
		}
		conn := core.NewConn(netconn, 4*1024)
		log.Debugf("[SYS LiveGo] new client, connect remote: ", conn.RemoteAddr().String(),
			"local:", conn.LocalAddr().String())
		go this.handleConn(conn)
	}
}

func (this *Rtmp_Server) handleConn(conn *core.Conn) (err error) {
	if err = conn.HandshakeServer(); err != nil {
		conn.Close()
		log.Errorf("[SYS LiveGo] handleConn HandshakeServer err: ", err)
		return err
	}
	connServer := core.NewConnServer(conn)
	if err = connServer.ReadMsg(); err != nil {
		conn.Close()
		log.Errorf("[SYS LiveGo] handleConn read msg err:%v", err)
		return err
	}

	appname, name, _ := connServer.GetInfo()
	// 校验数据流管道
	if ret := this.engine.CheckAppName(appname); !ret {
		err := fmt.Errorf("application name=%s is not configured", appname)
		conn.Close()
		log.Errorf("[SYS LiveGo] CheckAppName err: %v", err)
		return err
	}
	log.Debugf("[SYS LiveGo] handleConn: IsPublisher=%v", connServer.IsPublisher())
	if connServer.IsPublisher() {
		if this.engine.RTMPNoAuth() {
			key, err := this.engine.GetKey(name)
			if err != nil {
				err := fmt.Errorf("Cannot create key err=%s", err.Error())
				conn.Close()
				log.Errorf("[SYS LiveGo] GetKey err:%v", err)
				return err
			}
			name = key
		}
		channel, err := this.engine.GetChannel(name)
		if err != nil {
			err := fmt.Errorf("invalid key err=%s", err.Error())
			conn.Close()
			log.Errorf("[SYS LiveGo] CheckKey err:%v", err)
			return err
		}
		connServer.PublishInfo.Name = channel
		if pushlist, ret := this.engine.GetStaticPushUrlList(appname); ret && (pushlist != nil) {
			log.Debugf("[SYS LiveGo] GetStaticPushUrlList: %v", pushlist)
		}
		reader := NewVirReader(connServer, this.engine.Timeout())
		this.handler.HandleReader(reader)
		log.Debugf("[SYS LiveGo] new publisher: %+v", reader.Info())

		if this.getter != nil {
			writeType := reflect.TypeOf(this.getter)
			log.Debugf("[SYS LiveGo] handleConn:writeType=%v", writeType)
			writer := this.getter.GetWriter(reader.Info())
			this.handler.HandleWriter(writer)
		}
		if this.engine.FLVArchive() {
			flvWriter := new(flv.FlvDvr)
			this.handler.HandleWriter(flvWriter.GetWriter(this.engine.FLVDir(), reader.Info()))
		}
	} else {
		writer := NewVirWriter(connServer, this.engine.Timeout())
		log.Debugf("[SYS LiveGo] new player: %+v", writer.Info())
		this.handler.HandleWriter(writer)
	}

	return
}

type GetInFo interface {
	GetInfo() (string, string, string)
}

type StreamReadWriteCloser interface {
	GetInFo
	Close(error)
	Write(core.ChunkStream) error
	Read(c *core.ChunkStream) error
}
