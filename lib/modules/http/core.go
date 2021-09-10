package http

import (
	"encoding/xml"
	"fmt"
	"math"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
)

type (
	HandlerFunc   func(*Context)
	HandlersChain []HandlerFunc

	IHttp interface {
		core.IModule
		IRoutes
		addRoute(method, path string, handlers HandlersChain) (err error)
		Group(relativePath string, handlers ...HandlerFunc) *RouterGroup
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
	}
	H       map[string]interface{}
	OutJson struct {
		ErrorCode core.ErrorCode `json:"code"`
		Message   string         `json:"message"`
		Data      interface{}    `json:"data"`
	}
)

func (h H) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	start.Name = xml.Name{
		Space: "",
		Local: "map",
	}
	if err := e.EncodeToken(start); err != nil {
		return err
	}
	for key, value := range h {
		elem := xml.StartElement{
			Name: xml.Name{Space: "", Local: key},
			Attr: []xml.Attr{},
		}
		if err := e.EncodeElement(value, elem); err != nil {
			return err
		}
	}

	return e.EncodeToken(xml.EndElement{Name: start.Name})
}

const (
	abortIndex             int8 = math.MaxInt8 / 2 //HandlersChain 处理方法最大容量
	defaultMultipartMemory      = 32 << 20         // 32 MB
)

var (
	default404Body   = []byte("404 page not found")
	default405Body   = []byte("405 method not allowed")
	defaultAppEngine bool
	mimePlain        = []string{MIMEPlain}
)

func outErr(guard bool, text string) (err error) {
	if !guard {
		return fmt.Errorf(text)
	}
	return err
}

func cleanPath(p string) string {
	// Turn empty string into "/"
	if p == "" {
		return "/"
	}

	n := len(p)
	var buf []byte

	// Invariants:
	//      reading from path; r is index of next byte to process.
	//      writing to buf; w is index of next byte to write.

	// path must start with '/'
	r := 1
	w := 1

	if p[0] != '/' {
		r = 0
		buf = make([]byte, n+1)
		buf[0] = '/'
	}

	trailing := n > 1 && p[n-1] == '/'

	// A bit more clunky without a 'lazybuf' like the path package, but the loop
	// gets completely inlined (bufApp). So in contrast to the path package this
	// loop has no expensive function calls (except 1x make)

	for r < n {
		switch {
		case p[r] == '/':
			// empty path element, trailing slash is added after the end
			r++

		case p[r] == '.' && r+1 == n:
			trailing = true
			r++

		case p[r] == '.' && p[r+1] == '/':
			// . element
			r += 2

		case p[r] == '.' && p[r+1] == '.' && (r+2 == n || p[r+2] == '/'):
			// .. element: remove to last /
			r += 3

			if w > 1 {
				// can backtrack
				w--

				if buf == nil {
					for w > 1 && p[w] != '/' {
						w--
					}
				} else {
					for w > 1 && buf[w] != '/' {
						w--
					}
				}
			}

		default:
			// real path element.
			// add slash if needed
			if w > 1 {
				bufApp(&buf, p, w, '/')
				w++
			}

			// copy element
			for r < n && p[r] != '/' {
				bufApp(&buf, p, w, p[r])
				w++
				r++
			}
		}
	}

	// re-append trailing slash
	if trailing && w > 1 {
		bufApp(&buf, p, w, '/')
		w++
	}

	if buf == nil {
		return p[:w]
	}
	return string(buf[:w])
}
func bufApp(buf *[]byte, s string, w int, c byte) {
	if *buf == nil {
		if s[w] == c {
			return
		}

		*buf = make([]byte, len(s))
		copy(*buf, s[:w])
	}
	(*buf)[w] = c
}
func serveError(c *Context, code int, defaultMessage []byte) {
	c.writermem.status = code
	c.Next()
	if c.writermem.Written() {
		return
	}
	if c.writermem.Status() == code {
		c.writermem.Header()["Content-Type"] = mimePlain
		_, err := c.Writer.Write(defaultMessage)
		if err != nil {
			log.Errorf("cannot write message to writer during serve error: %v", err)
		}
		return
	}
	c.writermem.WriteHeaderNow()
	return
}
