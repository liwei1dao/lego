package discovery

import (
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/discovery/consul"
	"github.com/liwei1dao/lego/sys/discovery/dcore"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/codec/json"
)

func newSys(options *Options) (sys *Discovery, err error) {
	sys = &Discovery{options: options}
	switch options.StoreType {
	case StoreConsul:
		sys.store, err = consul.NewConsulStore(options.Endpoints, options.Config)
		break
	default:
		err = fmt.Errorf("StoreType:%d unsupported type", options.StoreType)
	}
	return
}

type Discovery struct {
	options                 *Options
	store                   dcore.IStore
	RetriesAfterWatchFailed int //-1 表示一直重试 0 表示不重试
	mu                      sync.Mutex
	pairsMu                 sync.RWMutex
	pairs                   []*core.ServiceNode
	chans                   []chan []*core.ServiceNode
	stopregisterSignal      chan struct{}
	stopwatchSignal         chan struct{}
}

func (this *Discovery) GetServices() []*core.ServiceNode {
	this.pairsMu.RLock()
	defer this.pairsMu.RUnlock()
	return this.pairs
}

func (this *Discovery) WatchService() chan []*core.ServiceNode {
	this.mu.Lock()
	defer this.mu.Unlock()

	ch := make(chan []*core.ServiceNode, 10)
	this.chans = append(this.chans, ch)
	return ch
}
func (this *Discovery) GetNodePath() string {
	return fmt.Sprintf("%s/%s/%s", this.options.BasePath, this.options.ServiceNode.Type, this.options.ServiceNode.Id)
}

func (this *Discovery) Start() (err error) {
	var (
		ps []*dcore.KVPair
	)
	ps, err = this.store.List(this.options.BasePath)
	if err != nil && err != dcore.ErrKeyNotFound {
		this.options.Log.Errorf("cannot get services of from registry: %v, err: %v", this.options.BasePath, err)
		return err
	}

	pairs := make([]*core.ServiceNode, 0, len(ps))
	prefix := this.options.BasePath + "/"
	for _, p := range ps {
		if !strings.HasPrefix(p.Key, prefix) { // avoid prefix issue of consul List
			continue
		}
		// k := strings.TrimPrefix(p.Key, prefix)
		pair := &core.ServiceNode{}
		if err = this.Unmarshal(p.Value, pair); err != nil {
			this.options.Log.Errorln(err)
		}
		pairs = append(pairs, pair)
	}
	this.pairsMu.Lock()
	this.pairs = pairs
	this.pairsMu.Unlock()
	if this.options.UpdateInterval > 0 && this.options.ServiceNode != nil {
		//先注册进去一次
		d, _ := this.Marshal(this.options.ServiceNode)
		nodePath := this.GetNodePath()
		if err = this.store.Put(nodePath, d, &dcore.WriteOptions{TTL: this.options.UpdateInterval * 2}); err != nil {
			this.options.Log.Errorf("consul put path %s err: %v", nodePath, err)
		}
		go func() {
			ticker := time.NewTicker(this.options.UpdateInterval)
			defer ticker.Stop()
			defer this.store.Close()
		locp:
			for {
				select {
				case <-this.stopregisterSignal:
					break locp
				case <-ticker.C:
					d, _ := this.Marshal(this.options.ServiceNode)
					nodePath := this.GetNodePath()
					err := this.store.Put(nodePath, d, &dcore.WriteOptions{TTL: this.options.UpdateInterval * 2})
					if err != nil {
						this.options.Log.Error("consul put ", log.Field{Key: "nodePath", Value: nodePath}, log.Field{Key: "value", Value: string(d)}, log.Field{Key: "err", Value: err.Error()})
					}
				}
			}
			this.options.Log.Debug("close Timed registration coroutine")
		}()
	}
	go this.watch()
	return nil
}

func (this *Discovery) Close() error {
	_ = this.store.Delete(this.options.BasePath)
	this.stopregisterSignal <- struct{}{}
	this.stopwatchSignal <- struct{}{}
	this.options.Log.Debug("Stop End !")
	return nil
}

//监控
func (this *Discovery) watch() {
	defer func() {
		this.store.Close()
	}()

	var err error
	var c <-chan []*dcore.KVPair
	var tempDelay time.Duration

	retry := this.RetriesAfterWatchFailed
	for this.RetriesAfterWatchFailed < 0 || retry >= 0 {
		c, err = this.store.WatchTree(this.options.BasePath, this.stopwatchSignal)
		if err != nil {
			if this.RetriesAfterWatchFailed > 0 {
				retry--
			}
			if tempDelay == 0 {
				tempDelay = 1 * time.Second
			} else {
				tempDelay *= 2
			}
			if max := 30 * time.Second; tempDelay > max {
				tempDelay = max
			}
			this.options.Log.Warnf("can not watchtree (with retry %d, sleep %v): %s: %v", retry, tempDelay, this.options.BasePath, err)
			time.Sleep(tempDelay)
			continue
		}
		break
	}

	if err != nil {
		this.options.Log.Errorf("can't watch %s: %v", this.options.BasePath, err)
		return
	}

	prefix := this.options.BasePath + "/"
	for ps := range c {
		var pairs []*core.ServiceNode // latest servers
		if ps == nil {
			this.pairsMu.Lock()
			this.pairs = pairs
			this.pairsMu.Unlock()
			continue
		}
		for _, p := range ps {
			if !strings.HasPrefix(p.Key, prefix) { // avoid prefix issue of consul List
				continue
			}
			// k := strings.TrimPrefix(p.Key, prefix)
			pair := &core.ServiceNode{}
			// this.options.Log.Debug("watch!", log.Field{Key: "key", Value: p.Key}, log.Field{Key: "value", Value: string(p.Value)})
			if err = this.Unmarshal(p.Value, pair); err != nil {
				// this.options.Log.Error("watch Unmarshal err !", log.Field{Key: "key", Value: p.Key}, log.Field{Key: "value", Value: string(p.Value)}, log.Field{Key: "err", Value: err})
				continue
			}
			pairs = append(pairs, pair)
		}
		this.pairsMu.Lock()
		this.pairs = pairs
		this.pairsMu.Unlock()

		this.mu.Lock()
		for _, ch := range this.chans {
			ch := ch
			go func() {
				defer func() {
					recover()
				}()
				select {
				case ch <- pairs:
				case <-time.After(time.Minute):
					this.options.Log.Warn("chan is full and new change has been dropped")
				}
			}()
		}
		this.mu.Unlock()
	}
	this.options.Log.Info("close watch coroutine")
}

///编解码***********************************************************************
func (this *Discovery) Marshal(v interface{}) ([]byte, error) {
	if this.options.Codec != nil {
		return this.options.Codec.Marshal(v)
	} else {
		return json.Marshal(v)
	}
}
func (this *Discovery) Unmarshal(data []byte, v interface{}) error {
	if this.options.Codec != nil {
		return this.options.Codec.Unmarshal(data, v)
	} else {
		return json.Unmarshal(data, v)
	}
}
