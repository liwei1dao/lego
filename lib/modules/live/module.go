package live

import (
	"net"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib/modules/live/av"
)

type Live struct {
	cbase.ModuleBase
	options    IOptions
	handler    av.Handler
	getter     av.GetWriter
	rtmpListen net.Listener
	cachecomp  ICacheComp
	apicomp    *ApiComp
	rtmpcomp   *RtmpComp
}

func (this *Live) NewOptions() (options core.IModuleOptions) {
	return new(Options)
}

func (this *Live) GetOptions() (options IOptions) {
	return this.options
}

func (this *Live) Init(service core.IService, module core.IModule, options core.IModuleOptions) (err error) {
	this.options = options.(IOptions)
	stream := NewRtmpStream(this.options)
	this.handler = stream
	return this.ModuleBase.Init(service, module, options)
}

func (this *Live) OnInstallComp() {
	this.ModuleBase.OnInstallComp()
	this.cachecomp = this.RegisterComp(new(CacheComp)).(ICacheComp)
	this.rtmpcomp = this.RegisterComp(new(RtmpComp)).(*RtmpComp)
	this.apicomp = this.RegisterComp(new(ApiComp)).(*ApiComp)
}

func (this *Live) GetHandler() av.Handler {
	return this.handler
}

func (this *Live) GetGetter() av.GetWriter {
	return this.getter
}

func (this *Live) GetCacheComp() ICacheComp {
	return this.cachecomp
}
