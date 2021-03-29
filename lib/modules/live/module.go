package live

import (
	"net"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib"
	"github.com/liwei1dao/lego/lib/modules/live/liveconn"
	"github.com/liwei1dao/lego/sys/log"
)

type Live struct {
	cbase.ModuleBase
	options    IOptions
	rtmpListen net.Listener
}

func (this *Live) GetType() core.M_Modules {
	return lib.SM_MonitorModule
}

func (this *Live) NewOptions() (options core.IModuleOptions) {
	return new(Options)
}

func (this *Live) Init(service core.IService, module core.IModule, options core.IModuleOptions) (err error) {
	this.options = options.(IOptions)
	return this.ModuleBase.Init(service, module, options)
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
			return
		}
		conn := liveconn.NewConn(netconn, 4*1024)
		log.Debugf("new client, connect remote: ", conn.RemoteAddr().String(),
			"local:", conn.LocalAddr().String())
		go this.handleConn(conn)
	}
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
	return
}
