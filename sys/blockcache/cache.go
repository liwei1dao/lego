package blockcache

import (
	"sync"
	"sync/atomic"
	"unsafe"

	"github.com/liwei1dao/lego/utils/container"
)

func newSys(options *Options) (sys *Cache, err error) {
	sys = &Cache{
		options:  options,
		inpip:    make(chan interface{}),
		outpip:   make(chan interface{}),
		outnotic: make(chan struct{}),
		element:  container.NewLKQueue(),
		free:     options.CacheMaxSzie,
	}
	go sys.run()
	return
}

type Item struct {
	Size  int64
	Value interface{}
}

type Cache struct {
	options   *Options
	inpip     chan interface{}
	outpip    chan interface{}
	outnotic  chan struct{}
	instate   int32
	outruning int32
	element   *container.LKQueue
	free      int64
	close     int32
	wg        sync.WaitGroup
}

func (this *Cache) In() chan<- interface{} {
	return this.inpip
}

func (this *Cache) Out() <-chan interface{} {
	return this.outpip
}

func (this *Cache) Close() {
	atomic.StoreInt32(&this.close, 1)
	close(this.inpip)
	close(this.outnotic)
	this.wg.Wait()
	close(this.outpip)
}

func (this *Cache) run() {
	for v := range this.inpip {
		siez := int64(unsafe.Sizeof(v))
		if siez > this.options.CacheMaxSzie { //异常数据
			this.options.Log.Errorf("item size:%d large CacheMaxSzie:%d", siez, this.options.CacheMaxSzie)
			continue
		} else if siez > atomic.LoadInt64(&this.free) { //空间不足
			atomic.StoreInt32(&this.instate, 1)
		locp:
			for _ = range this.outnotic {
				if siez > atomic.LoadInt64(&this.free) {
					atomic.StoreInt32(&this.instate, 1)
				} else {
					this.element.Enqueue(&Item{Size: siez, Value: v})
					atomic.AddInt64(&this.free, -1*siez)
					break locp
				}
			}
		} else {
			this.element.Enqueue(&Item{Size: siez, Value: v})
			atomic.AddInt64(&this.free, -1*siez)
		}
		if atomic.CompareAndSwapInt32(&this.outruning, 0, 1) {
			this.wg.Add(1)
			go func() {
			locp:
				for {
					v := this.element.Dequeue()
					if v != nil && atomic.LoadInt32(&this.close) == 0 {
						item := v.(*Item)
						atomic.AddInt64(&this.free, item.Size)
						if atomic.CompareAndSwapInt32(&this.instate, 1, 0) {
							this.outnotic <- struct{}{}
						}
						this.outpip <- item.Value
					} else {
						break locp
					}
				}
				atomic.StoreInt32(&this.outruning, 0)
				this.wg.Done()
			}()
		}
	}
}

///日志***********************************************************************
