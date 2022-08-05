package core

import (
	"fmt"
	"strings"
	"sync"

	"github.com/liwei1dao/lego/sys/log"
)

func NewStream(sys ISys, log log.ILogger) *Stream {
	return &Stream{
		sys: sys,
		log: log,
		ws:  &sync.Map{},
	}
}

type PackWriterCloser struct {
	init bool
	w    WriteCloser
}

func (this *PackWriterCloser) GetWriter() WriteCloser {
	return this.w
}

type Stream struct {
	sys     ISys
	log     log.ILogger
	isStart bool
	cache   *Cache
	r       ReadCloser
	ws      *sync.Map
	info    Info
}

func (this *Stream) GetWs() *sync.Map {
	return this.ws
}
func (this *Stream) Info() Info {
	return this.info
}

func (this *Stream) SetInfo(info Info) {
	this.info = info
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

func (this *Stream) GetReader() ReadCloser {
	return this.r
}

func (this *Stream) AddWriter(w WriteCloser) {
	info := w.Info()
	pw := &PackWriterCloser{w: w}
	this.ws.Store(info.UID, pw)
}

func (this *Stream) AddReader(r ReadCloser) {
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
				this.log.Debugf("[SYS LiveGo] write timeout remove")
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
	var p Packet

	this.log.Debugf("[SYS LiveGo] TransStart: %v", this.info)

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
					this.log.Debugf("[SYS LiveGo] [%s] send cache packet error: %v, remove", v.w.Info(), err)
					this.ws.Delete(key)
					return true
				}
				v.init = true
			} else {
				newPacket := p
				//writeType := reflect.TypeOf(v.w)
				//log.Debugf("w.Write: type=%v, %v", writeType, v.w.Info())
				if err = v.w.Write(&newPacket); err != nil {
					this.log.Debugf("[SYS LiveGo] [%s] write packet error: %v, remove", v.w.Info(), err)
					this.ws.Delete(key)
				}
			}
			return true
		})
	}
}

func (this *Stream) TransStop() {
	this.log.Debugf("[SYS LiveGo] TransStop: %s", this.info.Key)

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
	pushurllist := this.sys.GetStaticPush()
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
		staticpushObj, err := this.sys.GetStaticPushObject(pushurl)
		if (staticpushObj != nil) && (err == nil) {
			return true
		} else {
			this.log.Debugf("[SYS LiveGo] SendStaticPush GetStaticPushObject %s error", pushurl)
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

	this.log.Debugf("[SYS LiveGo] StartStaticPush: current streamname=%s， appname=%s", streamname, appname)
	pushurllist := this.sys.GetStaticPush()
	if pushurllist == nil || len(pushurllist) < 1 {
		this.log.Debugf("StartStaticPush: GetStaticPushList error")
		return
	}

	for _, pushurl := range pushurllist {
		pushurl := pushurl + "/" + streamname
		this.log.Debugf("[SYS LiveGo] StartStaticPush: static pushurl=%s", pushurl)

		staticpushObj := this.sys.GetAndCreateStaticPushObject(pushurl)
		if staticpushObj != nil {
			if err := staticpushObj.Start(); err != nil {
				this.log.Errorf("[SYS LiveGo] StartStaticPush: staticpushObj.Start %s error=%v", pushurl, err)
			} else {
				this.log.Debugf("[SYS LiveGo] StartStaticPush: staticpushObj.Start %s ok", pushurl)
			}
		} else {
			this.log.Debugf("[SYS LiveGo] StartStaticPush GetStaticPushObject %s error", pushurl)
		}
	}
}

func (this *Stream) SendStaticPush(packet Packet) {
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
	pushurllist := this.sys.GetStaticPush()
	if pushurllist == nil || len(pushurllist) < 1 {
		//log.Debugf("SendStaticPush: GetStaticPushList error=%v", err)
		return
	}

	for _, pushurl := range pushurllist {
		pushurl := pushurl + "/" + streamname
		//log.Debugf("SendStaticPush: static pushurl=%s", pushurl)

		staticpushObj, err := this.sys.GetStaticPushObject(pushurl)
		if (staticpushObj != nil) && (err == nil) {
			staticpushObj.WriteAvPacket(&packet)
			//log.Debugf("SendStaticPush: WriteAvPacket %s ", pushurl)
		} else {
			this.log.Debugf("SendStaticPush GetStaticPushObject %s error", pushurl)
		}
	}
}

func (this *Stream) StopStaticPush() {
	key := this.info.Key
	this.log.Debugf("StopStaticPush......%s", key)
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

	this.log.Debugf("StopStaticPush: current streamname=%s， appname=%s", streamname, appname)
	pushurllist := this.sys.GetStaticPush()
	if pushurllist == nil || len(pushurllist) < 1 {
		this.log.Debugf("StopStaticPush: GetStaticPushList error: no pushurllist")
		return
	}

	for _, pushurl := range pushurllist {
		pushurl := pushurl + "/" + streamname
		this.log.Debugf("StopStaticPush: static pushurl=%s", pushurl)

		staticpushObj, err := this.sys.GetStaticPushObject(pushurl)
		if (staticpushObj != nil) && (err == nil) {
			staticpushObj.Stop()
			this.sys.ReleaseStaticPushObject(pushurl)
			this.log.Debugf("StopStaticPush: staticpushObj.Stop %s ", pushurl)
		} else {
			this.log.Debugf("StopStaticPush GetStaticPushObject %s error", pushurl)
		}
	}
}

func (this *Stream) closeInter() {
	if this.r != nil {
		this.StopStaticPush()
		this.log.Debugf("[%v] publisher closed", this.r.Info())
	}

	this.ws.Range(func(key, val interface{}) bool {
		v := val.(*PackWriterCloser)
		if v.w != nil {
			v.w.Close(fmt.Errorf("closed"))
			if v.w.Info().IsInterval() {
				this.ws.Delete(key)
				this.log.Debugf("[%v] player closed and remove\n", v.w.Info())
			}
		}
		return true
	})
}
