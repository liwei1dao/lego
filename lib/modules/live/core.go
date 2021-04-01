package live

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/lib/modules/live/av"
)

type (
	ILive interface {
		core.IModule
		GetOptions() IOptions
		GetGetter() av.GetWriter
		GetHandler() av.Handler
		GetCacheComp() ICacheComp
	}
	ICacheComp interface {
		core.IModuleComp
		GetChannel(key string) (channel string, err error)
		SetChannelKey(channel string) (key string, err error)
		GetChannelKey(channel string) (newKey string, err error)
		DeleteChannel(channel string) bool
	}
)
