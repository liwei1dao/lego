package hls

import (
	"fmt"
	"net"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/log"
)

func NewServer(engine core.IEngine) *Server {
	ret := &Server{
		engine: engine,
		conns:  &sync.Map{},
	}
	go ret.checkStop()
	return ret
}

var (
	ErrNoPublisher         = fmt.Errorf("no publisher")
	ErrInvalidReq          = fmt.Errorf("invalid req url path")
	ErrNoSupportVideoCodec = fmt.Errorf("no support video codec")
	ErrNoSupportAudioCodec = fmt.Errorf("no support audio codec")
)

var crossdomainxml = []byte(`<?xml version="1.0" ?>
<cross-domain-policy>
	<allow-access-from domain="*" />
	<allow-http-request-headers-from domain="*" headers="*"/>
</cross-domain-policy>`)

type Server struct {
	engine   core.IEngine
	listener net.Listener
	conns    *sync.Map
}

func (this *Server) Serve(listener net.Listener) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		this.handle(w, r)
	})
	this.listener = listener

	if this.engine.UseHlsHttps() {
		http.ServeTLS(listener, mux, "server.crt", "server.key")
	} else {
		http.Serve(listener, mux)
	}

	return nil
}

func (this *Server) handle(w http.ResponseWriter, r *http.Request) {
	if path.Base(r.URL.Path) == "crossdomain.xml" {
		w.Header().Set("Content-Type", "application/xml")
		w.Write(crossdomainxml)
		return
	}
	switch path.Ext(r.URL.Path) {
	case ".m3u8":
		key, _ := this.parseM3u8(r.URL.Path)
		conn := this.getConn(key)
		if conn == nil {
			http.Error(w, ErrNoPublisher.Error(), http.StatusForbidden)
			return
		}
		tsCache := conn.GetCacheInc()
		if tsCache == nil {
			http.Error(w, ErrNoPublisher.Error(), http.StatusForbidden)
			return
		}
		body, err := tsCache.GenM3U8PlayList()
		if err != nil {
			log.Errorf("GenM3U8PlayList error:%v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Cache-Control", "no-cache")
		w.Header().Set("Content-Type", "application/x-mpegURL")
		w.Header().Set("Content-Length", strconv.Itoa(len(body)))
		w.Write(body)
	case ".ts":
		key, _ := this.parseTs(r.URL.Path)
		conn := this.getConn(key)
		if conn == nil {
			http.Error(w, ErrNoPublisher.Error(), http.StatusForbidden)
			return
		}
		tsCache := conn.GetCacheInc()
		item, err := tsCache.GetItem(r.URL.Path)
		if err != nil {
			log.Errorf("GetItem error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "video/mp2ts")
		w.Header().Set("Content-Length", strconv.Itoa(len(item.Data)))
		w.Write(item.Data)
	}
}

func (this *Server) parseM3u8(pathstr string) (key string, err error) {
	pathstr = strings.TrimLeft(pathstr, "/")
	key = strings.Split(pathstr, path.Ext(pathstr))[0]
	return
}

func (this *Server) getConn(key string) *Source {
	v, ok := this.conns.Load(key)
	if !ok {
		return nil
	}
	return v.(*Source)
}

func (this *Server) parseTs(pathstr string) (key string, err error) {
	pathstr = strings.TrimLeft(pathstr, "/")
	paths := strings.SplitN(pathstr, "/", 3)
	if len(paths) != 3 {
		err = fmt.Errorf("invalid path=%s", pathstr)
		return
	}
	key = paths[0] + "/" + paths[1]
	return
}

func (this *Server) checkStop() {
	for {
		<-time.After(5 * time.Second)

		this.conns.Range(func(key, val interface{}) bool {
			v := val.(*Source)
			if !v.Alive() && !this.engine.HLSKeepAfterEnd() {
				log.Debugf("[SYS LiveGo] check stop and remove:%v", v.Info())
				this.conns.Delete(key)
			}
			return true
		})
	}
}
