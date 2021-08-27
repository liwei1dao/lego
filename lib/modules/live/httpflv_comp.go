package live

import (
	"encoding/json"
	"net"
	"net/http"
	"strings"

	"github.com/liwei1dao/lego"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/log"
)

type httpflvstream struct {
	Key string `json:"key"`
	Id  string `json:"id"`
}

type httpflvstreams struct {
	Publishers []httpflvstream `json:"publishers"`
	Players    []httpflvstream `json:"players"`
}

type HttpFlvComp struct {
	cbase.ModuleCompBase
	options IOptions
	module  ILive
	listen  net.Listener
}

func (this *HttpFlvComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.options = options.(IOptions)
	this.module = module.(ILive)
	return
}

func (this *HttpFlvComp) Start() (err error) {
	err = this.ModuleCompBase.Start()
	if this.listen, err = net.Listen("tcp", this.options.GetHttpFlvAddr()); err == nil {
		go this.run()
	}
	return
}

func (this *HttpFlvComp) run() (err error) {
	defer lego.Recover()
	log.Infof("HTTP-FLV listen On %s", this.options.GetHttpFlvAddr())
	this.Serve()
	return
}

func (this *HttpFlvComp) Serve() error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		this.handleConn(w, r)
	})
	mux.HandleFunc("/streams", func(w http.ResponseWriter, r *http.Request) {
		this.getStream(w, r)
	})
	if err := http.Serve(this.listen, mux); err != nil {
		return err
	}
	return nil
}

// 获取发布和播放器的信息
func (this *HttpFlvComp) getStreams(w http.ResponseWriter, r *http.Request) *httpflvstreams {
	rtmpStream := this.module.GetHandler().(*RtmpStream)
	if rtmpStream == nil {
		return nil
	}
	msgs := new(httpflvstreams)

	rtmpStream.GetStreams().Range(func(key, val interface{}) bool {
		if s, ok := val.(*Stream); ok {
			if s.GetReader() != nil {
				msg := httpflvstream{key.(string), s.GetReader().Info().UID}
				msgs.Publishers = append(msgs.Publishers, msg)
			}
		}
		return true
	})

	rtmpStream.GetStreams().Range(func(key, val interface{}) bool {
		ws := val.(*Stream).GetWs()

		ws.Range(func(k, v interface{}) bool {
			if pw, ok := v.(*PackWriterCloser); ok {
				if pw.GetWriter() != nil {
					msg := httpflvstream{key.(string), pw.GetWriter().Info().UID}
					msgs.Players = append(msgs.Players, msg)
				}
			}
			return true
		})
		return true
	})

	return msgs
}

func (this *HttpFlvComp) getStream(w http.ResponseWriter, r *http.Request) {
	msgs := this.getStreams(w, r)
	if msgs == nil {
		return
	}
	resp, _ := json.Marshal(msgs)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

func (this *HttpFlvComp) handleConn(w http.ResponseWriter, r *http.Request) {
	defer lego.Recover()

	url := r.URL.String()
	u := r.URL.Path
	if pos := strings.LastIndex(u, "."); pos < 0 || u[pos:] != ".flv" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	path := strings.TrimSuffix(strings.TrimLeft(u, "/"), ".flv")
	paths := strings.SplitN(path, "/", 2)
	log.Debugf("url:%s path:%s paths:%v", u, path, paths)

	if len(paths) != 2 {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}

	// 判断视屏流是否发布,如果没有发布,直接返回404
	msgs := this.getStreams(w, r)
	if msgs == nil || len(msgs.Publishers) == 0 {
		http.Error(w, "invalid path", http.StatusNotFound)
		return
	} else {
		include := false
		for _, item := range msgs.Publishers {
			if item.Key == path {
				include = true
				break
			}
		}
		if include == false {
			http.Error(w, "invalid path", http.StatusNotFound)
			return
		}
	}

	w.Header().Set("Access-Control-Allow-Origin", "*")
	writer := NewFLVWriter(paths[0], paths[1], url, w)

	this.module.GetHandler().HandleWriter(writer)
	writer.Wait()
}
