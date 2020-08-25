package http

import (
	"context"
	"fmt"
	"html/template"
	"net/http"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib/modules/http/render"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/utils"
)

type Http struct {
	cbase.ModuleBase
	RouterGroup
	service            core.IService
	http               *http.Server
	wg                 sync.WaitGroup
	certPath           string
	keyPath            string
	MaxMultipartMemory int64 //上传文件最大尺寸
	allNoRoute         HandlersChain
	allNoMethod        HandlersChain
	noMethod           HandlersChain
	noRoute            HandlersChain
	delims             render.Delims
	FuncMap            template.FuncMap
	HTMLRender         render.HTMLRender
	pool               sync.Pool
	trees              methodTrees
	Ip                 string
}

func (this *Http) Init(service core.IService, module core.IModule, setting map[string]interface{}) (err error) {
	this.service = service
	this.RouterGroup = RouterGroup{
		Handlers: nil,
		basePath: "/",
		root:     true,
	}
	this.MaxMultipartMemory = defaultMultipartMemory
	this.trees = make(methodTrees, 0, 9)
	this.RouterGroup.engine = this
	this.pool.New = func() interface{} {
		return this.allocateContext()
	}

	var httpAddr string
	if _, ok := setting["HttpAddr"]; !ok {
		return fmt.Errorf("Http Module Init HttpAddr 'Config' Is Null")
	}
	httpAddr = setting["HttpAddr"].(string)
	if setting["CertPath"] != nil && setting["KeyPath"] != nil {
		this.certPath = setting["CertPath"].(string)
		this.keyPath = setting["KeyPath"].(string)
	}
	this.http = &http.Server{
		Addr:    httpAddr,
		Handler: this,
	}
	if err = this.ModuleBase.Init(service, module, setting); err != nil {
		return
	}
	this.Ip = utils.GetOutboundIP()
	return
}
func (this *Http) Start() (err error) {
	err = this.ModuleBase.Start()
	this.wg.Add(1)
	go this.starthttp()
	return
}
func (this *Http) Destroy() (err error) {
	this.wg.Add(1)
	go this.closehttp()
	this.wg.Wait()
	err = this.ModuleBase.Destroy()
	return
}
func (this *Http) starthttp() {
	var err error
	if this.certPath != "" && this.keyPath != "" {
		err = this.http.ListenAndServeTLS(this.service.GetWorkPath()+this.certPath, this.service.GetWorkPath()+this.keyPath)
	} else {
		err = this.http.ListenAndServe()
	}
	if err != nil {
		log.Errorf("启动http服务错误%s", err)
	}
	this.wg.Done()
}
func (this *Http) closehttp() {
	//使用context控制srv.Shutdown的超时时间
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := this.http.Shutdown(ctx)
	if err != nil {
		log.Errorf("关闭Http服务组建失败%s", err.Error())
	}
	this.http.Close()
	this.wg.Done()
}

func (this *Http) Use(middleware ...HandlerFunc) IRoutes {
	this.RouterGroup.Use(middleware...)
	this.rebuild404Handlers()
	this.rebuild405Handlers()
	return this
}

//添加到路由树中
func (this *Http) addRoute(method, path string, handlers HandlersChain) (err error) {
	if err = outErr(path[0] == '/', "path must begin with '/'"); err != nil {
		return
	}
	if err = outErr(method != "", "HTTP method can not be empty"); err != nil {
		return
	}
	if err = outErr(len(handlers) > 0, "there must be at least one handler"); err != nil {
		return
	}

	root := this.trees.get(method)
	if root == nil {
		root = new(node)
		this.trees = append(this.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)
	return
}
func (this *Http) allocateContext() *Context {
	return &Context{engine: this}
}
func (this *Http) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := this.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = req
	c.reset()
	this.handleHTTPRequest(c)
	this.pool.Put(c)
}

func (this *Http) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	unescape := false
	rPath = cleanPath(rPath)
	t := this.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		handlers, params, _ := root.getValue(rPath, c.Params, unescape)
		if handlers != nil {
			c.handlers = handlers
			c.Params = params
			c.Next()
			c.writermem.WriteHeaderNow()
			return
		}
		break
	}
	c.handlers = this.allNoRoute
	serveError(c, http.StatusNotFound, default404Body)
}

// NoRoute adds handlers for NoRoute. It return a 404 code by default.
func (this *Http) NoRoute(handlers ...HandlerFunc) {
	this.noRoute = handlers
	this.rebuild404Handlers()
}

// NoMethod sets the handlers called when... TODO.
func (this *Http) NoMethod(handlers ...HandlerFunc) {
	this.noMethod = handlers
	this.rebuild405Handlers()
}

func (this *Http) rebuild404Handlers() {
	this.allNoRoute = this.combineHandlers(this.noRoute)
}
func (this *Http) rebuild405Handlers() {
	this.allNoMethod = this.combineHandlers(this.noMethod)
}

func (this *Http) LoadHTMLFiles(files ...string) {
	templ := template.Must(template.New("").Delims(this.delims.Left, this.delims.Right).Funcs(this.FuncMap).ParseFiles(files...))
	this.SetHTMLTemplate(templ)
}
func (this *Http) LoadHTMLGlob(pattern string) {
	left := this.delims.Left
	right := this.delims.Right
	templ := template.Must(template.New("").Delims(left, right).Funcs(this.FuncMap).ParseGlob(pattern))
	this.SetHTMLTemplate(templ)
}
func (this *Http) SetHTMLTemplate(templ *template.Template) {
	this.HTMLRender = render.HTMLProduction{Template: templ.Funcs(this.FuncMap)}
}
