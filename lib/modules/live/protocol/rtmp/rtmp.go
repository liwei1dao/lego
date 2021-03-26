package rtmp

import (
	"fmt"
	"net"
	"reflect"

	"github.com/liwei1dao/lego/lib/modules/live/av"
	"github.com/liwei1dao/lego/lib/modules/live/configure"
	"github.com/liwei1dao/lego/lib/modules/live/protocol/rtmp/core"
	"github.com/liwei1dao/lego/sys/log"
)

func NewRtmpServer(h av.Handler, getter av.GetWriter) *Server {
	return &Server{
		handler: h,
		getter:  getter,
	}
}

type Server struct {
	handler av.Handler
	getter  av.GetWriter
}

func (s *Server) Serve(listener net.Listener) (err error) {
	defer func() {
		if r := recover(); r != nil {
			log.Errorf("rtmp serve panic: ", r)
		}
	}()

	for {
		var netconn net.Conn
		netconn, err = listener.Accept()
		if err != nil {
			return
		}
		conn := core.NewConn(netconn, 4*1024)
		log.Debugf("new client, connect remote: ", conn.RemoteAddr().String(),
			"local:", conn.LocalAddr().String())
		go s.handleConn(conn)
	}
}

func (s *Server) handleConn(conn *core.Conn) error {
	if err := conn.HandshakeServer(); err != nil {
		conn.Close()
		log.Errorf("handleConn HandshakeServer err: ", err)
		return err
	}
	connServer := core.NewConnServer(conn)

	if err := connServer.ReadMsg(); err != nil {
		conn.Close()
		log.Errorf("handleConn read msg err: ", err)
		return err
	}

	appname, name, _ := connServer.GetInfo()

	if ret := configure.CheckAppName(appname); !ret {
		err := fmt.Errorf("application name=%s is not configured", appname)
		conn.Close()
		log.Errorf("CheckAppName err: ", err)
		return err
	}

	log.Debugf("handleConn: IsPublisher=%v", connServer.IsPublisher())
	if connServer.IsPublisher() {
		if configure.Config.GetBool("rtmp_noauth") {
			key, err := configure.RoomKeys.GetKey(name)
			if err != nil {
				err := fmt.Errorf("Cannot create key err=%s", err.Error())
				conn.Close()
				log.Error("GetKey err: ", err)
				return err
			}
			name = key
		}
		channel, err := configure.RoomKeys.GetChannel(name)
		if err != nil {
			err := fmt.Errorf("invalid key err=%s", err.Error())
			conn.Close()
			log.Errorf("CheckKey err: ", err)
			return err
		}
		connServer.PublishInfo.Name = channel
		if pushlist, ret := configure.GetStaticPushUrlList(appname); ret && (pushlist != nil) {
			log.Debugf("GetStaticPushUrlList: %v", pushlist)
		}
		reader := NewVirReader(connServer)
		s.handler.HandleReader(reader)
		log.Debugf("new publisher: %+v", reader.Info())

		if s.getter != nil {
			writeType := reflect.TypeOf(s.getter)
			log.Debugf("handleConn:writeType=%v", writeType)
			writer := s.getter.GetWriter(reader.Info())
			s.handler.HandleWriter(writer)
		}
		if configure.Config.GetBool("flv_archive") {
			flvWriter := new(flv.FlvDvr)
			s.handler.HandleWriter(flvWriter.GetWriter(reader.Info()))
		}
	} else {
		writer := NewVirWriter(connServer)
		log.Debugf("new player: %+v", writer.Info())
		s.handler.HandleWriter(writer)
	}

	return nil
}
