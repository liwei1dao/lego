package gin

import "net/http"

type (
	ISys interface {
		IRoutes
	}
	IRoutes interface {
		Use(...HandlerFunc) IRoutes
		Handle(string, string, ...HandlerFunc) IRoutes
		Any(string, ...HandlerFunc) IRoutes
		GET(string, ...HandlerFunc) IRoutes
		POST(string, ...HandlerFunc) IRoutes
		DELETE(string, ...HandlerFunc) IRoutes
		PATCH(string, ...HandlerFunc) IRoutes
		PUT(string, ...HandlerFunc) IRoutes
		OPTIONS(string, ...HandlerFunc) IRoutes
		HEAD(string, ...HandlerFunc) IRoutes
		StaticFile(string, string) IRoutes
		StaticFileFS(string, string, http.FileSystem) IRoutes
		Static(string, string) IRoutes
		StaticFS(string, http.FileSystem) IRoutes
	}
)

var defsys ISys

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}
