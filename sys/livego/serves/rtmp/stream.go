package rtmp

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/lib/modules/live/rtmprelay"
	"github.com/liwei1dao/lego/sys/livego/packet"
	"github.com/liwei1dao/lego/sys/log"
)

var (
	EmptyID = ""
)

func NewStream(gopNum int) *Stream {
	return &Stream{
		cache: NewCache(gopNum),
		ws:    &sync.Map{},
	}
}

type Stream struct {
	isStart bool
	cache   *Cache
	r       packet.ReadCloser
	ws      *sync.Map
	info    packet.Info
}

func (this *Stream) ID() string {
	if this.r != nil {
		return this.r.Info().UID
	}
	return EmptyID
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
				log.Infof("[SYS LiveGo] write timeout remove")
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

func (this *Stream) TransStop() {
	log.Debugf("TransStop: %s", this.info.Key)
	if this.isStart && this.r != nil {
		this.r.Close(fmt.Errorf("stop old"))
	}
	this.isStart = false
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
func (this *Stream) AddWriter(w packet.WriteCloser) {
	info := w.Info()
	pw := &PackWriterCloser{w: w}
	this.ws.Store(info.UID, pw)
}

func (this *Stream) AddReader(r packet.ReadCloser) {
	this.r = r
	go this.TransStart()
}

func (this *Stream) TransStart() {
	this.isStart = true
	var p packet.Packet

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
					log.Debugf("[%s] send cache packet error: %v, remove", v.w.Info(), err)
					this.ws.Delete(key)
					return true
				}
				v.init = true
			} else {
				newPacket := p
				//writeType := reflect.TypeOf(v.w)
				//log.Debugf("w.Write: type=%v, %v", writeType, v.w.Info())
				if err = v.w.Write(&newPacket); err != nil {
					log.Debugf("[%s] write packet error: %v, remove", v.w.Info(), err)
					this.ws.Delete(key)
				}
			}
			return true
		})
	}
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

	log.Debugf("StartStaticPush: current streamname=%s， appname=%s", streamname, appname)
	pushurllist, err := rtmprelay.GetStaticPushList(appname)
	if err != nil || len(pushurllist) < 1 {
		log.Debugf("StartStaticPush: GetStaticPushList error=%v", err)
		return
	}

	for _, pushurl := range pushurllist {
		pushurl := pushurl + "/" + streamname
		log.Debugf("StartStaticPush: static pushurl=%s", pushurl)

		staticpushObj := rtmprelay.GetAndCreateStaticPushObject(pushurl)
		if staticpushObj != nil {
			if err := staticpushObj.Start(); err != nil {
				log.Debugf("StartStaticPush: staticpushObj.Start %s error=%v", pushurl, err)
			} else {
				log.Debugf("StartStaticPush: staticpushObj.Start %s ok", pushurl)
			}
		} else {
			log.Debugf("StartStaticPush GetStaticPushObject %s error", pushurl)
		}
	}
}

type PackWriterCloser struct {
	init bool
	w    packet.WriteCloser
}

func NewRtmpStream(gopNum int) *RtmpStream {
	ret := &RtmpStream{
		gopNum:  gopNum,
		streams: &sync.Map{},
	}
	go ret.CheckAlive()
	return ret
}

type RtmpStream struct {
	gopNum  int
	streams *sync.Map //key
}

func (this *RtmpStream) CheckAlive() {
	for {
		<-time.After(5 * time.Second)
		this.streams.Range(func(key, val interface{}) bool {
			v := val.(*Stream)
			if v.CheckAlive() == 0 {
				this.streams.Delete(key)
			}
			return true
		})
	}
}

func (this *RtmpStream) HandleReader(r packet.ReadCloser) {
	info := r.Info()
	log.Debugf("[SYS LiveGo] HandleReader: info[%v]", info)

	var stream *Stream
	i, ok := this.streams.Load(info.Key)
	if stream, ok = i.(*Stream); ok {
		stream.TransStop()
		id := stream.ID()
		if id != EmptyID && id != info.UID {
			ns := NewStream(this.gopNum)
			stream.Copy(ns)
			stream = ns
			this.streams.Store(info.Key, ns)
		}
	} else {
		stream = NewStream(this.gopNum)
		this.streams.Store(info.Key, stream)
		stream.info = info
	}

	stream.AddReader(r)
}
