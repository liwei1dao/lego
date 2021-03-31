package live

import (
	"fmt"
	"net"
	"reflect"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib"
	"github.com/liwei1dao/lego/lib/modules/live/av"
	"github.com/liwei1dao/lego/lib/modules/live/container/flv"
	"github.com/liwei1dao/lego/lib/modules/live/liveconn"
	"github.com/liwei1dao/lego/sys/log"
)

type Live struct {
	cbase.ModuleBase
	options    IOptions
	handler    av.Handler
	getter     av.GetWriter
	rtmpListen net.Listener
	cachecomp  *CacheComp
}

func (this *Live) GetType() core.M_Modules {
	return lib.SM_MonitorModule
}

func (this *Live) NewOptions() (options core.IModuleOptions) {
	return new(Options)
}

func (this *Live) Init(service core.IService, module core.IModule, options core.IModuleOptions) (err error) {
	this.options = options.(IOptions)
	stream := NewRtmpStream(this.options)
	this.handler = stream
	return this.ModuleBase.Init(service, module, options)
}

func (this *Live) OnInstallComp() {
	this.ModuleBase.OnInstallComp()
	this.cachecomp = this.RegisterComp(new(CacheComp)).(*CacheComp)
}

func (this *Live) Start() (err error) {
	this.rtmpListen, err = net.Listen("tcp", this.options.GetRtmpAddr())
	return this.ModuleBase.Start()
}

func (this *Live) Run(closeSig chan bool) (err error) {
	for {
		var netconn net.Conn
		netconn, err = this.rtmpListen.Accept()
		if err != nil {
			break
		}
		conn := liveconn.NewConn(netconn, 4*1024)
		log.Debugf("new client, connect remote: ", conn.RemoteAddr().String(),
			"local:", conn.LocalAddr().String())
		go this.handleConn(conn)
	}
	closeSig <- true
	return
}

func (this *Live) handleConn(conn *liveconn.Conn) (err error) {
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
			key, err := this.cachecomp.GetChannelKey(name)
			if err != nil {
				err := fmt.Errorf("Cannot create key err=%s", err.Error())
				conn.Close()
				log.Errorf("GetKey err: ", err)
				return err
			}
			name = key
		}
		channel, err := this.cachecomp.GetChannel(name)
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
		this.handler.HandleReader(reader)
		log.Debugf("new publisher: %+v", reader.Info())

		if this.getter != nil {
			writeType := reflect.TypeOf(this.getter)
			log.Debugf("handleConn:writeType=%v", writeType)
			writer := this.getter.GetWriter(reader.Info())
			this.handler.HandleWriter(writer)
		}
		if this.options.GetFLVArchive() {
			flvWriter := new(flv.FlvDvr)
			this.handler.HandleWriter(flvWriter.GetWriter(this.options.GetFLVDir(), reader.Info()))
		}
	} else {
		writer := NewVirWriter(this.options.GetWriteTimeout(), connServer)
		log.Debugf("new player: %+v", writer.Info())
		this.handler.HandleWriter(writer)
	}
	return
}
