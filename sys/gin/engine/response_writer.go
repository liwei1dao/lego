package engine

import (
	"bufio"
	"io"
	"net"
	"net/http"

	"github.com/liwei1dao/lego/sys/log"
)

const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

type IResponseWriter interface {
	http.ResponseWriter
	http.Hijacker
	http.Flusher
	http.CloseNotifier
	/*
		返回当前请求的 HTTP 响应状态码。
	*/
	Status() int
	/*
		大小返回已经写入响应 http 正文的字节数
	*/
	Size() int
	// WriteString writes the string into the response body.
	WriteString(string) (int, error)
	// Written returns true if the response body was already written.
	Written() bool
	/*
		WriteHeaderNow 强制写入 http 标头（状态码 + 标头）。
	*/
	WriteHeaderNow()
	// Pusher get the http.Pusher for server push
	Pusher() http.Pusher
}
type ResponseWriter struct {
	http.ResponseWriter
	log    log.ILogger
	size   int
	status int
}

func (this *ResponseWriter) reset(writer http.ResponseWriter) {
	this.ResponseWriter = writer
	this.size = noWritten
	this.status = defaultStatus
}

func (this *ResponseWriter) WriteHeader(code int) {
	if code > 0 && this.status != code {
		if this.Written() {
			this.log.Warnf("Headers were already written. Wanted to override status code %d with %d", this.status, code)
		}
		this.status = code
	}
}
func (this *ResponseWriter) WriteHeaderNow() {
	if !this.Written() {
		this.size = 0
		this.ResponseWriter.WriteHeader(this.status)
	}
}

func (this *ResponseWriter) Write(data []byte) (n int, err error) {
	this.WriteHeaderNow()
	n, err = this.ResponseWriter.Write(data)
	this.size += n
	return
}

func (this *ResponseWriter) WriteString(s string) (n int, err error) {
	this.WriteHeaderNow()
	n, err = io.WriteString(this.ResponseWriter, s)
	this.size += n
	return
}

func (this *ResponseWriter) Status() int {
	return this.status
}

func (this *ResponseWriter) Size() int {
	return this.size
}
func (this *ResponseWriter) Written() bool {
	return this.size != noWritten
}

// Hijack implements the http.Hijacker interface.
func (this *ResponseWriter) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	if this.size < 0 {
		this.size = 0
	}
	return this.ResponseWriter.(http.Hijacker).Hijack()
}

func (this *ResponseWriter) CloseNotify() <-chan bool {
	return this.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (this *ResponseWriter) Flush() {
	this.WriteHeaderNow()
	this.ResponseWriter.(http.Flusher).Flush()
}

func (this *ResponseWriter) Pusher() (pusher http.Pusher) {
	if pusher, ok := this.ResponseWriter.(http.Pusher); ok {
		return pusher
	}
	return nil
}
