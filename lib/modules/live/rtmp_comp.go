package live

import (
	"fmt"
	"net"
	"reflect"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib/modules/live/container/flv"
	"github.com/liwei1dao/lego/lib/modules/live/liveconn"
	"github.com/liwei1dao/lego/sys/log"
)

type RtmpComp struct {
	cbase.ModuleCompBase
	options IOptions
	module  ILive
	listen  net.Listener
}

func (this *RtmpComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.options = options.(IOptions)
	this.module = module.(ILive)
	return
}

func (this *RtmpComp) Start() (err error) {
	err = this.ModuleCompBase.Start()
	if this.listen, err = net.Listen("tcp", this.options.GetRtmpAddr()); err == nil {
		go this.run()
	}
	return
}

func (this *RtmpComp) run() {
	var (
		err error
	)
	for {
		var netconn net.Conn
		netconn, err = this.listen.Accept()
		if err != nil {
			break
		}
		conn := liveconn.NewConn(netconn, 4*1024)
		log.Debugf("new client, connect remote: ", conn.RemoteAddr().String(),
			"local:", conn.LocalAddr().String())
		go this.handleConn(conn)
	}
	if err != nil {
		log.Fatalf("rtmp is err:%v", err)
	}
}

func (this *RtmpComp) handleConn(conn *liveconn.Conn) (err error) {
	if err = conn.HandshakeServer(); err != nil {
		conn.Close()
		log.Errorf("handleConn HandshakeServer err: ", err)
		return err
	}
	connServer := liveconn.NewConnServer(conn)

	if err := connServer.ReadMsg(); err != nil {
		conn.Close()
		log.Errorf("handleConn read msg err: ", err)
		return err
	}

	appname, name, _ := connServer.GetInfo()
	if ret := this.options.CheckAppName(appname); !ret {
		err := fmt.Errorf("application name=%s is not configured", appname)
		conn.Close()
		log.Errorf("CheckAppName err: ", err)
		return err
	}
	log.Debugf("handleConn: IsPublisher=%v", connServer.IsPublisher())
	if connServer.IsPublisher() {
		if this.options.GetRtmpNoAuth() {
			key, err := this.module.GetCacheComp().GetChannelKey(name)
			if err != nil {
				err := fmt.Errorf("Cannot create key err=%s", err.Error())
				conn.Close()
				log.Errorf("GetKey err: ", err)
				return err
			}
			name = key
		}
		channel, err := this.module.GetCacheComp().GetChannel(name)
		if err != nil {
			err := fmt.Errorf("invalid key err=%s", err.Error())
			conn.Close()
			log.Errorf("CheckKey err: ", err)
			return err
		}
		connServer.PublishInfo.Name = channel
		if pushlist, ret := this.options.GetStaticPushUrlList(appname); ret && (pushlist != nil) {
			log.Debugf("GetStaticPushUrlList: %v", pushlist)
		}
		reader := NewVirReader(connServer)
		this.module.GetHandler().HandleReader(reader)
		log.Debugf("new publisher: %+v", reader.Info())

		if this.module.GetGetter() != nil {
			writeType := reflect.TypeOf(this.module.GetGetter())
			log.Debugf("handleConn:writeType=%v", writeType)
			writer := this.module.GetGetter().GetWriter(reader.Info())
			this.module.GetHandler().HandleWriter(writer)
		}
		if this.options.GetFLVArchive() {
			flvWriter := new(flv.FlvDvr)
			this.module.GetHandler().HandleWriter(flvWriter.GetWriter(this.options.GetFLVDir(), reader.Info()))
		}
	} else {
		writer := NewVirWriter(this.options.GetWriteTimeout(), connServer)
		log.Debugf("new player: %+v", writer.Info())
		this.module.GetHandler().HandleWriter(writer)
	}
	return
}
