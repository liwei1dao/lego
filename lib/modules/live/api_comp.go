package live

import (
	"encoding/json"
	"fmt"
	"net"
	"net/http"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	"github.com/form3tech-oss/jwt-go"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/lib/modules/live/rtmprelay"
	"github.com/liwei1dao/lego/sys/log"
)

type Response struct {
	w      http.ResponseWriter
	Status int         `json:"status"`
	Data   interface{} `json:"data"`
}

func (r *Response) SendJson() (int, error) {
	resp, _ := json.Marshal(r)
	r.w.Header().Set("Content-Type", "application/json")
	r.w.WriteHeader(r.Status)
	return r.w.Write(resp)
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

//主机信息监控
type ApiComp struct {
	cbase.ModuleCompBase
	options IOptions
	module  ILive
	listen  net.Listener
	session map[string]*rtmprelay.RtmpRelay
}

func (this *ApiComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.options = options.(IOptions)
	this.module = module.(ILive)
	return
}

func (this *ApiComp) Start() (err error) {
	err = this.ModuleCompBase.Start()
	if this.listen, err = net.Listen("tcp", this.options.GetAPIAddr()); err == nil {
		go this.run()
	}
	return
}

func (this *ApiComp) run() (err error) {
	defer cbase.Recover()
	log.Infof("HTTP-API listen On %s", this.options.GetAPIAddr())
	this.Serve()
	return
}

func (this *ApiComp) JWTMiddleware(next http.Handler) http.Handler {
	isJWT := len(this.options.GetJWTAlgorithm()) > 0
	if !isJWT {
		return next
	}

	log.Info("Using JWT middleware")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var algorithm jwt.SigningMethod
		if len(this.options.GetJWTAlgorithm()) > 0 {
			algorithm = jwt.GetSigningMethod(this.options.GetJWTAlgorithm())
		}

		if algorithm == nil {
			algorithm = jwt.SigningMethodHS256
		}

		jwtMiddleware := jwtmiddleware.New(jwtmiddleware.Options{
			Extractor: jwtmiddleware.FromFirst(jwtmiddleware.FromAuthHeader, jwtmiddleware.FromParameter("jwt")),
			ValidationKeyGetter: func(token *jwt.Token) (interface{}, error) {
				return []byte(this.options.GetJWTAlgorithm()), nil
			},
			SigningMethod: algorithm,
			ErrorHandler: func(w http.ResponseWriter, r *http.Request, err string) {
				res := &Response{
					w:      w,
					Status: 403,
					Data:   err,
				}
				res.SendJson()
			},
		})

		jwtMiddleware.HandlerWithNext(w, r, next.ServeHTTP)
	})
}

func (this *ApiComp) Serve() error {
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
	http.Serve(this.listen, this.JWTMiddleware(mux))
	return nil
}

//http://127.0.0.1:8090/control/push?&oper=start&app=live&name=123456&url=rtmp://192.168.16.136/live/123456
func (this *ApiComp) handlePush(w http.ResponseWriter, req *http.Request) {
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

	log.Debugf("control push: oper=%v, app=%v, name=%v, url=%v", oper, app, name, url)
	if (len(app) <= 0) || (len(name) <= 0) || (len(url) <= 0) {
		res.Data = "control push parameter error, please check them."
		return
	}

	localurl := "rtmp://127.0.0.1" + this.options.GetRtmpAddr() + "/" + app + "/" + name
	remoteurl := url

	keyString := "push:" + app + "/" + name
	if oper == "stop" {
		pushRtmprelay, found := this.session[keyString]
		if !found {
			retString = fmt.Sprintf("<h1>session key[%s] not exist, please check it again.</h1>", keyString)
			res.Data = retString
			return
		}
		log.Debugf("rtmprelay stop push %s from %s", remoteurl, localurl)
		pushRtmprelay.Stop()

		delete(this.session, keyString)
		retString = fmt.Sprintf("<h1>push url stop %s ok</h1></br>", url)
		res.Data = retString
		log.Debugf("push stop return %s", retString)
	} else {
		pushRtmprelay := rtmprelay.NewRtmpRelay(&localurl, &remoteurl)
		log.Debugf("rtmprelay start push %s from %s", remoteurl, localurl)
		err = pushRtmprelay.Start()
		if err != nil {
			retString = fmt.Sprintf("push error=%v", err)
		} else {
			retString = fmt.Sprintf("<h1>push url start %s ok</h1></br>", url)
			this.session[keyString] = pushRtmprelay
		}

		res.Data = retString
		log.Debugf("push start return %s", retString)
	}
}

//http://127.0.0.1:8090/control/pull?&oper=start&app=live&name=123456&url=rtmp://192.168.16.136/live/123456
func (this *ApiComp) handlePull(w http.ResponseWriter, req *http.Request) {
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

	log.Debugf("control pull: oper=%v, app=%v, name=%v, url=%v", oper, app, name, url)
	if (len(app) <= 0) || (len(name) <= 0) || (len(url) <= 0) {
		res.Status = 400
		res.Data = "control push parameter error, please check them."
		return
	}

	remoteurl := "rtmp://127.0.0.1" + this.options.GetRtmpAddr() + "/" + app + "/" + name
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
		log.Debugf("rtmprelay stop push %s from %s", remoteurl, localurl)
		pullRtmprelay.Stop()

		delete(this.session, keyString)
		retString = fmt.Sprintf("<h1>push url stop %s ok</h1></br>", url)
		res.Status = 400
		res.Data = retString
		log.Debugf("pull stop return %s", retString)
	} else {
		pullRtmprelay := rtmprelay.NewRtmpRelay(&localurl, &remoteurl)
		log.Debugf("rtmprelay start push %s from %s", remoteurl, localurl)
		err = pullRtmprelay.Start()
		if err != nil {
			retString = fmt.Sprintf("push error=%v", err)
		} else {
			this.session[keyString] = pullRtmprelay
			retString = fmt.Sprintf("<h1>push url start %s ok</h1></br>", url)
		}
		res.Status = 400
		res.Data = retString
		log.Debugf("pull start return %s", retString)
	}
}

//http://127.0.0.1:8090/control/get?room=ROOM_NAME
func (this *ApiComp) handleGet(w http.ResponseWriter, r *http.Request) {
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

	msg, err := this.module.GetCacheComp().GetChannelKey(room)
	if err != nil {
		msg = err.Error()
		res.Status = 400
	}
	res.Data = msg
}

//http://127.0.0.1:8090/control/reset?room=ROOM_NAME
func (this *ApiComp) handleReset(w http.ResponseWriter, r *http.Request) {
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

	msg, err := this.module.GetCacheComp().GetChannelKey(room)

	if err != nil {
		msg = err.Error()
		res.Status = 400
	}

	res.Data = msg
}

//http://127.0.0.1:8090/control/delete?room=ROOM_NAME
func (this *ApiComp) handleDelete(w http.ResponseWriter, r *http.Request) {
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

	if this.module.GetCacheComp().DeleteChannel(room) {
		res.Data = "Ok"
		return
	}
	res.Status = 404
	res.Data = "room not found"
}

//http://127.0.0.1:8090/stat/livestat
func (this *ApiComp) GetLiveStatics(w http.ResponseWriter, req *http.Request) {
	res := &Response{
		w:      w,
		Data:   nil,
		Status: 200,
	}

	defer res.SendJson()

	rtmpStream := this.module.GetHandler().(*RtmpStream)
	if rtmpStream == nil {
		res.Status = 500
		res.Data = "Get rtmp stream information error"
		return
	}

	msgs := new(streams)

	rtmpStream.GetStreams().Range(func(key, val interface{}) bool {
		if s, ok := val.(*Stream); ok {
			if s.GetReader() != nil {
				switch s.GetReader().(type) {
				case *VirReader:
					v := s.GetReader().(*VirReader)
					msg := stream{key.(string), v.Info().URL, v.ReadBWInfo.StreamId, v.ReadBWInfo.VideoDatainBytes, v.ReadBWInfo.VideoSpeedInBytesperMS,
						v.ReadBWInfo.AudioDatainBytes, v.ReadBWInfo.AudioSpeedInBytesperMS}
					msgs.Publishers = append(msgs.Publishers, msg)
				}
			}
		}
		return true
	})

	rtmpStream.GetStreams().Range(func(key, val interface{}) bool {
		ws := val.(*Stream).GetWs()
		ws.Range(func(k, v interface{}) bool {
			if pw, ok := v.(*PackWriterCloser); ok {
				if pw.GetWriter() != nil {
					switch pw.GetWriter().(type) {
					case *VirWriter:
						v := pw.GetWriter().(*VirWriter)
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

	//resp, _ := json.Marshal(msgs)
	res.Data = msgs
}
