package gin

import (
	"math"
	"net"
	"net/http"
	"strings"

	"github.com/liwei1dao/lego/sys/gin/binding"
)

const (
	MIMEJSON              = binding.MIMEJSON
	MIMEHTML              = binding.MIMEHTML
	MIMEXML               = binding.MIMEXML
	MIMEXML2              = binding.MIMEXML2
	MIMEPlain             = binding.MIMEPlain
	MIMEPOSTForm          = binding.MIMEPOSTForm
	MIMEMultipartPOSTForm = binding.MIMEMultipartPOSTForm
	MIMEYAML              = binding.MIMEYAML
)
const abortIndex int8 = math.MaxInt8 >> 1

type Context struct {
	writermem    responseWriter
	Request      *http.Request
	Writer       ResponseWriter
	Params       Params
	handlers     HandlersChain
	index        int8
	fullPath     string
	engine       *Engine
	params       *Params
	skippedNodes *[]skippedNode

	Keys   map[string]interface{}
	Errors errorMsgs
}

func (this *Context) Next() {
	this.index++
	for this.index < int8(len(this.handlers)) {
		this.handlers[this.index](this)
		this.index++
	}
}
func (this *Context) Param(key string) string {
	return this.Params.ByName(key)
}
func (this *Context) File(filepath string) {
	http.ServeFile(this.Writer, this.Request, filepath)
}

func (this *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		this.Request.URL.Path = old
	}(this.Request.URL.Path)
	this.Request.URL.Path = filepath
	http.FileServer(fs).ServeHTTP(this.Writer, this.Request)
}

func (c *Context) RemoteIP() string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(c.Request.RemoteAddr))
	if err != nil {
		return ""
	}
	return ip
}
func (this *Context) ClientIP() string {
	// 检查我们是否在受信任的平台上运行，如果出错则继续向后运行
	if this.engine.TrustedPlatform != "" {
		// Developers can define their own header of Trusted Platform or use predefined constants
		if addr := this.requestHeader(this.engine.TrustedPlatform); addr != "" {
			return addr
		}
	}

	/*
	   // 它还检查 remoteIP 是否是受信任的代理。
	   // 为了执行此验证，它将查看 IP 是否包含在至少一个 CIDR 块中
	   // 由 Engine.SetTrustedProxies() 定义
	*/
	remoteIP := net.ParseIP(this.RemoteIP())
	if remoteIP == nil {
		return ""
	}
	trusted := this.engine.isTrustedProxy(remoteIP)

	if trusted && this.engine.ForwardedByClientIP && this.engine.RemoteIPHeaders != nil {
		for _, headerName := range this.engine.RemoteIPHeaders {
			ip, valid := this.engine.validateHeader(this.requestHeader(headerName))
			if valid {
				return ip
			}
		}
	}
	return remoteIP.String()
}

func (this *Context) requestHeader(key string) string {
	return this.Request.Header.Get(key)
}

func (this *Context) reset() {
	this.Writer = &this.writermem
	this.Params = this.Params[:0]
	this.handlers = nil
	this.index = -1

	this.fullPath = ""
	this.Keys = nil
	this.Errors = this.Errors[:0]
	// this.Accepted = nil
	// this.queryCache = nil
	// this.formCache = nil
	// this.sameSite = 0
	*this.params = (*this.params)[:0]
	*this.skippedNodes = (*this.skippedNodes)[:0]
}
