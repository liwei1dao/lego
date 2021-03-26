package live

import "github.com/liwei1dao/lego/core"

type (
	ILive interface {
		core.IModule
		GetOptions() IOptions
	}
)

var module ILive
