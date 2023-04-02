package gin

import (
	"fmt"
	"net/http"
	"reflect"
	"sort"
	"strings"

	"github.com/liwei1dao/lego/sys/gin/engine"
	"github.com/liwei1dao/lego/utils/crypto/md5"
)

type ISys interface {
	engine.IRoutes
	HandleContext(c *engine.Context)
	LoadHTMLGlob(pattern string)
	Close() (err error)
}

var defsys ISys

func OnInit(config map[string]interface{}, opt ...Option) (err error) {
	var option *Options
	if option, err = newOptions(config, opt...); err != nil {
		return
	}
	defsys, err = newSys(option)
	return
}

func NewSys(opt ...Option) (sys ISys, err error) {
	var option *Options
	if option, err = newOptionsByOption(opt...); err != nil {
		return
	}
	sys, err = newSys(option)
	return
}

func LoadHTMLGlob(pattern string) {
	defsys.LoadHTMLGlob(pattern)
}

func HandleContext(c *engine.Context) {
	defsys.HandleContext(c)
}

func Close() (err error) {
	return defsys.Close()
}
func NoRoute(handlers ...engine.HandlerFunc) {
	defsys.NoRoute(handlers...)
}
func Use(handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.Use(handlers...)
}
func Handle(httpMethod string, relativePath string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.Handle(httpMethod, relativePath, handlers...)
}
func Any(relativePath string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.Any(relativePath, handlers...)
}
func GET(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.GET(httpMethod, handlers...)
}
func POST(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.POST(httpMethod, handlers...)
}
func DELETE(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.DELETE(httpMethod, handlers...)
}
func PATCH(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.PATCH(httpMethod, handlers...)
}
func PUT(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.PUT(httpMethod, handlers...)
}
func OPTIONS(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.OPTIONS(httpMethod, handlers...)
}
func HEAD(httpMethod string, handlers ...engine.HandlerFunc) engine.IRoutes {
	return defsys.HEAD(httpMethod, handlers...)
}
func StaticFile(relativePath string, filepath string) engine.IRoutes {
	return defsys.StaticFile(relativePath, filepath)
}
func StaticFileFS(relativePath string, filepath string, fs http.FileSystem) engine.IRoutes {
	return defsys.StaticFileFS(relativePath, filepath, fs)
}

func Static(relativePath string, root string) engine.IRoutes {
	return defsys.Static(relativePath, root)
}

func StaticFS(relativePath string, fs http.FileSystem) engine.IRoutes {
	return defsys.StaticFS(relativePath, fs)
}

//签名接口
func ParamSign(key string, param map[string]interface{}) (origin, sign string) {
	var keys []string
	for k, _ := range param {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	builder := strings.Builder{}
	for _, v := range keys {
		builder.WriteString(v)
		builder.WriteString("=")
		switch reflect.TypeOf(param[v]).Kind() {
		case reflect.Int,
			reflect.Int8,
			reflect.Int16,
			reflect.Int32,
			reflect.Int64,
			reflect.Uint,
			reflect.Uint8,
			reflect.Uint16,
			reflect.Uint32,
			reflect.Uint64:
			builder.WriteString(fmt.Sprintf("%d", param[v]))
			break
		case reflect.Float32,
			reflect.Float64:
			builder.WriteString(fmt.Sprintf("%v", param[v]))
		case reflect.Bool:
			builder.WriteString(fmt.Sprintf("%v", param[v]))
		case reflect.Slice, reflect.Array:
			s := reflect.ValueOf(param[v])
			valueStr := ""
			for i := 0; i < s.Len(); i++ {
				valueStr += fmt.Sprintf("%v,", s.Index(i).Interface())
			}
			if s.Len() > 0 {
				valueStr = valueStr[0 : len(valueStr)-1]
			}
			builder.WriteString(fmt.Sprintf("%s", valueStr))
			break
		default:
			builder.WriteString(fmt.Sprintf("%s", param[v]))
			break
		}
		builder.WriteString("&")
	}
	builder.WriteString("key=" + key)
	origin = builder.String()
	sign = md5.MD5EncToLower(origin)
	return
}
