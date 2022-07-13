package cachego

import (
	"sync"
	"time"
)

func newSys(options *Options) (sys *CacheGo, err error) {
	sys = &CacheGo{
		options: options,
		items:   make(map[string]Item),
	}
	err = sys.init()
	return
}

type CacheGo struct {
	options   *Options
	items     map[string]Item
	mu        sync.RWMutex
	onEvicted func(string, interface{})
	stop      chan bool
}

func (this *CacheGo) init() (err error) {
	if this.options.CleanupInterval > 0 {
		go this.run()
	}
	return
}

func (this *CacheGo) run() {
	ticker := time.NewTicker(time.Duration(this.options.CleanupInterval) * time.Second)
	for {
		select {
		case <-ticker.C:
			this.DeleteExpired()
		case <-this.stop:
			ticker.Stop()
			return
		}
	}
}

func (this *CacheGo) Clsoe() {
	this.stop <- true
}

///日志***********************************************************************
func (this *CacheGo) Debugf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Debugf("[SYS CacheGo] "+format, a...)
	}
}
func (this *CacheGo) Infof(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Infof("[SYS CacheGo] "+format, a...)
	}
}
func (this *CacheGo) Warnf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Warnf("[SYS CacheGo] "+format, a...)
	}
}
func (this *CacheGo) Errorf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Errorf("[SYS CacheGo] "+format, a...)
	}
}
func (this *CacheGo) Panicf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Panicf("[SYS CacheGo] "+format, a...)
	}
}
func (this *CacheGo) Fatalf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Fatalf("[SYS CacheGo] "+format, a...)
	}
}
