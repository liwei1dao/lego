package livego

import (
	"fmt"
	"strings"
	"sync"

	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/log"
)

func NewStream() *Stream {
	return &Stream{
		ws: &sync.Map{},
	}
}

type PackWriterCloser struct {
	init bool
	w    core.WriteCloser
}
type Stream struct {
	server  core.IServer
	isStart bool
	cache   *Cache
	r       core.ReadCloser
	ws      *sync.Map
	info    core.Info
}

func (this *Stream) ID() string {
	if this.r != nil {
		return this.r.Info().UID
	}
	return EmptyID
}

func (this *Stream) Copy(dst *Stream) {
	dst.info = this.info
	this.ws.Range(func(key, val interface{}) bool {
		v := val.(*PackWriterCloser)
		this.ws.Delete(key)
		v.w.CalcBaseTimestamp()
		dst.AddWriter(v.w)
		return true
	})
}

func (this *Stream) AddWriter(w core.WriteCloser) {
	info := w.Info()
	pw := &PackWriterCloser{w: w}
	this.ws.Store(info.UID, pw)
}

func (this *Stream) AddReader(r core.ReadCloser) {
	this.r = r
	go this.TransStart()
}

func (this *Stream) CheckAlive() (n int) {
	if this.r != nil && this.isStart {
		if this.r.Alive() {
			n++
		} else {
			this.r.Close(fmt.Errorf("read timeout"))
		}
	}

	this.ws.Range(func(key, val interface{}) bool {
		v := val.(*PackWriterCloser)
		if v.w != nil {
			//Alive from RWBaser, check last frame now - timestamp, if > timeout then Remove it
			if !v.w.Alive() {
				log.Debugf("[SYS LiveGo] write timeout remove")
				this.ws.Delete(key)
				v.w.Close(fmt.Errorf("write timeout"))
				return true
			}
			n++
		}
		return true
	})

	return
}

func (this *Stream) TransStart() {
	this.isStart = true
	var p core.Packet

	log.Debugf("[SYS LiveGo] TransStart: %v", this.info)

	this.StartStaticPush()

	for {
		if !this.isStart {
			this.closeInter()
			return
		}
		err := this.r.Read(&p)
		if err != nil {
			this.closeInter()
			this.isStart = false
			return
		}

		if this.IsSendStaticPush() {
			this.SendStaticPush(p)
		}

		this.cache.Write(p)

		this.ws.Range(func(key, val interface{}) bool {
			v := val.(*PackWriterCloser)
			if !v.init {
				//log.Debugf("cache.send: %v", v.w.Info())
				if err = this.cache.Send(v.w); err != nil {
					log.Debugf("[SYS LiveGo] [%s] send cache packet error: %v, remove", v.w.Info(), err)
					this.ws.Delete(key)
					return true
				}
				v.init = true
			} else {
				newPacket := p
				//writeType := reflect.TypeOf(v.w)
				//log.Debugf("w.Write: type=%v, %v", writeType, v.w.Info())
				if err = v.w.Write(&newPacket); err != nil {
					log.Debugf("[SYS LiveGo] [%s] write packet error: %v, remove", v.w.Info(), err)
					this.ws.Delete(key)
				}
			}
			return true
		})
	}
}
func (this *Stream) TransStop() {
	log.Debugf("[SYS LiveGo] TransStop: %s", this.info.Key)

	if this.isStart && this.r != nil {
		this.r.Close(fmt.Errorf("stop old"))
	}

	this.isStart = false
}

func (this *Stream) IsSendStaticPush() bool {
	key := this.info.Key
	dscr := strings.Split(key, "/")
	if len(dscr) < 1 {
		return false
	}
	pushurllist := this.server.GetStaticPush()
	if pushurllist == nil || len(pushurllist) < 1 {
		return false
	}

	index := strings.Index(key, "/")
	if index < 0 {
		return false
	}

	streamname := key[index+1:]

	for _, pushurl := range pushurllist {
		pushurl := pushurl + "/" + streamname
		staticpushObj, err := this.server.GetStaticPushObject(pushurl)
		if (staticpushObj != nil) && (err == nil) {
			return true
		} else {
			log.Debugf("[SYS LiveGo] SendStaticPush GetStaticPushObject %s error", pushurl)
		}
	}
	return false
}

/*检测本application下是否配置static_push,
如果配置, 启动push远端的连接*/
func (this *Stream) StartStaticPush() {
	key := this.info.Key

	dscr := strings.Split(key, "/")
	if len(dscr) < 1 {
		return
	}

	index := strings.Index(key, "/")
	if index < 0 {
		return
	}

	streamname := key[index+1:]
	appname := dscr[0]

	log.Debugf("[SYS LiveGo] StartStaticPush: current streamname=%s， appname=%s", streamname, appname)
	pushurllist := this.server.GetStaticPush()
	if pushurllist == nil || len(pushurllist) < 1 {
		log.Debugf("StartStaticPush: GetStaticPushList error")
		return
	}

	for _, pushurl := range pushurllist {
		pushurl := pushurl + "/" + streamname
		log.Debugf("[SYS LiveGo] StartStaticPush: static pushurl=%s", pushurl)

		staticpushObj := this.server.GetAndCreateStaticPushObject(pushurl)
		if staticpushObj != nil {
			if err := staticpushObj.Start(); err != nil {
				log.Errorf("[SYS LiveGo] StartStaticPush: staticpushObj.Start %s error=%v", pushurl, err)
			} else {
				log.Debugf("[SYS LiveGo] StartStaticPush: staticpushObj.Start %s ok", pushurl)
			}
		} else {
			log.Debugf("[SYS LiveGo] StartStaticPush GetStaticPushObject %s error", pushurl)
		}
	}
}

func (this *Stream) SendStaticPush(packet core.Packet) {
	key := this.info.Key

	dscr := strings.Split(key, "/")
	if len(dscr) < 1 {
		return
	}

	index := strings.Index(key, "/")
	if index < 0 {
		return
	}

	streamname := key[index+1:]
	pushurllist := this.server.GetStaticPush()
	if pushurllist == nil || len(pushurllist) < 1 {
		//log.Debugf("SendStaticPush: GetStaticPushList error=%v", err)
		return
	}

	for _, pushurl := range pushurllist {
		pushurl := pushurl + "/" + streamname
		//log.Debugf("SendStaticPush: static pushurl=%s", pushurl)

		staticpushObj, err := this.server.GetStaticPushObject(pushurl)
		if (staticpushObj != nil) && (err == nil) {
			staticpushObj.WriteAvPacket(&packet)
			//log.Debugf("SendStaticPush: WriteAvPacket %s ", pushurl)
		} else {
			log.Debugf("SendStaticPush GetStaticPushObject %s error", pushurl)
		}
	}
}

func (this *Stream) StopStaticPush() {
	key := this.info.Key
	if this.server.GetDebug() {
		log.Debugf("StopStaticPush......%s", key)
	}

	dscr := strings.Split(key, "/")
	if len(dscr) < 1 {
		return
	}

	index := strings.Index(key, "/")
	if index < 0 {
		return
	}

	streamname := key[index+1:]
	appname := dscr[0]

	log.Debugf("StopStaticPush: current streamname=%s， appname=%s", streamname, appname)
	pushurllist := this.server.GetStaticPush()
	if pushurllist == nil || len(pushurllist) < 1 {
		log.Debugf("StopStaticPush: GetStaticPushList error: no pushurllist")
		return
	}

	for _, pushurl := range pushurllist {
		pushurl := pushurl + "/" + streamname
		log.Debugf("StopStaticPush: static pushurl=%s", pushurl)

		staticpushObj, err := this.server.GetStaticPushObject(pushurl)
		if (staticpushObj != nil) && (err == nil) {
			staticpushObj.Stop()
			this.server.ReleaseStaticPushObject(pushurl)
			log.Debugf("StopStaticPush: staticpushObj.Stop %s ", pushurl)
		} else {
			log.Debugf("StopStaticPush GetStaticPushObject %s error", pushurl)
		}
	}
}

func (this *Stream) closeInter() {
	if this.r != nil {
		this.StopStaticPush()
		log.Debugf("[%v] publisher closed", this.r.Info())
	}

	this.ws.Range(func(key, val interface{}) bool {
		v := val.(*PackWriterCloser)
		if v.w != nil {
			v.w.Close(fmt.Errorf("closed"))
			if v.w.Info().IsInterval() {
				this.ws.Delete(key)
				log.Debugf("[%v] player closed and remove\n", v.w.Info())
			}
		}
		return true
	})
}
