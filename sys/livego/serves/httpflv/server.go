package httpflv

import (
	"net"
	"net/http"
	"strings"

	"github.com/liwei1dao/lego/utils/codec/json"

	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/log"
)

type stream struct {
	Key string `json:"key"`
	Id  string `json:"id"`
}

type streams struct {
	Publishers []stream `json:"publishers"`
	Players    []stream `json:"players"`
}

func NewServer(sys core.ISys, log log.ILogger) (server *Server, err error) {
	server = &Server{
		sys: sys,
	}
	err = server.init()
	return
}

type Server struct {
	sys core.ISys
	log log.ILogger
}

func (this *Server) init() (err error) {
	var flvListen net.Listener
	if flvListen, err = net.Listen("tcp", this.sys.GetHTTPFLVAddr()); err != nil {
		this.log.Errorf("HttpFlvServer init err%v", err)
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				this.log.Errorf("HTTP-FLV server panic: ", r)
			}
		}()
		this.log.Infof("HTTP-FLV listen On ", this.sys.GetHTTPFLVAddr())
		this.Serve(flvListen)
	}()
	return
}

func (this *Server) Serve(l net.Listener) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		this.handleConn(w, r)
	})
	mux.HandleFunc("/streams", func(w http.ResponseWriter, r *http.Request) {
		this.getStream(w, r)
	})
	if err := http.Serve(l, mux); err != nil {
		return err
	}
	return nil
}

func (this *Server) handleConn(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			this.log.Errorf("http flv handleConn panic:%v", r)
		}
	}()

	url := r.URL.String()
	u := r.URL.Path
	if pos := strings.LastIndex(u, "."); pos < 0 || u[pos:] != ".flv" {
		http.Error(w, "invalid path", http.StatusBadRequest)
		return
	}
	path := strings.TrimSuffix(strings.TrimLeft(u, "/"), ".flv")
	paths := strings.SplitN(path, "/", 2)
	this.log.Debugf("url:%s path:%s paths:%s", u, path, paths)

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
	writer := NewFLVWriter(this.sys, this.log, paths[0], paths[1], url, w)

	this.sys.HandleWriter(writer)
	writer.Wait()
}

func (this *Server) getStream(w http.ResponseWriter, r *http.Request) {
	msgs := this.getStreams(w, r)
	if msgs == nil {
		return
	}
	resp, _ := json.Marshal(msgs)
	w.Header().Set("Content-Type", "application/json")
	w.Write(resp)
}

// 获取发布和播放器的信息
func (this *Server) getStreams(w http.ResponseWriter, r *http.Request) *streams {
	msgs := new(streams)
	this.sys.GetStreams().Range(func(key, val interface{}) bool {
		if s, ok := val.(*core.Stream); ok {
			if s.GetReader() != nil {
				msg := stream{key.(string), s.GetReader().Info().UID}
				msgs.Publishers = append(msgs.Publishers, msg)
			}
		}
		return true
	})

	this.sys.GetStreams().Range(func(key, val interface{}) bool {
		ws := val.(*core.Stream).GetWs()

		ws.Range(func(k, v interface{}) bool {
			if pw, ok := v.(*core.PackWriterCloser); ok {
				if pw.GetWriter() != nil {
					msg := stream{key.(string), pw.GetWriter().Info().UID}
					msgs.Players = append(msgs.Players, msg)
				}
			}
			return true
		})
		return true
	})

	return msgs
}
