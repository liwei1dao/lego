package engine

import (
	"html/template"
	"net"
	"net/http"
	"path"
	"strings"
	"sync"

	"github.com/liwei1dao/lego/sys/gin/render"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/codec"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

/*
	默认文件上传的最大尺寸
*/
const defaultMultipartMemory = 32 << 20 // 32 MB
/*
	默认可信代理
*/
var defaultTrustedCIDRs = []*net.IPNet{
	{ // 0.0.0.0/0 (IPv4)
		IP:   net.IP{0x0, 0x0, 0x0, 0x0},
		Mask: net.IPMask{0x0, 0x0, 0x0, 0x0},
	},
	{ // ::/0 (IPv6)
		IP:   net.IP{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
		Mask: net.IPMask{0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0, 0x0},
	},
}

func NewEngine(log log.ILogger) (engine *Engine) {
	engine = &Engine{
		RouterGroup: RouterGroup{
			Handlers: nil,
			basePath: "/",
			root:     true,
		},
		log:                    log,
		FuncMap:                template.FuncMap{},
		RedirectTrailingSlash:  true,
		RedirectFixedPath:      false,
		HandleMethodNotAllowed: false,
		ForwardedByClientIP:    true,
		RemoteIPHeaders:        []string{"X-Forwarded-For", "X-Real-IP"},
		UseRawPath:             false,
		RemoveExtraSlash:       false,
		UnescapePathValues:     true,
		MaxMultipartMemory:     defaultMultipartMemory,
		trees:                  make(methodTrees, 0, 9),
		delims:                 render.Delims{Left: "{{", Right: "}}"},
		secureJSONPrefix:       "while(1);",
		trustedProxies:         []string{"0.0.0.0/0"},
		trustedCIDRs:           defaultTrustedCIDRs,
	}
	engine.RouterGroup.engine = engine
	engine.pool.New = func() interface{} {
		return engine.allocateContext()
	}
	return
}

var (
	default404Body = []byte("404 page not found")
	default405Body = []byte("405 method not allowed")
)
var mimePlain = []string{MIMEPlain}

type Engine struct {
	RouterGroup
	log        log.ILogger
	UseRawPath bool
	/*
		如果启用，路由器尝试修复当前请求路径，如果没有
		如果没有
		已为其注册句柄。
		第一个多余的路径元素，如 ../ 或 // 被删除。
		之后路由器对清理后的路径进行不区分大小写的查找。
		如果可以找到该路由的句柄，则路由器进行重定向
		到正确的路径，GET 请求的状态码为 301，而 GET 请求的状态码为 307
		所有其他请求方法。
		例如 /FOO 和 /..//Foo 可以重定向到 /foo。
		RedirectTrailingSlash 与此选项无关。
	*/
	RedirectFixedPath bool
	/*
		如果为真，路径值将不转义
	*/
	UnescapePathValues bool
	/*
		如果当前路由无法匹配，但启用自动重定向
		带有（不带）尾部斜杠的路径的处理程序存在。
		例如，如果 /foo/ 被请求，但路由只存在于 /foo，则
		对于 GET 请求，客户端被重定向到 /foo，http 状态码为 301
		对于所有其他请求方法，则为 307。
	*/
	RedirectTrailingSlash bool //

	/*
		如果启用，路由器检查是否允许其他方法
		当前路由，如果当前请求无法路由。
		如果是这种情况，则使用“不允许的方法”回答请求
		和 HTTP 状态码 405。
		如果不允许其他方法，则将请求委托给 NotFound
		处理程序。
	*/
	HandleMethodNotAllowed bool
	/*
		可以从 URL 中解析出一个参数，即使带有额外的斜杠。
	*/
	RemoveExtraSlash bool
	/*
		TrustedPlatform 如果设置为值 gin.Platform* 的常量，则信任由设置的标头
		那个平台，比如判断客户端IP
	*/
	TrustedPlatform string
	/*
		ForwardedByClientIP 如果启用，客户端 IP 将从请求的标头中解析
		匹配存储在 `(*gin.Engine).RemoteIPHeaders` 中的那些。 如果没有 IP
		fetched, 它回退到从获取的 IP
		`(*gin.Context).Request.RemoteAddr`。
		ForwardedByClientIP 布尔值
	*/
	ForwardedByClientIP bool
	/*
		RemoteIPHeaders 用于获取客户端 IP 时的 headers 列表
		`(*gin.Engine).ForwardedByClientIP` 为 `true` 并且
		`(*gin.Context).Request.RemoteAddr` 被至少一个匹配
		由 `(*gin.Engine).SetTrustedProxies()` 定义的列表的网络来源。
	*/
	RemoteIPHeaders []string
	/*
		文件上传的最大尺寸
	*/
	MaxMultipartMemory int64
	/*
		是否使用H2C
	*/
	UseH2C           bool
	delims           render.Delims
	secureJSONPrefix string
	HTMLRender       render.HTMLRender
	FuncMap          template.FuncMap
	noRoute          HandlersChain
	noMethod         HandlersChain
	allNoRoute       HandlersChain
	allNoMethod      HandlersChain
	pool             sync.Pool
	trees            methodTrees
	maxParams        uint16
	maxSections      uint16
	trustedProxies   []string
	trustedCIDRs     []*net.IPNet
}

func (this *Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c := this.pool.Get().(*Context)
	c.writermem.reset(w)
	c.Request = req
	c.reset()
	this.handleHTTPRequest(c)
	this.pool.Put(c)
}

func (this *Engine) Handler() http.Handler {
	if !this.UseH2C {
		return this
	}

	h2s := &http2.Server{}
	return h2c.NewHandler(this, h2s)
}

/*
	使用中间件
*/
func (this *Engine) Use(middleware ...HandlerFunc) IRoutes {
	this.RouterGroup.Use(middleware...)
	this.rebuild404Handlers()
	this.rebuild405Handlers()
	return this
}

/*
	LoadHTMLGlob 加载由 glob 模式标识的 HTML 文件
	并将结果与 HTML 渲染器相关联。
*/
func (this *Engine) LoadHTMLGlob(pattern string) {
	left := this.delims.Left
	right := this.delims.Right
	templ := template.Must(template.New("").Delims(left, right).Funcs(this.FuncMap).ParseGlob(pattern))

	if this.log.Enabled(log.DebugLevel) {
		this.debugPrintLoadTemplate(templ)
		this.HTMLRender = render.HTMLDebug{Glob: pattern, FuncMap: this.FuncMap, Delims: this.delims}
		return
	}

	this.SetHTMLTemplate(templ)
}

/*
LoadHTMLFiles 加载一段 HTML 文件
并将结果与 HTML 渲染器相关联。
*/
func (this *Engine) LoadHTMLFiles(files ...string) {
	if this.log.Enabled(log.DebugLevel) {
		this.HTMLRender = render.HTMLDebug{Files: files, FuncMap: this.FuncMap, Delims: this.delims}
		return
	}
	templ := template.Must(template.New("").Delims(this.delims.Left, this.delims.Right).Funcs(this.FuncMap).ParseFiles(files...))
	this.SetHTMLTemplate(templ)
}

func (this *Engine) SetHTMLTemplate(templ *template.Template) {
	if len(this.trees) > 0 {
		this.log.Warnf(`Since SetHTMLTemplate() is NOT thread-safe. It should only be called
			at initialization. ie. before any route is registered or the router is listening in a socket:
				router := gin.Default()
				router.SetHTMLTemplate(template) // << good place
			`)

	}
	this.HTMLRender = render.HTMLProduction{Template: templ.Funcs(this.FuncMap)}
}

/*
	设置template FuncMap
*/
func (engine *Engine) SetFuncMap(funcMap template.FuncMap) {
	engine.FuncMap = funcMap
}

/*
	404 处理路由
*/
func (this *Engine) NoRoute(handlers ...HandlerFunc) {
	this.noRoute = handlers
	this.rebuild404Handlers()
}

/*
	没有找到对应的方法
*/
func (this *Engine) NoMethod(handlers ...HandlerFunc) {
	this.noMethod = handlers
	this.rebuild405Handlers()
}

func (engine *Engine) Routes() (routes RoutesInfo) {
	for _, tree := range engine.trees {
		routes = iterate("", tree.method, routes, tree.root)
	}
	return routes
}

/*
	设置信任代理
*/
func (this *Engine) SetTrustedProxies(trustedProxies []string) error {
	this.trustedProxies = trustedProxies
	return this.parseTrustedProxies()
}

func (this *Engine) addRoute(method, path string, handlers HandlersChain) {
	assert1(path[0] == '/', "path must begin with '/'")
	assert1(method != "", "HTTP method can not be empty")
	assert1(len(handlers) > 0, "there must be at least one handler")
	if this.log.Enabled(log.DebugLevel) {
		nuHandlers := len(handlers)
		handlerName := nameOfFunction(handlers.Last())
		this.log.Debugf("%s:%s --> %s handlers:%d", method, path, handlerName, nuHandlers)
	}
	root := this.trees.get(method)
	if root == nil {
		root = new(node)
		root.fullPath = "/"
		this.trees = append(this.trees, methodTree{method: method, root: root})
	}
	root.addRoute(path, handlers)
	// Update maxParams
	if paramsCount := countParams(path); paramsCount > this.maxParams {
		this.maxParams = paramsCount
	}

	if sectionsCount := countSections(path); sectionsCount > this.maxSections {
		this.maxSections = sectionsCount
	}
}

//重定向
func (this *Engine) HandleContext(c *Context) {
	oldIndexValue := c.index
	c.reset()
	this.handleHTTPRequest(c)
	c.index = oldIndexValue
}

func (this *Engine) handleHTTPRequest(c *Context) {
	httpMethod := c.Request.Method
	rPath := c.Request.URL.Path
	unescape := false
	if this.UseRawPath && len(c.Request.URL.RawPath) > 0 {
		rPath = c.Request.URL.RawPath
		unescape = this.UnescapePathValues
	}
	if this.RemoveExtraSlash {
		rPath = cleanPath(rPath)
	}
	t := this.trees
	for i, tl := 0, len(t); i < tl; i++ {
		if t[i].method != httpMethod {
			continue
		}
		root := t[i].root
		// Find route in tree
		value := root.getValue(rPath, c.params, c.skippedNodes, unescape)
		if value.params != nil {
			c.Params = *value.params
		}
		if value.handlers != nil {
			c.handlers = value.handlers
			c.fullPath = value.fullPath
			c.Next()
			c.writermem.WriteHeaderNow()
			return
		}
		if httpMethod != http.MethodConnect && rPath != "/" {
			if value.tsr && this.RedirectTrailingSlash {
				this.redirectTrailingSlash(c)
				return
			}
			if this.RedirectFixedPath && this.redirectFixedPath(c, root, this.RedirectFixedPath) {
				return
			}
		}
		break
	}

	if this.HandleMethodNotAllowed {
		for _, tree := range this.trees {
			if tree.method == httpMethod {
				continue
			}
			if value := tree.root.getValue(rPath, nil, c.skippedNodes, unescape); value.handlers != nil {
				c.handlers = this.allNoMethod
				this.serveError(c, http.StatusMethodNotAllowed, default405Body)
				return
			}
		}
	}
	c.handlers = this.allNoRoute
	this.serveError(c, http.StatusNotFound, default404Body)
}

func (this *Engine) rebuild404Handlers() {
	this.allNoRoute = this.combineHandlers(this.noRoute)
}

func (this *Engine) rebuild405Handlers() {
	this.allNoMethod = this.combineHandlers(this.noMethod)
}

func (this *Engine) IsUnsafeTrustedProxies() bool {
	return this.isTrustedProxy(net.ParseIP("0.0.0.0")) || this.isTrustedProxy(net.ParseIP("::"))
}

//validateHeader 将解析 X-Forwarded-For 标头并返回受信任的客户端 IP 地址
func (this *Engine) validateHeader(header string) (clientIP string, valid bool) {
	if header == "" {
		return "", false
	}
	items := strings.Split(header, ",")
	for i := len(items) - 1; i >= 0; i-- {
		ipStr := strings.TrimSpace(items[i])
		ip := net.ParseIP(ipStr)
		if ip == nil {
			break
		}

		// X-Forwarded-For is appended by proxy
		// Check IPs in reverse order and stop when find untrusted proxy
		if (i == 0) || (!this.isTrustedProxy(ip)) {
			return ipStr, true
		}
	}
	return "", false
}

///目标Ip是否是可信
func (this *Engine) isTrustedProxy(ip net.IP) bool {
	if this.trustedCIDRs == nil {
		return false
	}
	for _, cidr := range this.trustedCIDRs {
		if cidr.Contains(ip) {
			return true
		}
	}
	return false
}

func (this *Engine) serveError(c *Context, code int, defaultMessage []byte) {
	c.writermem.status = code
	c.Next()
	if c.writermem.Written() {
		return
	}
	if c.writermem.Status() == code {
		c.writermem.Header()["Content-Type"] = mimePlain
		_, err := c.Writer.Write(defaultMessage)
		if err != nil {
			this.log.Errorf("[SYS-Gin] cannot write message to writer during serve error: %v", err)
		}
		return
	}
	c.writermem.WriteHeaderNow()
}

func (this *Engine) redirectFixedPath(c *Context, root *node, trailingSlash bool) bool {
	req := c.Request
	rPath := req.URL.Path

	if fixedPath, ok := root.findCaseInsensitivePath(cleanPath(rPath), trailingSlash); ok {
		req.URL.Path = codec.BytesToString(fixedPath)
		this.redirectRequest(c)
		return true
	}
	return false
}

func (this *Engine) redirectTrailingSlash(c *Context) {
	req := c.Request
	p := req.URL.Path
	if prefix := path.Clean(c.Request.Header.Get("X-Forwarded-Prefix")); prefix != "." {
		p = prefix + "/" + req.URL.Path
	}
	req.URL.Path = p + "/"
	if length := len(p); length > 1 && p[length-1] == '/' {
		req.URL.Path = p[:length-1]
	}
	this.redirectRequest(c)
}

func (this *Engine) redirectRequest(c *Context) {
	req := c.Request
	rPath := req.URL.Path
	rURL := req.URL.String()

	code := http.StatusMovedPermanently // Permanent redirect, request with GET method
	if req.Method != http.MethodGet {
		code = http.StatusTemporaryRedirect
	}
	this.log.Debugf("redirecting request %d: %s --> %s", code, rPath, rURL)
	http.Redirect(c.Writer, req, rURL, code)
	c.writermem.WriteHeaderNow()
}

func (this *Engine) parseTrustedProxies() error {
	trustedCIDRs, err := this.prepareTrustedCIDRs()
	this.trustedCIDRs = trustedCIDRs
	return err
}

func (this *Engine) prepareTrustedCIDRs() ([]*net.IPNet, error) {
	if this.trustedProxies == nil {
		return nil, nil
	}

	cidr := make([]*net.IPNet, 0, len(this.trustedProxies))
	for _, trustedProxy := range this.trustedProxies {
		if !strings.Contains(trustedProxy, "/") {
			ip := parseIP(trustedProxy)
			if ip == nil {
				return cidr, &net.ParseError{Type: "IP address", Text: trustedProxy}
			}

			switch len(ip) {
			case net.IPv4len:
				trustedProxy += "/32"
			case net.IPv6len:
				trustedProxy += "/128"
			}
		}
		_, cidrNet, err := net.ParseCIDR(trustedProxy)
		if err != nil {
			return cidr, err
		}
		cidr = append(cidr, cidrNet)
	}
	return cidr, nil
}

func (this *Engine) allocateContext() *Context {
	v := make(Params, 0, this.maxParams)
	skippedNodes := make([]skippedNode, 0, this.maxSections)
	return &Context{Log: this.log, engine: this, params: &v, skippedNodes: &skippedNodes}
}

//日志接口-------------------------------------------------------------
func (this *Engine) debugPrintLoadTemplate(tmpl *template.Template) {
	var buf strings.Builder
	for _, tmpl := range tmpl.Templates() {
		buf.WriteString("\t- ")
		buf.WriteString(tmpl.Name())
		buf.WriteString("\n")
	}
	format := "Loaded HTML Templates (%d): \n%s\n"
	if !strings.HasSuffix(format, "\n") {
		format += "\n"
	}
	this.log.Debugf(format, len(tmpl.Templates()), buf.String())

}
