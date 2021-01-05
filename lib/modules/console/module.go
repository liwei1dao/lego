package console

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/lib/modules/http"
)

type Console struct {
	http.Http
}

func (this *Console) Init(service core.IService, module core.IModule, options core.IModuleOptions) (err error) {
	err = this.Http.Init(service, module, options)
	return
}
