package discovery

import (
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpcl/core"
)

type Discovery struct {
	sys                     core.ISys
	store                   core.IStore
	RetriesAfterWatchFailed int //-1 表示一直重试 0 表示不重试
	mu                      sync.Mutex
	pairsMu                 sync.RWMutex
	pairs                   []*core.KVPair
	chans                   []chan []*core.KVPair
	stopregisterSignal      chan struct{}
	stopwatchSignal         chan struct{}
}

func (this *Discovery) GetServices() (services []*core.ServiceNode) {
	var (
		err error
	)
	this.pairsMu.RLock()
	services = make([]*core.ServiceNode, len(this.pairs))
	for _, v := range this.pairs {
		node := &core.ServiceNode{}
		if err = this.sys.Decoder().Decoder(v.Value, node); err != nil {
			this.sys.Errorf("Decoder err:%v", err)
		}
	}
	this.pairsMu.RUnlock()
	return
}

func (this *Discovery) Start() (err error) {
	var (
		ps []*core.KVPair
	)
	ps, err = this.store.List(this.sys.GetBasePath())
	if err != nil && err != core.ErrKeyNotFound {
		log.Infof("cannot get services of from registry: %v, err: %v", this.sys.GetBasePath(), err)
		return err
	}

	pairs := make([]*core.KVPair, 0, len(ps))
	prefix := this.sys.GetBasePath() + "/"
	for _, p := range ps {
		if !strings.HasPrefix(p.Key, prefix) { // avoid prefix issue of consul List
			continue
		}
		k := strings.TrimPrefix(p.Key, prefix)
		pair := &core.KVPair{Key: k, Value: p.Value}

		pairs = append(pairs, pair)
	}
	this.pairsMu.Lock()
	this.pairs = pairs
	this.pairsMu.Unlock()
	if this.sys.GetUpdateInterval() > 0 {
		go func() {
			ticker := time.NewTicker(this.sys.GetUpdateInterval())
			defer ticker.Stop()
			defer this.store.Close()
		locp:
			for {
				select {
				case <-this.stopregisterSignal:
					break locp
				case <-ticker.C:
					d, _ := this.sys.Encoder().Encoder(this.sys.GetServiceNode())
					nodePath := this.sys.GetNodePath()
					err := this.store.Put(nodePath, d, &core.WriteOptions{TTL: this.sys.GetUpdateInterval() * 2})
					if err != nil {
						log.Errorf("cannot re-create consul path %s: %v", nodePath, err)
					}
				}
			}
			this.sys.Debugf("close Timed registration coroutine")
		}()
	}
	return nil
}

func (this *Discovery) Stop() error {
	_ = this.store.Delete(this.sys.GetNodePath())
	this.stopregisterSignal <- struct{}{}
	this.stopwatchSignal <- struct{}{}
	this.sys.Debugf("Stop End !")
	return nil
}

//监控
func (this *Discovery) watch() {
	defer func() {
		this.store.Close()
	}()

	var err error
	var c <-chan []*core.KVPair
	var tempDelay time.Duration

	retry := this.RetriesAfterWatchFailed
	for this.RetriesAfterWatchFailed < 0 || retry >= 0 {
		c, err = this.store.WatchTree(this.sys.GetBasePath(), this.stopwatchSignal)
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
			this.sys.Warnf("can not watchtree (with retry %d, sleep %v): %s: %v", retry, tempDelay, this.sys.GetNodePath(), err)
			time.Sleep(tempDelay)
			continue
		}
		break
	}

	if err != nil {
		this.sys.Errorf("can't watch %s: %v", this.sys.GetBasePath(), err)
		return
	}

	prefix := this.sys.GetBasePath() + "/"
	for ps := range c {
		var pairs []*core.KVPair // latest servers
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
			k := strings.TrimPrefix(p.Key, prefix)
			pair := &core.KVPair{Key: k, Value: p.Value}
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
					this.sys.Warnf("chan is full and new change has been dropped")
				}
			}()
		}
		this.mu.Unlock()
	}
	this.sys.Infof("close watch coroutine")
}
