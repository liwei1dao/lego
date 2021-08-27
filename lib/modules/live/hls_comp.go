package live

import (
	"fmt"
	"net"
	"net/http"
	"path"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib/modules/live/av"
	"github.com/liwei1dao/lego/sys/log"
)

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

//主机信息监控
type HlsComp struct {
	cbase.ModuleCompBase
	options IOptions
	module  ILive
	listen  net.Listener
	conns   *sync.Map
}

func (this *HlsComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.options = options.(IOptions)
	this.module = module.(ILive)
	this.conns = &sync.Map{}
	go this.checkStop()
	return
}

func (this *HlsComp) Start() (err error) {
	err = this.ModuleCompBase.Start()
	if this.listen, err = net.Listen("tcp", this.options.GetHlsAddr()); err == nil {
		go this.run()
	}
	return
}

func (this *HlsComp) checkStop() {
	for {
		<-time.After(5 * time.Second)

		this.conns.Range(func(key, val interface{}) bool {
			v := val.(*Source)
			if !v.Alive() && !this.options.GetHlsKeepAfterEnd() {
				log.Debugf("check stop and remove:%v", v.Info())
				this.conns.Delete(key)
			}
			return true
		})
	}
}

func (this *HlsComp) run() (err error) {
	defer lego.Recover()
	log.Infof("HLS listen On %s", this.options.GetHlsAddr())
	this.Serve()
	return
}

func (this *HlsComp) GetWriter(info av.Info) av.WriteCloser {
	var s *Source
	v, ok := this.conns.Load(info.Key)
	if !ok {
		log.Debug("new hls source")
		s = NewSource(this.options, info)
		this.conns.Store(info.Key, s)
	} else {
		s = v.(*Source)
	}
	return s
}

func (this *HlsComp) Serve() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		this.handle(w, r)
	})
	http.Serve(this.listen, mux)
	return nil
}

func (this *HlsComp) handle(w http.ResponseWriter, r *http.Request) {
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
			log.Debugf("GenM3U8PlayList error:%v", err)
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
			log.Debugf("GetItem error: %v", err)
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Content-Type", "video/mp2ts")
		w.Header().Set("Content-Length", strconv.Itoa(len(item.Data)))
		w.Write(item.Data)
	}
}

func (this *HlsComp) parseM3u8(pathstr string) (key string, err error) {
	pathstr = strings.TrimLeft(pathstr, "/")
	key = strings.Split(pathstr, path.Ext(pathstr))[0]
	return
}

func (this *HlsComp) getConn(key string) *Source {
	v, ok := this.conns.Load(key)
	if !ok {
		return nil
	}
	return v.(*Source)
}

func (this *HlsComp) parseTs(pathstr string) (key string, err error) {
	pathstr = strings.TrimLeft(pathstr, "/")
	paths := strings.SplitN(pathstr, "/", 3)
	if len(paths) != 3 {
		err = fmt.Errorf("invalid path=%s", pathstr)
		return
	}
	key = paths[0] + "/" + paths[1]

	return
}
