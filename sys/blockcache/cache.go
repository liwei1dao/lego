package blockcache

import (
	"sync/atomic"
	"unsafe"
)

func newSys(options Options) (sys *Cache, err error) {
	sys = &Cache{options: options}
	return
}

type Item struct {
	Size  uint64
	Value interface{}
}

type Cache struct {
	options  Options
	inpip    chan interface{}
	outpip   chan interface{}
	head     uint64
	tail     uint64
	element  []*Item
	capacity uint64
	size     uint64
	free     uint64
}

func (this *Cache) In() chan<- interface{} {
	return this.inpip
}

func (this *Cache) Out() <-chan interface{} {
	return this.outpip
}
func (this *Cache) Run() {
	for v := range this.inpip {
		siez := uint64(unsafe.Sizeof(v))
		if siez > this.options.CacheMaxSzie { //异常数据
			this.Errorf("item size:%d large CacheMaxSzie:%d", siez, this.options.CacheMaxSzie)
		} else if siez > this.free {

		} else {

		}
	}

}

func (this *Cache) Push(v interface{}) bool {
	oldTail := atomic.LoadUint64(&this.tail)
	oldHead := atomic.LoadUint64(&this.head)
	if this.isFull(oldTail, oldHead) {
		return false
	}

	newTail := (oldTail + 1) & this.mask
	// ----------- BEGIN --------------
	tailNode := atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&this.element[newTail])))
	if tailNode != nil {
		return false
	}
	// ----------- END --------------
	if !atomic.CompareAndSwapUint64(&this.tail, oldTail, newTail) {
		return false
	}

	atomic.StorePointer((*unsafe.Pointer)(unsafe.Pointer(&this.element[newTail])), unsafe.Pointer(&v))
	return true
}
func (this *Cache) isEmpty(tail uint64, head uint64) bool {
	return tail-head == 0
}
func (this *Cache) isFull(tail uint64, head uint64) bool {
	return tail-head == this.capacity-1
}

///日志***********************************************************************
func (this *Cache) Debugf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Debugf("[SYS BlockCache] "+format, a)
	}
}
func (this *Cache) Infof(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Infof("[SYS BlockCache] "+format, a)
	}
}
func (this *Cache) Warnf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Warnf("[SYS BlockCache] "+format, a)
	}
}
func (this *Cache) Errorf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Errorf("[SYS BlockCache] "+format, a)
	}
}
func (this *Cache) Panicf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Panicf("[SYS BlockCache] "+format, a)
	}
}
func (this *Cache) Fatalf(format string, a ...interface{}) {
	if this.options.Debug {
		this.options.Log.Fatalf("[SYS BlockCache] "+format, a)
	}
}
