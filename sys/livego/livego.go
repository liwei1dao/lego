package livego

import (
	"fmt"
	"net"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/livego/container/flv"
	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/utils/container/id"
)

func newSys(options Options) (sys *LiveGo, err error) {
	sys = &LiveGo{
		options:       options,
		keyLock:       new(sync.RWMutex),
		keys:          make(map[string]string),
		channelLock:   new(sync.RWMutex),
		channels:      make(map[string]string),
		streams:       new(sync.Map),
		mapLock:       new(sync.RWMutex),
		staticPushMap: make(map[string]*core.StaticPush),
	}
	err = sys.init()
	return
}

type LiveGo struct {
	options Options
	///流频道管理
	keyLock     *sync.RWMutex
	keys        map[string]string
	channelLock *sync.RWMutex
	channels    map[string]string
	///流管理
	streams *sync.Map
	///静态推送管理
	mapLock       *sync.RWMutex
	staticPushMap map[string](*core.StaticPush)
}

func (this *LiveGo) init() (err error) {
	var (
		rtmpListen net.Listener
	)
	if rtmpListen, err = net.Listen("tcp", this.options.RTMPAddr); err != nil {
		this.Errorf("init err:%v", err)
		return
	}
	this.Infof("RTMP Listen On:%s", this.options.RTMPAddr)
	go this.Serve(rtmpListen)
	return
}

func (this *LiveGo) Serve(listener net.Listener) (err error) {
	for {
		var netconn net.Conn
		if netconn, err = listener.Accept(); err != nil {
			this.Errorf("Serve err: ", err)
			return
		}
		conn := core.NewConn(netconn, this)
		this.Debugf("new client, connect remote: ", conn.RemoteAddr().String(),
			"local:", conn.LocalAddr().String())
		go this.handleConn(conn)
	}
}

func (this *LiveGo) handleConn(conn *core.Conn) (err error) {
	if err := conn.HandshakeServer(); err != nil {
		conn.Close()
		this.Errorf("handleConn HandshakeServer err: ", err)
		return err
	}
	connServer := core.NewConnServer(conn)
	if err := connServer.ReadMsg(); err != nil {
		conn.Close()
		this.Errorf("handleConn read msg err:%v", err)
		return err
	}
	appname, name, _ := connServer.GetInfo()
	if appname != this.options.Appname {
		err := fmt.Errorf("application name=%s is not configured", appname)
		conn.Close()
		this.Errorf("CheckAppName err: ", err)
	}
	this.Debugf("handleConn: IsPublisher=%v", connServer.IsPublisher())
	if connServer.IsPublisher() {
		if this.options.RTMPNoAuth {
			key, err := this.GetKey(name)
			if err != nil {
				err := fmt.Errorf("Cannot create err=%s", err.Error())
				conn.Close()
				this.Errorf("GetKey err:%v", err)
				return err
			}
			name = key
		}
		channel, err := this.GetChannel(name)
		if err != nil {
			err := fmt.Errorf("invalid key err=%s", err.Error())
			conn.Close()
			this.Errorf("CheckKey err:%v", err)
			return err
		}
		connServer.PublishInfo.Name = channel
		if this.options.StaticPush != nil && len(this.options.StaticPush) > 0 {
			this.Debugf("GetStaticPushUrlList: %v", this.options.StaticPush)
		}
		reader := NewReader(connServer)
		this.HandleReader(reader)
		this.Debugf("new publisher: %+v", reader.Info())
		if this.GetFLVArchive() {
			flvWriter := flv.NewFlvDvr(this)
			this.HandleWriter(flvWriter.GetWriter(reader.Info()))
		}
	} else {
		writer := NewWriter(connServer)
		this.Debugf("new player: %+v", writer.Info())
		this.HandleWriter(writer)
	}
	return
}

///流管理***********************************************************************
func (this *LiveGo) GetStreams() *sync.Map {
	return this.streams
}

//监测存活
func (this *LiveGo) CheckAlive() {
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

//读处理
func (this *LiveGo) HandleReader(r core.ReadCloser) {
	info := r.Info()
	this.Debugf("HandleReader: info[%v]", info)
	var stm *Stream
	i, ok := this.streams.Load(info.Key)
	if stm, ok = i.(*Stream); ok {
		stm.TransStop()
		id := stm.ID()
		if id != EmptyID && id != info.UID {
			ns := NewStream()
			stm.Copy(ns)
			stm = ns
			this.streams.Store(info.Key, ns)
		}
	} else {
		stm = NewStream()
		this.streams.Store(info.Key, stm)
		stm.info = info
	}
	stm.AddReader(r)
}

//写处理
func (this *LiveGo) HandleWriter(w core.WriteCloser) {
	info := w.Info()
	this.Debugf("HandleWriter: info[%v]", info)

	var s *Stream
	item, ok := this.streams.Load(info.Key)
	if !ok {
		this.Debugf("HandleWriter: not found create new info[%v]", info)
		s = NewStream()
		this.streams.Store(info.Key, s)
		s.info = info
	} else {
		s = item.(*Stream)
		s.AddWriter(w)
	}
}

///参数列表***********************************************************************
func (this *LiveGo) GetAppname() string {
	return this.options.Appname
}
func (this *LiveGo) GetHls() bool {
	return this.options.Hls
}
func (this *LiveGo) GetUseHlsHttps() bool {
	return this.options.UseHlsHttps
}
func (this *LiveGo) GetHlsServerCrt() string {
	return this.options.HlsServerCrt
}
func (this *LiveGo) GetHlsServerKey() string {
	return this.options.HlsServerKey
}
func (this *LiveGo) GetFlv() bool {
	return this.options.Flv
}
func (this *LiveGo) GetApi() bool {
	return this.options.Api
}
func (this *LiveGo) GetApiAddr() string {
	return this.options.APIAddr
}
func (this *LiveGo) GetJWTSecret() string {
	return this.options.JWTSecret
}
func (this *LiveGo) GetJWTAlgorithm() string {
	return this.options.JWTAlgorithm
}
func (this *LiveGo) GetStaticPush() []string {
	return this.options.StaticPush
}
func (this *LiveGo) GetRTMPAddr() string {
	return this.options.RTMPAddr
}

func (this *LiveGo) GetRTMPNoAuth() bool {
	return this.options.RTMPNoAuth
}
func (this *LiveGo) GetFLVDir() string {
	return this.options.FLVDir
}
func (this *LiveGo) GetFLVArchive() bool {
	return this.options.FLVArchive
}
func (this *LiveGo) GetEnableTLSVerify() bool {
	return this.options.EnableTLSVerify
}
func (this *LiveGo) GetTimeout() int {
	return this.options.Timeout
}
func (this *LiveGo) GetConnBuffSzie() int {
	return this.options.ConnBuffSzie
}
func (this *LiveGo) GetDebug() bool {
	return this.options.Debug
}

///日志***********************************************************************
func (this *LiveGo) Debugf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Debugf("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Infof(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Infof("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Warnf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Warnf("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Errorf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Errorf("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Panicf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Panicf("[SYS LiveGo] "+format, a)
	}
}
func (this *LiveGo) Fatalf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Fatalf("[SYS LiveGo] "+format, a)
	}
}

///房间***********************************************************************
func (this *LiveGo) SetKey(channel string) (key string, err error) {
	key = id.NewXId()
	this.keyLock.Lock()
	this.keys[key] = channel
	this.keyLock.Unlock()
	this.channelLock.Lock()
	this.channels[channel] = key
	this.channelLock.Unlock()
	return
}
func (this *LiveGo) GetKey(channel string) (key string, err error) {
	var ok bool
	this.channelLock.RLock()
	key, ok = this.channels[channel]
	this.channelLock.RUnlock()
	if !ok {
		err = fmt.Errorf("no channel:%s", channel)
	}
	return
}
func (this *LiveGo) GetChannel(key string) (channel string, err error) {
	var ok bool
	this.channelLock.RLock()
	channel, ok = this.keys[key]
	this.channelLock.RUnlock()
	if !ok {
		err = fmt.Errorf("no key:%s", channel)
	}
	return
}
func (this *LiveGo) DeleteChannel(channel string) (ok bool) {
	this.channelLock.Lock()
	delete(this.channels, channel)
	this.channelLock.Unlock()
	ok = true
	return
}
func (this *LiveGo) DeleteKey(key string) (ok bool) {
	this.keyLock.Lock()
	delete(this.keys, key)
	this.keyLock.Unlock()
	ok = true
	return
}

///静态推送管理***********************************************************************
func (this *LiveGo) GetAndCreateStaticPushObject(rtmpurl string) (pushobj *core.StaticPush) {
	this.mapLock.RLock()
	staticpush, ok := this.staticPushMap[rtmpurl]
	this.Debugf("GetAndCreateStaticPushObject: %s, return %v", rtmpurl, ok)
	if !ok {
		this.mapLock.RUnlock()
		newStaticpush := core.NewStaticPush(rtmpurl)
		this.mapLock.Lock()
		this.staticPushMap[rtmpurl] = newStaticpush
		this.mapLock.Unlock()
		return newStaticpush
	}
	this.mapLock.RUnlock()
	return staticpush
}

func (this *LiveGo) GetStaticPushObject(rtmpurl string) (pushobj *core.StaticPush, err error) {
	this.mapLock.RLock()
	if staticpush, ok := this.staticPushMap[rtmpurl]; ok {
		this.mapLock.RUnlock()
		return staticpush, nil
	}
	this.mapLock.RUnlock()

	return nil, fmt.Errorf("G_StaticPushMap[%s] not exist....", rtmpurl)
}

func (this *LiveGo) ReleaseStaticPushObject(rtmpurl string) {
	this.mapLock.RLock()
	if _, ok := this.staticPushMap[rtmpurl]; ok {
		this.mapLock.RUnlock()
		this.Debugf("ReleaseStaticPushObject %s ok", rtmpurl)
		this.mapLock.Lock()
		delete(this.staticPushMap, rtmpurl)
		this.mapLock.Unlock()
	} else {
		this.mapLock.RUnlock()
		this.Debugf("ReleaseStaticPushObject: not find %s", rtmpurl)
	}
	return
}
