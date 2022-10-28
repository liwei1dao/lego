package gin

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/gin-gonic/autotls"
	"github.com/liwei1dao/lego/sys/gin/engine"
	"github.com/liwei1dao/lego/sys/gin/middleware/logger"
	"github.com/liwei1dao/lego/sys/gin/middleware/recovery"
)

func newSys(options *Options) (sys *Gin, err error) {
	sys = &Gin{
		options: options,
	}
	sys.engine = engine.NewEngine(options.Log)
	///添加基础中间件
	sys.engine.Use(logger.Logger([]string{}), recovery.Recovery())
	if options.CertFile != "" && options.KeyFile != "" {
		sys.RunTLS(options.ListenPort, options.CertFile, options.KeyFile)
	} else if options.LetEncrypt {
		sys.RunLetEncrypt(options.Domain...)
	} else {
		sys.Run(options.ListenPort)
	}
	return
}

type Gin struct {
	options *Options
	server  *http.Server
	engine  *engine.Engine
}

func (this *Gin) Run(listenPort int) (err error) {
	// if this.engine.IsUnsafeTrustedProxies() {
	// 	this.Warnf("You trusted all proxies, this is NOT safe. We recommend you to set a value.\n" +
	// 		"Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.")

	// }
	this.options.Log.Debugf("Listening and serving HTTP on:%d", listenPort)
	this.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", listenPort),
		Handler: this.engine.Handler(),
	}
	go func() {
		if err := this.server.ListenAndServe(); err != nil && errors.Is(err, http.ErrServerClosed) {
			this.options.Log.Errorln(err)
		}
	}()
	// err = http.ListenAndServe(fmt.Sprintf(":%d", this.options.ListenPort), this.Handler())
	return
}

func (this *Gin) RunTLS(listenPort int, certFile, keyFile string) (err error) {
	this.options.Log.Debugf("Listening and serving HTTPS on :%d", listenPort)
	// if this.engine.IsUnsafeTrustedProxies() {
	// 	this.Warnf("You trusted all proxies, this is NOT safe. We recommend you to set a value.\n" +
	// 		"Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.")
	// }
	this.server = &http.Server{
		Addr:    fmt.Sprintf(":%d", listenPort),
		Handler: this.engine.Handler(),
	}
	go func() {
		if err := this.server.ListenAndServeTLS(certFile, keyFile); err != nil && errors.Is(err, http.ErrServerClosed) {
			this.options.Log.Errorln(err)
		}
	}()
	// err = http.ListenAndServeTLS(addr, certFile, keyFile, this.Handler())
	return
}

func (this *Gin) RunLetEncrypt(domain ...string) {
	this.options.Log.Debugf("Listening and serving LetEncrypt on :%v", domain)
	go func() {
		if err := autotls.Run(this.engine, domain...); err != nil && errors.Is(err, http.ErrServerClosed) {
			this.options.Log.Errorln(err)
		}
	}()
}

func (this *Gin) RunListener(listener net.Listener) (err error) {
	this.options.Log.Debugf("Listening and serving HTTP on listener what's bind with address@%s", listener.Addr())
	defer func() {
		if err != nil {
			this.options.Log.Errorln(err)
		}
	}()
	// if this.engine.IsUnsafeTrustedProxies() {
	// 	this.Warnf("You trusted all proxies, this is NOT safe. We recommend you to set a value.\n" +
	// 		"Please check https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies for details.")

	// }
	err = http.Serve(listener, this.engine.Handler())
	return
}
func (this *Gin) HandleContext(c *engine.Context) {
	this.HandleContext(c)
}
func (this *Gin) LoadHTMLGlob(pattern string) {
	this.engine.LoadHTMLGlob(pattern)
}

func (this *Gin) Close() (err error) {
	if err = this.server.Shutdown(context.Background()); err != nil {
		this.options.Log.Errorln(err)
	}
	this.server.Close()
	return
}
func (this *Gin) Register(rcvr interface{}) {
	this.engine.Register(rcvr)
}

func (this *Gin) NoRoute(handlers ...engine.HandlerFunc) {
	this.engine.NoRoute(handlers...)
}

func (this *Gin) Group(relativePath string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return this.engine.Group(relativePath, handlers...)
}
func (this *Gin) Use(handlers ...engine.HandlerFunc) engine.IRoutes {
	return this.engine.Use(handlers...)
}
func (this *Gin) Handle(httpMethod string, relativePath string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return this.engine.Handle(httpMethod, relativePath, handlers...)
}
func (this *Gin) Any(relativePath string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return this.engine.Any(relativePath, handlers...)
}
func (this *Gin) GET(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return this.engine.GET(httpMethod, handlers...)
}
func (this *Gin) POST(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return this.engine.POST(httpMethod, handlers...)
}
func (this *Gin) DELETE(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return this.engine.DELETE(httpMethod, handlers...)
}
func (this *Gin) PATCH(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.PATCH(httpMethod, handlers...)
}
func (this *Gin) PUT(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return this.engine.PUT(httpMethod, handlers...)
}
func (this *Gin) OPTIONS(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return this.engine.OPTIONS(httpMethod, handlers...)
}
func (this *Gin) HEAD(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return this.engine.HEAD(httpMethod, handlers...)
}
func (this *Gin) StaticFile(relativePath string, filepath string) engine.IRoutes {
	return this.engine.StaticFile(relativePath, filepath)
}
func (this *Gin) StaticFileFS(relativePath string, filepath string, fs http.FileSystem) engine.IRoutes {
	return this.engine.StaticFileFS(relativePath, filepath, fs)
}

func (this *Gin) Static(relativePath string, root string) engine.IRoutes {
	return this.engine.Static(relativePath, root)
}

func (this *Gin) StaticFS(relativePath string, fs http.FileSystem) engine.IRoutes {
	return this.engine.StaticFS(relativePath, fs)
}
