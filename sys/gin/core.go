package gin

import "net/http"

type (
	ISys interface {
		IRoutes
		Close() (err error)
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
	RouteInfo struct {
		Method      string
		Path        string
		Handler     string
		HandlerFunc HandlerFunc
	}
	RoutesInfo []RouteInfo
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

func Close() (err error) {
	return defsys.Close()
}

func Use(handlers ...HandlerFunc) IRoutes {
	return defsys.Use(handlers...)
}
func Handle(httpMethod string, relativePath string, handlers ...HandlerFunc) IRoutes {
	return defsys.Handle(httpMethod, relativePath, handlers...)
}
func Any(relativePath string, handlers ...HandlerFunc) IRoutes {
	return defsys.Any(relativePath, handlers...)
}
func GET(httpMethod string, handlers ...HandlerFunc) IRoutes {
	return defsys.GET(httpMethod, handlers...)
}
func POST(httpMethod string, handlers ...HandlerFunc) IRoutes {
	return defsys.POST(httpMethod, handlers...)
}
func DELETE(httpMethod string, handlers ...HandlerFunc) IRoutes {
	return defsys.DELETE(httpMethod, handlers...)
}
func PATCH(httpMethod string, handlers ...HandlerFunc) IRoutes {
	return defsys.PATCH(httpMethod, handlers...)
}
func PUT(httpMethod string, handlers ...HandlerFunc) IRoutes {
	return defsys.PUT(httpMethod, handlers...)
}
func OPTIONS(httpMethod string, handlers ...HandlerFunc) IRoutes {
	return defsys.OPTIONS(httpMethod, handlers...)
}
func HEAD(httpMethod string, handlers ...HandlerFunc) IRoutes {
	return defsys.HEAD(httpMethod, handlers...)
}
func StaticFile(relativePath string, filepath string) IRoutes {
	return defsys.StaticFile(relativePath, filepath)
}
func StaticFileFS(relativePath string, filepath string, fs http.FileSystem) IRoutes {
	return defsys.StaticFileFS(relativePath, filepath, fs)
}

func Static(relativePath string, root string) IRoutes {
	return defsys.Static(relativePath, root)
}

func StaticFS(relativePath string, fs http.FileSystem) IRoutes {
	return defsys.StaticFS(relativePath, fs)
}
