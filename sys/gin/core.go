package gin

import (
	"net/http"

	"github.com/liwei1dao/lego/sys/gin/engine"
)

type ISys interface {
	engine.IRoutes
	Close() (err error)
}

var defsys ISys

func OnInit(config map[string]interface{}, opt ...Option) (err error) {
	var option *Options
	if option, err = newOptions(config, opt...); err != nil {
		return
	}
	defsys, err = newSys(option)
	return
}

func NewSys(opt ...Option) (sys ISys, err error) {
	var option *Options
	if option, err = newOptionsByOption(opt...); err != nil {
		return
	}
	sys, err = newSys(option)
	return
}

func Close() (err error) {
	return defsys.Close()
}

func Use(handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.Use(handlers...)
}
func Handle(httpMethod string, relativePath string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.Handle(httpMethod, relativePath, handlers...)
}
func Any(relativePath string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.Any(relativePath, handlers...)
}
func GET(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.GET(httpMethod, handlers...)
}
func POST(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.POST(httpMethod, handlers...)
}
func DELETE(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.DELETE(httpMethod, handlers...)
}
func PATCH(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.PATCH(httpMethod, handlers...)
}
func PUT(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.PUT(httpMethod, handlers...)
}
func OPTIONS(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.OPTIONS(httpMethod, handlers...)
}
func HEAD(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.HEAD(httpMethod, handlers...)
}
func StaticFile(relativePath string, filepath string) engine.IRoutes {
	return defsys.StaticFile(relativePath, filepath)
}
func StaticFileFS(relativePath string, filepath string, fs http.FileSystem) engine.IRoutes {
	return defsys.StaticFileFS(relativePath, filepath, fs)
}

func Static(relativePath string, root string) engine.IRoutes {
	return defsys.Static(relativePath, root)
}

func StaticFS(relativePath string, fs http.FileSystem) engine.IRoutes {
	return defsys.StaticFS(relativePath, fs)
}
