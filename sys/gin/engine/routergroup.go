package engine

import (
	"net/http"
	"path"
	"reflect"
	"regexp"
	"strings"
)

var (
	regEnLetter = regexp.MustCompile("^[A-Z]+$")
	anyMethods  = []string{
		http.MethodGet, http.MethodPost, http.MethodPut, http.MethodPatch,
		http.MethodHead, http.MethodOptions, http.MethodDelete, http.MethodConnect,
		http.MethodTrace,
	}
)

type RouterGroup struct {
	Handlers HandlersChain
	basePath string
	engine   *Engine
	root     bool
}

func (this *RouterGroup) BasePath() string {
	return this.basePath
}
func (this *RouterGroup) NoRoute(handlers ...HandlerFunc) {
	this.engine.NoRoute(handlers...)
}
func (this *RouterGroup) Use(middleware ...HandlerFunc) IRoutes {
	this.Handlers = append(this.Handlers, middleware...)
	return this.returnObj()
}

func (this *RouterGroup) Group(relativePath string, handlers ...HandlerFunc) IRoutes {
	return &RouterGroup{
		Handlers: this.combineHandlers(handlers),
		basePath: this.calculateAbsolutePath(relativePath),
		engine:   this.engine,
	}
}

func (this *RouterGroup) Handle(httpMethod, relativePath string, handlers ...HandlerFunc) IRoutes {
	if matched := regEnLetter.MatchString(httpMethod); !matched {
		panic("http method " + httpMethod + " is not valid")
	}
	return this.handle(httpMethod, relativePath, handlers)
}

func (this *RouterGroup) POST(relativePath string, handlers ...HandlerFunc) IRoutes {
	return this.handle(http.MethodPost, relativePath, handlers)
}

func (this *RouterGroup) GET(relativePath string, handlers ...HandlerFunc) IRoutes {
	return this.handle(http.MethodGet, relativePath, handlers)
}

func (this *RouterGroup) DELETE(relativePath string, handlers ...HandlerFunc) IRoutes {
	return this.handle(http.MethodDelete, relativePath, handlers)
}

func (this *RouterGroup) PATCH(relativePath string, handlers ...HandlerFunc) IRoutes {
	return this.handle(http.MethodPatch, relativePath, handlers)
}

func (this *RouterGroup) PUT(relativePath string, handlers ...HandlerFunc) IRoutes {
	return this.handle(http.MethodPut, relativePath, handlers)
}

func (this *RouterGroup) OPTIONS(relativePath string, handlers ...HandlerFunc) IRoutes {
	return this.handle(http.MethodOptions, relativePath, handlers)
}

func (this *RouterGroup) HEAD(relativePath string, handlers ...HandlerFunc) IRoutes {
	return this.handle(http.MethodHead, relativePath, handlers)
}

func (this *RouterGroup) Any(relativePath string, handlers ...HandlerFunc) IRoutes {
	for _, method := range anyMethods {
		this.handle(method, relativePath, handlers)
	}
	return this.returnObj()
}

func (this *RouterGroup) StaticFile(relativePath, filepath string) IRoutes {
	return this.staticFileHandler(relativePath, func(c *Context) {
		c.File(filepath)
	})
}

func (this *RouterGroup) StaticFileFS(relativePath, filepath string, fs http.FileSystem) IRoutes {
	return this.staticFileHandler(relativePath, func(c *Context) {
		c.FileFromFS(filepath, fs)
	})
}

func (this *RouterGroup) Static(relativePath, root string) IRoutes {
	return this.StaticFS(relativePath, Dir(root, false))
}

func (this *RouterGroup) StaticFS(relativePath string, fs http.FileSystem) IRoutes {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static folder")
	}
	handler := this.createStaticHandler(relativePath, fs)
	urlPattern := path.Join(relativePath, "/*filepath")

	// Register GET and HEAD handlers
	this.GET(urlPattern, handler)
	this.HEAD(urlPattern, handler)
	return this.returnObj()
}

func (this *RouterGroup) handle(httpMethod, relativePath string, handlers HandlersChain) IRoutes {
	absolutePath := this.calculateAbsolutePath(relativePath)
	handlers = this.combineHandlers(handlers)
	this.engine.addRoute(httpMethod, absolutePath, handlers)
	return this.returnObj()
}

func (group *RouterGroup) createStaticHandler(relativePath string, fs http.FileSystem) HandlerFunc {
	absolutePath := group.calculateAbsolutePath(relativePath)
	fileServer := http.StripPrefix(absolutePath, http.FileServer(fs))
	return func(c *Context) {
		if _, noListing := fs.(*onlyFilesFS); noListing {
			c.Writer.WriteHeader(http.StatusNotFound)
		}
		file := c.Param("filepath")
		f, err := fs.Open(file)
		if err != nil {
			c.Writer.WriteHeader(http.StatusNotFound)
			c.handlers = group.engine.noRoute
			c.index = -1
			return
		}
		f.Close()

		fileServer.ServeHTTP(c.Writer, c.Request)
	}
}

func (this *RouterGroup) combineHandlers(handlers HandlersChain) HandlersChain {
	finalSize := len(this.Handlers) + len(handlers)
	assert1(finalSize < int(abortIndex), "too many handlers")
	mergedHandlers := make(HandlersChain, finalSize)
	copy(mergedHandlers, this.Handlers)
	copy(mergedHandlers[len(this.Handlers):], handlers)
	return mergedHandlers
}

func (this *RouterGroup) calculateAbsolutePath(relativePath string) string {
	return joinPaths(this.basePath, relativePath)
}

func (this *RouterGroup) staticFileHandler(relativePath string, handler HandlerFunc) IRoutes {
	if strings.Contains(relativePath, ":") || strings.Contains(relativePath, "*") {
		panic("URL parameters can not be used when serving a static file")
	}
	this.GET(relativePath, handler)
	this.HEAD(relativePath, handler)
	return this.returnObj()
}

func (this *RouterGroup) returnObj() IRoutes {
	if this.root {
		return this.engine
	}
	return this
}

func (this *RouterGroup) Register(rcvr interface{}) {
	typ := reflect.TypeOf(rcvr)
	vof := reflect.ValueOf(rcvr)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mname := method.Name
		mtype := method.Type
		if method.PkgPath != "" {
			continue
		}
		if mtype.NumIn() != 2 {
			continue
		}
		context := mtype.In(1)
		if context.String() != "*engine.Context" {
			continue
		}
		if mtype.NumOut() != 0 {
			continue
		}
		this.POST(strings.ToLower(mname), vof.MethodByName(mname).Interface().(func(*Context)))
	}
}
