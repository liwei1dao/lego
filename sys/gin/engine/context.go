package engine

import (
	"errors"
	"io"
	"io/ioutil"
	"math"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/gin/binding"
	"github.com/liwei1dao/lego/sys/gin/render"
	"github.com/liwei1dao/lego/sys/log"
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

func newContext(log log.ILogger, engine *Engine, params *Params, skippedNodes *[]skippedNode) *Context {
	return &Context{
		engine:       engine,
		params:       params,
		skippedNodes: skippedNodes,
		writermem:    ResponseWriter{log: log},
	}
}

type Context struct {
	Log          log.ILogger
	engine       *Engine
	writermem    ResponseWriter
	Request      *http.Request
	Writer       IResponseWriter
	Params       Params
	handlers     HandlersChain
	index        int8
	fullPath     string
	params       *Params
	skippedNodes *[]skippedNode
	mu           sync.RWMutex
	Keys         map[string]interface{}
	Errors       errorMsgs
	Accepted     []string
	queryCache   url.Values
	formCache    url.Values
	sameSite     http.SameSite
}

func (this *Context) Copy() *Context {
	cp := Context{
		writermem: this.writermem,
		Request:   this.Request,
		Params:    this.Params,
		engine:    this.engine,
	}
	cp.writermem.ResponseWriter = nil
	cp.Writer = &cp.writermem
	cp.index = abortIndex
	cp.handlers = nil
	cp.Keys = map[string]interface{}{}
	for k, v := range this.Keys {
		cp.Keys[k] = v
	}
	paramCopy := make([]Param, len(cp.Params))
	copy(paramCopy, cp.Params)
	cp.Params = paramCopy
	return &cp
}

func (this *Context) HandlerName() string {
	return nameOfFunction(this.handlers.Last())
}

func (this *Context) HandlerNames() []string {
	hn := make([]string, 0, len(this.handlers))
	for _, val := range this.handlers {
		hn = append(hn, nameOfFunction(val))
	}
	return hn
}

func (this *Context) Handler() HandlerFunc {
	return this.handlers.Last()
}

/*
	FullPath 返回匹配的路由完整路径。 对于未找到的路线
	返回一个空字符串。
*/
func (c *Context) FullPath() string {
	return c.fullPath
}

func (this *Context) Next() {
	this.index++
	for this.index < int8(len(this.handlers)) {
		this.handlers[this.index](this)
		this.index++
	}
}

/*
	如果当前上下文被中止，IsAborted 返回 true。
*/
func (this *Context) IsAborted() bool {
	return this.index >= abortIndex
}

/*
	Abort 防止挂起的处理程序被调用。 请注意，这不会停止当前处理程序。
	假设你有一个授权中间件来验证当前请求是否被授权。
	如果授权失败（例如：密码不匹配），调用 Abort 以确保剩余的 handlers
	因为这个请求没有被调用。
*/
func (this *Context) Abort() {
	this.index = abortIndex
}

/*
	AbortWithStatus 调用 `Abort()` 并使用指定的状态代码写入标头。
	例如，验证请求失败的尝试可以使用：context.AbortWithStatus(401)。
*/
func (this *Context) AbortWithStatus(code int) {
	this.Status(code)
	this.Writer.WriteHeaderNow()
	this.Abort()
}

func (this *Context) AbortWithStatusJSON(code int, jsonObj interface{}) {
	this.Abort()
	this.JSON(code, jsonObj)
}

func (this *Context) AbortWithError(code int, err error) *Error {
	this.AbortWithStatus(code)
	return this.Error(err)
}

func (this *Context) Set(key string, value interface{}) {
	this.mu.Lock()
	if this.Keys == nil {
		this.Keys = make(map[string]interface{})
	}

	this.Keys[key] = value
	this.mu.Unlock()
}

func (this *Context) SetUserId(uid string) {
	this.Set("UserId", uid)
}

func (this *Context) Get(key string) (value interface{}, exists bool) {
	this.mu.RLock()
	value, exists = this.Keys[key]
	this.mu.RUnlock()
	return
}

/*
如果存在，MustGet 返回给定键的值，否则抛出异常。
*/
func (this *Context) MustGet(key string) interface{} {
	if value, exists := this.Get(key); exists {
		return value
	}
	panic("Key \"" + key + "\" does not exist")
}

func (this *Context) GetString(key string) (s string) {
	if val, ok := this.Get(key); ok && val != nil {
		s, _ = val.(string)
	}
	return
}

func (this *Context) GetBool(key string) (b bool) {
	if val, ok := this.Get(key); ok && val != nil {
		b, _ = val.(bool)
	}
	return
}

func (this *Context) GetInt(key string) (i int) {
	if val, ok := this.Get(key); ok && val != nil {
		i, _ = val.(int)
	}
	return
}
func (this *Context) GetInt64(key string) (i64 int64) {
	if val, ok := this.Get(key); ok && val != nil {
		i64, _ = val.(int64)
	}
	return
}
func (this *Context) GetUint(key string) (ui uint) {
	if val, ok := this.Get(key); ok && val != nil {
		ui, _ = val.(uint)
	}
	return
}

func (c *Context) GetUInt32(key string) (i uint32) {
	if val, ok := c.Get(key); ok && val != nil {
		i, _ = val.(uint32)
	}
	return
}

func (this *Context) GetUint64(key string) (ui64 uint64) {
	if val, ok := this.Get(key); ok && val != nil {
		ui64, _ = val.(uint64)
	}
	return
}
func (this *Context) GetFloat64(key string) (f64 float64) {
	if val, ok := this.Get(key); ok && val != nil {
		f64, _ = val.(float64)
	}
	return
}
func (this *Context) GetTime(key string) (t time.Time) {
	if val, ok := this.Get(key); ok && val != nil {
		t, _ = val.(time.Time)
	}
	return
}
func (this *Context) GetDuration(key string) (d time.Duration) {
	if val, ok := this.Get(key); ok && val != nil {
		d, _ = val.(time.Duration)
	}
	return
}
func (this *Context) GetStringSlice(key string) (ss []string) {
	if val, ok := this.Get(key); ok && val != nil {
		ss, _ = val.([]string)
	}
	return
}
func (this *Context) GetStringMap(key string) (sm map[string]interface{}) {
	if val, ok := this.Get(key); ok && val != nil {
		sm, _ = val.(map[string]interface{})
	}
	return
}
func (this *Context) GetStringMapString(key string) (sms map[string]string) {
	if val, ok := this.Get(key); ok && val != nil {
		sms, _ = val.(map[string]string)
	}
	return
}

func (this *Context) GetStringMapStringSlice(key string) (smss map[string][]string) {
	if val, ok := this.Get(key); ok && val != nil {
		smss, _ = val.(map[string][]string)
	}
	return
}

func (this *Context) GetUserId() string {
	return this.GetString("UserId")
}

func (this *Context) Header(key, value string) {
	if value == "" {
		this.Writer.Header().Del(key)
		return
	}
	this.Writer.Header().Set(key, value)
}

// Status sets the HTTP response code.
func (this *Context) Status(code int) {
	this.Writer.WriteHeader(code)
}

func (this *Context) Param(key string) string {
	return this.Params.ByName(key)
}

func (this *Context) AddParam(key, value string) {
	this.Params = append(this.Params, Param{Key: key, Value: value})
}

func (this *Context) Query(key string) (value string) {
	value, _ = this.GetQuery(key)
	return
}

func (this *Context) DefaultQuery(key, defaultValue string) string {
	if value, ok := this.GetQuery(key); ok {
		return value
	}
	return defaultValue
}

func (this *Context) GetQuery(key string) (string, bool) {
	if values, ok := this.GetQueryArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (this *Context) initQueryCache() {
	if this.queryCache == nil {
		if this.Request != nil {
			this.queryCache = this.Request.URL.Query()
		} else {
			this.queryCache = url.Values{}
		}
	}
}

func (this *Context) GetQueryArray(key string) (values []string, ok bool) {
	this.initQueryCache()
	values, ok = this.queryCache[key]
	return
}

func (this *Context) QueryMap(key string) (dicts map[string]string) {
	dicts, _ = this.GetQueryMap(key)
	return
}

func (this *Context) GetQueryMap(key string) (map[string]string, bool) {
	this.initQueryCache()
	return this.get(this.queryCache, key)
}

func (this *Context) PostForm(key string) (value string) {
	value, _ = this.GetPostForm(key)
	return
}

func (this *Context) GetPostForm(key string) (string, bool) {
	if values, ok := this.GetPostFormArray(key); ok {
		return values[0], ok
	}
	return "", false
}

func (this *Context) initFormCache() {
	if this.formCache == nil {
		this.formCache = make(url.Values)
		req := this.Request
		if err := req.ParseMultipartForm(this.engine.MaxMultipartMemory); err != nil {
			if !errors.Is(err, http.ErrNotMultipart) {
				this.Log.Errorf("error on parse multipart form array: %v", err)
			}
		}
		this.formCache = req.PostForm
	}
}

func (this *Context) GetPostFormArray(key string) (values []string, ok bool) {
	this.initFormCache()
	values, ok = this.formCache[key]
	return
}

func (this *Context) PostFormMap(key string) (dicts map[string]string) {
	dicts, _ = this.GetPostFormMap(key)
	return
}

func (this *Context) GetPostFormMap(key string) (map[string]string, bool) {
	this.initFormCache()
	return this.get(this.formCache, key)
}

func (this *Context) FormFile(name string) (multipart.File, *multipart.FileHeader, error) {
	if this.Request.MultipartForm == nil {
		if err := this.Request.ParseMultipartForm(this.engine.MaxMultipartMemory); err != nil {
			return nil, nil, err
		}
	}
	f, fh, err := this.Request.FormFile(name)
	if err != nil {
		return nil, nil, err
	}
	return f, fh, err
}

/*
	MultipartForm 是解析后的多部分表单，包括文件上传。
*/
func (this *Context) MultipartForm() (*multipart.Form, error) {
	err := this.Request.ParseMultipartForm(this.engine.MaxMultipartMemory)
	return this.Request.MultipartForm, err
}

/*
	保存上传文件
*/
func (this *Context) SaveUploadedFile(file *multipart.FileHeader, dst string) error {
	src, err := file.Open()
	if err != nil {
		return err
	}
	defer src.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, src)
	return err
}

func (this *Context) GetRawData() ([]byte, error) {
	return ioutil.ReadAll(this.Request.Body)
}

//序列化--------------------------------------------------------------------------------------------
func (this *Context) Bind(obj interface{}) error {
	b := binding.Default(this.Request.Method, this.ContentType())
	return this.MustBindWith(obj, b)
}

func (this *Context) ShouldBindJSON(obj interface{}) error {
	return this.ShouldBindWith(obj, binding.JSON)
}

func (this *Context) MustBindWith(obj interface{}, b binding.Binding) error {
	if err := this.ShouldBindWith(obj, b); err != nil {
		this.AbortWithError(http.StatusBadRequest, err).SetType(ErrorTypeBind) // nolint: errcheck
		return err
	}
	return nil
}

func (this *Context) ShouldBindWith(obj interface{}, b binding.Binding) error {
	return b.Bind(this.Request, obj)
}
func (this *Context) ShouldBindUri(obj interface{}) error {
	m := make(map[string][]string)
	for _, v := range this.Params {
		m[v.Key] = []string{v.Value}
	}
	return binding.Uri.BindUri(m, obj)
}

func (this *Context) BindJSON(obj interface{}) error {
	return this.MustBindWith(obj, binding.JSON)
}
func (this *Context) BindXML(obj interface{}) error {
	return this.MustBindWith(obj, binding.XML)
}
func (this *Context) BindQuery(obj interface{}) error {
	return this.MustBindWith(obj, binding.Query)
}

func (this *Context) BindYAML(obj interface{}) error {
	return this.MustBindWith(obj, binding.YAML)
}

func (this *Context) BindHeader(obj interface{}) error {
	return this.MustBindWith(obj, binding.Header)
}

func (this *Context) BindUri(obj interface{}) error {
	if err := this.ShouldBindUri(obj); err != nil {
		this.AbortWithError(http.StatusBadRequest, err).SetType(ErrorTypeBind) // nolint: errcheck
		return err
	}
	return nil
}

//输出-----------------------------------------------------------------------------------------
func (this *Context) HTML(code int, name string, obj interface{}) {
	instance := this.engine.HTMLRender.Instance(name, obj)
	this.Render(code, instance)
}

func (this *Context) IndentedJSON(code int, obj interface{}) {
	this.Render(code, render.IndentedJSON{Data: obj})
}

func (this *Context) SecureJSON(code int, obj interface{}) {
	this.Render(code, render.SecureJSON{Prefix: this.engine.secureJSONPrefix, Data: obj})
}

func (this *Context) JSONP(code int, obj interface{}) {
	callback := this.DefaultQuery("callback", "")
	if callback == "" {
		this.Render(code, render.JSON{Data: obj})
		return
	}
	this.Render(code, render.JsonpJSON{Callback: callback, Data: obj})
}

func (this *Context) JSON(code int, obj interface{}) {
	this.Render(code, render.JSON{Data: obj})
}

func (this *Context) AsciiJSON(code int, obj interface{}) {
	this.Render(code, render.AsciiJSON{Data: obj})
}

func (this *Context) PureJSON(code int, obj interface{}) {
	this.Render(code, render.PureJSON{Data: obj})
}

func (this *Context) XML(code int, obj interface{}) {
	this.Render(code, render.XML{Data: obj})
}

func (this *Context) YAML(code int, obj interface{}) {
	this.Render(code, render.YAML{Data: obj})
}

func (this *Context) ProtoBuf(code int, obj interface{}) {
	this.Render(code, render.ProtoBuf{Data: obj})
}

func (this *Context) String(code int, format string, values ...interface{}) {
	this.Render(code, render.String{Format: format, Data: values})
}

func (this *Context) Redirect(code int, location string) {
	this.Render(-1, render.Redirect{
		Code:     code,
		Location: location,
		Request:  this.Request,
	})
}

func (this *Context) Data(code int, contentType string, data []byte) {
	this.Render(code, render.Data{
		ContentType: contentType,
		Data:        data,
	})
}

func (this *Context) DataFromReader(code int, contentLength int64, contentType string, reader io.Reader, extraHeaders map[string]string) {
	this.Render(code, render.Reader{
		Headers:       extraHeaders,
		ContentType:   contentType,
		ContentLength: contentLength,
		Reader:        reader,
	})
}

func (this *Context) File(filepath string) {
	http.ServeFile(this.Writer, this.Request, filepath)
}

/*渲染页面接口*/
func (this *Context) Render(code int, r render.Render) {
	this.Status(code)

	if !bodyAllowedForStatus(code) {
		r.WriteContentType(this.Writer)
		this.Writer.WriteHeaderNow()
		return
	}
	if err := r.Render(this.Writer); err != nil {
		panic(err)
	}
}

func (this *Context) FileFromFS(filepath string, fs http.FileSystem) {
	defer func(old string) {
		this.Request.URL.Path = old
	}(this.Request.URL.Path)
	this.Request.URL.Path = filepath
	http.FileServer(fs).ServeHTTP(this.Writer, this.Request)
}

/*
以高效的方式将指定的文件写入正文流
在客户端，通常会使用给定的文件名下载文件
*/
func (this *Context) FileAttachment(filepath, filename string) {
	if isASCII(filename) {
		this.Writer.Header().Set("Content-Disposition", `attachment; filename="`+filename+`"`)
	} else {
		this.Writer.Header().Set("Content-Disposition", `attachment; filename*=UTF-8''`+url.QueryEscape(filename))
	}
	http.ServeFile(this.Writer, this.Request, filepath)
}

func (this *Context) Stream(step func(w io.Writer) bool) bool {
	w := this.Writer
	clientGone := w.CloseNotify()
	for {
		select {
		case <-clientGone:
			return true
		default:
			keepOpen := step(w)
			w.Flush()
			if !keepOpen {
				return false
			}
		}
	}
}

type Negotiate struct {
	Offered  []string
	HTMLName string
	HTMLData interface{}
	JSONData interface{}
	XMLData  interface{}
	YAMLData interface{}
	Data     interface{}
}

func (this *Context) Negotiate(code int, config Negotiate) {
	switch this.NegotiateFormat(config.Offered...) {
	case binding.MIMEJSON:
		data := chooseData(config.JSONData, config.Data)
		this.JSON(code, data)

	case binding.MIMEHTML:
		data := chooseData(config.HTMLData, config.Data)
		this.HTML(code, config.HTMLName, data)

	case binding.MIMEXML:
		data := chooseData(config.XMLData, config.Data)
		this.XML(code, data)

	case binding.MIMEYAML:
		data := chooseData(config.YAMLData, config.Data)
		this.YAML(code, data)

	default:
		this.AbortWithError(http.StatusNotAcceptable, errors.New("the accepted formats are not offered by the server")) // nolint: errcheck
	}
}

func (this *Context) NegotiateFormat(offered ...string) string {
	assert1(len(offered) > 0, "you must provide at least one offer")

	if this.Accepted == nil {
		this.Accepted = parseAccept(this.requestHeader("Accept"))
	}
	if len(this.Accepted) == 0 {
		return offered[0]
	}
	for _, accepted := range this.Accepted {
		for _, offer := range offered {
			// According to RFC 2616 and RFC 2396, non-ASCII characters are not allowed in headers,
			// therefore we can just iterate over the string without casting it into []rune
			i := 0
			for ; i < len(accepted); i++ {
				if accepted[i] == '*' || offer[i] == '*' {
					return offer
				}
				if accepted[i] != offer[i] {
					break
				}
			}
			if i == len(accepted) {
				return offer
			}
		}
	}
	return ""
}

func (this *Context) SetAccepted(formats ...string) {
	this.Accepted = formats
}

func (this *Context) Deadline() (deadline time.Time, ok bool) {
	if this.Request == nil || this.Request.Context() == nil {
		return
	}
	return this.Request.Context().Deadline()
}

func (this *Context) Done() <-chan struct{} {
	if this.Request == nil || this.Request.Context() == nil {
		return nil
	}
	return this.Request.Context().Done()
}

func (this *Context) Err() error {
	if this.Request == nil || this.Request.Context() == nil {
		return nil
	}
	return this.Request.Context().Err()
}

func (c *Context) Value(key interface{}) interface{} {
	if key == 0 {
		return c.Request
	}
	if keyAsString, ok := key.(string); ok {
		if val, exists := c.Get(keyAsString); exists {
			return val
		}
	}
	if c.Request == nil || c.Request.Context() == nil {
		return nil
	}
	return c.Request.Context().Value(key)
}

func (this *Context) ContentType() string {
	return filterFlags(this.requestHeader("Content-Type"))
}

func (this *Context) RemoteIP() string {
	ip, _, err := net.SplitHostPort(strings.TrimSpace(this.Request.RemoteAddr))
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

func (this *Context) IsWebsocket() bool {
	if strings.Contains(strings.ToLower(this.requestHeader("Connection")), "upgrade") &&
		strings.EqualFold(this.requestHeader("Upgrade"), "websocket") {
		return true
	}
	return false
}

func (this *Context) SetSameSite(samesite http.SameSite) {
	this.sameSite = samesite
}

func (this *Context) SetCookie(name, value string, maxAge int, path, domain string, secure, httpOnly bool) {
	if path == "" {
		path = "/"
	}
	http.SetCookie(this.Writer, &http.Cookie{
		Name:     name,
		Value:    url.QueryEscape(value),
		MaxAge:   maxAge,
		Path:     path,
		Domain:   domain,
		SameSite: this.sameSite,
		Secure:   secure,
		HttpOnly: httpOnly,
	})
}

func (this *Context) Cookie(name string) (string, error) {
	cookie, err := this.Request.Cookie(name)
	if err != nil {
		return "", err
	}
	val, _ := url.QueryUnescape(cookie.Value)
	return val, nil
}

func (this *Context) Error(err error) *Error {
	if err == nil {
		panic("err is nil")
	}

	var parsedError *Error
	ok := errors.As(err, &parsedError)
	if !ok {
		parsedError = &Error{
			Err:  err,
			Type: ErrorTypePrivate,
		}
	}

	this.Errors = append(this.Errors, parsedError)
	return parsedError
}

func (this *Context) get(m map[string][]string, key string) (map[string]string, bool) {
	dicts := make(map[string]string)
	exist := false
	for k, v := range m {
		if i := strings.IndexByte(k, '['); i >= 1 && k[0:i] == key {
			if j := strings.IndexByte(k[i+1:], ']'); j >= 1 {
				exist = true
				dicts[k[i+1:][:j]] = v[0]
			}
		}
	}
	return dicts, exist
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
	this.Accepted = nil
	this.queryCache = nil
	this.formCache = nil
	this.sameSite = 0
	*this.params = (*this.params)[:0]
	*this.skippedNodes = (*this.skippedNodes)[:0]
}

/*
	bodyAllowedForStatus 是 http.bodyAllowedForStatus 非导出函数的副本。
*/
func bodyAllowedForStatus(status int) bool {
	switch {
	case status >= 100 && status <= 199:
		return false
	case status == http.StatusNoContent:
		return false
	case status == http.StatusNotModified:
		return false
	}
	return true
}
