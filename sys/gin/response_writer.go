package gin

import (
	"net/http"

	"github.com/liwei1dao/lego/sys/log"
)

const (
	noWritten     = -1
	defaultStatus = http.StatusOK
)

type ResponseWriter interface {
	http.ResponseWriter
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
	/*
		WriteHeaderNow 强制写入 http 标头（状态码 + 标头）。
	*/
	WriteHeaderNow()
}
type responseWriter struct {
	http.ResponseWriter
	size   int
	status int
}

func (this *responseWriter) reset(writer http.ResponseWriter) {
	this.ResponseWriter = writer
	this.size = noWritten
	this.status = defaultStatus
}

func (w *responseWriter) WriteHeader(code int) {
	if code > 0 && w.status != code {
		if w.Written() {
			log.Warnf("Headers were already written. Wanted to override status code %d with %d", w.status, code)
		}
		w.status = code
	}
}

func (w *responseWriter) Status() int {
	return w.status
}

func (w *responseWriter) Size() int {
	return w.size
}
func (w *responseWriter) Written() bool {
	return w.size != noWritten
}

func (w *responseWriter) CloseNotify() <-chan bool {
	return w.ResponseWriter.(http.CloseNotifier).CloseNotify()
}

func (w *responseWriter) Flush() {
	w.WriteHeaderNow()
	w.ResponseWriter.(http.Flusher).Flush()
}

func (this *responseWriter) WriteHeaderNow() {
	if !this.Written() {
		this.size = 0
		this.ResponseWriter.WriteHeader(this.status)
	}
}
