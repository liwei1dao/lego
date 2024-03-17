package api

import (
	"errors"
	"fmt"
	"net"
	"net/http"

	"github.com/liwei1dao/lego/utils/codec/json"

	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/livego/serves/rtmp"
	"github.com/liwei1dao/lego/sys/log"
)

type Response struct {
	w      http.ResponseWriter
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func (this *Response) SendJson() (int, error) {
	resp, _ := json.Marshal(this)
	this.w.Header().Set("Content-Type", "application/json")
	this.w.WriteHeader(this.Status)
	return this.w.Write(resp)
}

func NewServer(sys core.ISys, log log.ILogger) (server *Server, err error) {
	server = &Server{
		sys: sys,
		log: log,
	}
	err = server.init()
	return
}

type Server struct {
	sys     core.ISys
	log     log.ILogger
	session map[string]*core.RtmpRelay
}

func (this *Server) init() (err error) {
	var (
		opListen net.Listener
	)
	if this.sys.GetApiAddr() != "" {
		opListen, err = net.Listen("tcp", this.sys.GetApiAddr())
		if err != nil {
			this.log.Errorf("init ApiServer err:%v", err)
			return
		}
		go func() {
			defer func() {
				if r := recover(); r != nil {
					this.log.Errorf("HTTP-API server panic: ", r)
				}
			}()
			this.log.Infof("HTTP-API listen On ", this.sys.GetApiAddr())
			this.Serve(opListen)
		}()
	} else {
		err = errors.New("HTTP-API ApiAddr is null")
	}
	return
}

func (this *Server) Serve(l net.Listener) (err error) {
	mux := http.NewServeMux()
	mux.Handle("/statics/", http.StripPrefix("/statics/", http.FileServer(http.Dir("statics"))))
	mux.HandleFunc("/control/push", func(w http.ResponseWriter, r *http.Request) {
		this.handlePush(w, r)
	})
	mux.HandleFunc("/control/pull", func(w http.ResponseWriter, r *http.Request) {
		this.handlePull(w, r)
	})
	mux.HandleFunc("/control/get", func(w http.ResponseWriter, r *http.Request) {
		this.handleGet(w, r)
	})
	mux.HandleFunc("/control/reset", func(w http.ResponseWriter, r *http.Request) {
		this.handleReset(w, r)
	})
	mux.HandleFunc("/control/delete", func(w http.ResponseWriter, r *http.Request) {
		this.handleDelete(w, r)
	})
	mux.HandleFunc("/stat/livestat", func(w http.ResponseWriter, r *http.Request) {
		this.GetLiveStatics(w, r)
	})
	http.Serve(l, mux)
	return
}

//http://127.0.0.1:8090/control/push?&oper=start&app=live&name=123456&url=rtmp://192.168.16.136/live/123456
func (this *Server) handlePush(w http.ResponseWriter, req *http.Request) {
	var retString string
	var err error

	res := &Response{
		w:      w,
		Data:   nil,
		Status: 200,
	}

	defer res.SendJson()

	if req.ParseForm() != nil {
		res.Data = "url: /control/push?&oper=start&app=live&name=123456&url=rtmp://192.168.16.136/live/123456"
		return
	}

	oper := req.Form.Get("oper")
	app := req.Form.Get("app")
	name := req.Form.Get("name")
	url := req.Form.Get("url")

	this.log.Debugf("control push: oper=%v, app=%v, name=%v, url=%v", oper, app, name, url)
	if (len(app) <= 0) || (len(name) <= 0) || (len(url) <= 0) {
		res.Data = "control push parameter error, please check them."
		return
	}

	localurl := "rtmp://127.0.0.1" + this.sys.GetRTMPAddr() + "/" + app + "/" + name
	remoteurl := url

	keyString := "push:" + app + "/" + name
	if oper == "stop" {
		pushRtmprelay, found := this.session[keyString]
		if !found {
			retString = fmt.Sprintf("<h1>session key[%s] not exist, please check it again.</h1>", keyString)
			res.Data = retString
			return
		}
		this.log.Debugf("rtmprelay stop push %s from %s", remoteurl, localurl)
		pushRtmprelay.Stop()

		delete(this.session, keyString)
		retString = fmt.Sprintf("<h1>push url stop %s ok</h1></br>", url)
		res.Data = retString
		this.log.Debugf("push stop return %s", retString)
	} else {
		pushRtmprelay := core.NewRtmpRelay(this.sys, this.log, localurl, remoteurl)
		this.log.Debugf("rtmprelay start push %s from %s", remoteurl, localurl)
		err = pushRtmprelay.Start()
		if err != nil {
			retString = fmt.Sprintf("push error=%v", err)
		} else {
			retString = fmt.Sprintf("<h1>push url start %s ok</h1></br>", url)
			this.session[keyString] = pushRtmprelay
		}

		res.Data = retString
		this.log.Debugf("push start return %s", retString)
	}
}

//http://127.0.0.1:8090/control/pull?&oper=start&app=live&name=123456&url=rtmp://192.168.16.136/live/123456
func (this *Server) handlePull(w http.ResponseWriter, req *http.Request) {
	var retString string
	var err error

	res := &Response{
		w:      w,
		Data:   nil,
		Status: 200,
	}

	defer res.SendJson()

	if req.ParseForm() != nil {
		res.Status = 400
		res.Data = "url: /control/pull?&oper=start&app=live&name=123456&url=rtmp://192.168.16.136/live/123456"
		return
	}

	oper := req.Form.Get("oper")
	app := req.Form.Get("app")
	name := req.Form.Get("name")
	url := req.Form.Get("url")

	this.log.Debugf("control pull: oper=%v, app=%v, name=%v, url=%v", oper, app, name, url)
	if (len(app) <= 0) || (len(name) <= 0) || (len(url) <= 0) {
		res.Status = 400
		res.Data = "control push parameter error, please check them."
		return
	}

	remoteurl := "rtmp://127.0.0.1" + this.sys.GetRTMPAddr() + "/" + app + "/" + name
	localurl := url

	keyString := "pull:" + app + "/" + name
	if oper == "stop" {
		pullRtmprelay, found := this.session[keyString]

		if !found {
			retString = fmt.Sprintf("session key[%s] not exist, please check it again.", keyString)
			res.Status = 400
			res.Data = retString
			return
		}
		this.log.Debugf("rtmprelay stop push %s from %s", remoteurl, localurl)
		pullRtmprelay.Stop()

		delete(this.session, keyString)
		retString = fmt.Sprintf("<h1>push url stop %s ok</h1></br>", url)
		res.Status = 400
		res.Data = retString
		this.log.Debugf("pull stop return %s", retString)
	} else {
		pullRtmprelay := core.NewRtmpRelay(this.sys, this.log, localurl, remoteurl)
		this.log.Debugf("rtmprelay start push %s from %s", remoteurl, localurl)
		err = pullRtmprelay.Start()
		if err != nil {
			res.Status = 400
			retString = fmt.Sprintf("push error=%v", err)
		} else {
			this.session[keyString] = pullRtmprelay
			retString = fmt.Sprintf("<h1>pull url start %s ok</h1></br>", url)
		}

		res.Data = retString
		this.log.Debugf("pull start return %s", retString)
	}
}

//http://127.0.0.1:8090/control/get?room=ROOM_NAME
func (this *Server) handleGet(w http.ResponseWriter, r *http.Request) {
	res := &Response{
		w:      w,
		Data:   nil,
		Status: 200,
	}
	defer res.SendJson()

	if err := r.ParseForm(); err != nil {
		res.Status = 400
		res.Data = "url: /control/get?room=<ROOM_NAME>"
		return
	}

	room := r.Form.Get("room")

	if len(room) == 0 {
		res.Status = 400
		res.Data = "url: /control/get?room=<ROOM_NAME>"
		return
	}

	msg, err := this.sys.GetKey(room)
	if err != nil {
		msg = err.Error()
		res.Status = 400
	}
	res.Data = msg
}

//http://127.0.0.1:8090/control/reset?room=ROOM_NAME
func (this *Server) handleReset(w http.ResponseWriter, r *http.Request) {
	res := &Response{
		w:      w,
		Data:   nil,
		Status: 200,
	}
	defer res.SendJson()

	if err := r.ParseForm(); err != nil {
		res.Status = 400
		res.Data = "url: /control/reset?room=<ROOM_NAME>"
		return
	}
	room := r.Form.Get("room")

	if len(room) == 0 {
		res.Status = 400
		res.Data = "url: /control/reset?room=<ROOM_NAME>"
		return
	}

	msg, err := this.sys.SetKey(room)

	if err != nil {
		msg = err.Error()
		res.Status = 400
	}

	res.Data = msg
}

//http://127.0.0.1:8090/control/delete?room=ROOM_NAME
func (this *Server) handleDelete(w http.ResponseWriter, r *http.Request) {
	res := &Response{
		w:      w,
		Data:   nil,
		Status: 200,
	}
	defer res.SendJson()

	if err := r.ParseForm(); err != nil {
		res.Status = 400
		res.Data = "url: /control/delete?room=<ROOM_NAME>"
		return
	}

	room := r.Form.Get("room")

	if len(room) == 0 {
		res.Status = 400
		res.Data = "url: /control/delete?room=<ROOM_NAME>"
		return
	}

	if this.sys.DeleteChannel(room) {
		res.Data = "Ok"
		return
	}
	res.Status = 404
	res.Data = "room not found"
}

//http://127.0.0.1:8090/stat/livestat
func (this *Server) GetLiveStatics(w http.ResponseWriter, req *http.Request) {
	res := &Response{
		w:      w,
		Data:   nil,
		Status: 200,
	}

	defer res.SendJson()

	room := ""
	if err := req.ParseForm(); err == nil {
		room = req.Form.Get("room")
	}

	msgs := new(streams)
	if room == "" {
		this.sys.GetStreams().Range(func(key, val interface{}) bool {
			if s, ok := val.(*core.Stream); ok {
				if s.GetReader() != nil {
					switch s.GetReader().(type) {
					case *rtmp.Reader:
						v := s.GetReader().(*rtmp.Reader)
						msg := stream{key.(string), v.Info().URL, v.ReadBWInfo.StreamId, v.ReadBWInfo.VideoDatainBytes, v.ReadBWInfo.VideoSpeedInBytesperMS,
							v.ReadBWInfo.AudioDatainBytes, v.ReadBWInfo.AudioSpeedInBytesperMS}
						msgs.Publishers = append(msgs.Publishers, msg)
					}
				}
			}
			return true
		})
		this.sys.GetStreams().Range(func(key, val interface{}) bool {
			ws := val.(*core.Stream).GetWs()
			ws.Range(func(k, v interface{}) bool {
				if pw, ok := v.(*core.PackWriterCloser); ok {
					if pw.GetWriter() != nil {
						switch pw.GetWriter().(type) {
						case *rtmp.Writer:
							v := pw.GetWriter().(*rtmp.Writer)
							msg := stream{key.(string), v.Info().URL, v.WriteBWInfo.StreamId, v.WriteBWInfo.VideoDatainBytes, v.WriteBWInfo.VideoSpeedInBytesperMS,
								v.WriteBWInfo.AudioDatainBytes, v.WriteBWInfo.AudioSpeedInBytesperMS}
							msgs.Players = append(msgs.Players, msg)
						}
					}
				}
				return true
			})
			return true
		})
	} else {
		// Warning: The room should be in the "live/stream" format!
		roomInfo, exists := (this.sys.GetStreams()).Load(room)
		if exists == false {
			res.Status = 404
			res.Data = "room not found or inactive"
			return
		}

		if s, ok := roomInfo.(*core.Stream); ok {
			if s.GetReader() != nil {
				switch s.GetReader().(type) {
				case *rtmp.Reader:
					v := s.GetReader().(*rtmp.Reader)
					msg := stream{room, v.Info().URL, v.ReadBWInfo.StreamId, v.ReadBWInfo.VideoDatainBytes, v.ReadBWInfo.VideoSpeedInBytesperMS,
						v.ReadBWInfo.AudioDatainBytes, v.ReadBWInfo.AudioSpeedInBytesperMS}
					msgs.Publishers = append(msgs.Publishers, msg)
				}
			}

			s.GetWs().Range(func(k, v interface{}) bool {
				if pw, ok := v.(*core.PackWriterCloser); ok {
					if pw.GetWriter() != nil {
						switch pw.GetWriter().(type) {
						case *rtmp.Writer:
							v := pw.GetWriter().(*rtmp.Writer)
							msg := stream{room, v.Info().URL, v.WriteBWInfo.StreamId, v.WriteBWInfo.VideoDatainBytes, v.WriteBWInfo.VideoSpeedInBytesperMS,
								v.WriteBWInfo.AudioDatainBytes, v.WriteBWInfo.AudioSpeedInBytesperMS}
							msgs.Players = append(msgs.Players, msg)
						}
					}
				}
				return true
			})
		}
	}

	//resp, _ := json.Marshal(msgs)
	res.Data = msgs
}

type stream struct {
	Key             string `json:"key"`
	Url             string `json:"url"`
	StreamId        uint32 `json:"stream_id"`
	VideoTotalBytes uint64 `json:"video_total_bytes"`
	VideoSpeed      uint64 `json:"video_speed"`
	AudioTotalBytes uint64 `json:"audio_total_bytes"`
	AudioSpeed      uint64 `json:"audio_speed"`
}

type streams struct {
	Publishers []stream `json:"publishers"`
	Players    []stream `json:"players"`
}
