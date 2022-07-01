package consul

import (
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/rpcl/core"
)

func NewConsulDiscovery() (discovery *ConsulDiscovery, err error) {
	return
}

type ConsulDiscovery struct {
	sys                     core.ISys
	store                   core.IStore
	pairsMu                 sync.RWMutex
	RetriesAfterWatchFailed int //-1 表示一直重试 0 表示不重试
	basePath                string
	mu                      sync.Mutex
	pairs                   []*core.KVPair
	chans                   []chan []*core.KVPair
	stopCh                  chan struct{}
}

func (this *ConsulDiscovery) Start() {

}

func (this *ConsulDiscovery) GetServices() []*core.KVPair {
	this.pairsMu.RLock()
	defer this.pairsMu.RUnlock()
	return this.pairs
}

//监控
func (this *ConsulDiscovery) watch() {
	defer func() {
		this.store.Close()
	}()
	for {
		var err error
		var c <-chan []*core.KVPair
		var tempDelay time.Duration

		retry := this.RetriesAfterWatchFailed
		for this.RetriesAfterWatchFailed < 0 || retry >= 0 {
			c, err = this.store.WatchTree(this.basePath, nil)
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
				this.sys.Warnf("can not watchtree (with retry %d, sleep %v): %s: %v", retry, tempDelay, this.basePath, err)
				time.Sleep(tempDelay)
				continue
			}
			break
		}
		if err != nil {
			this.sys.Errorf("can't watch %s: %v", this.basePath, err)
			return
		}

		prefix := this.basePath + "/"
	readChanges:
		for {
			select {
			case <-this.stopCh:
				this.sys.Infof("discovery has been closed")
				return
			case ps, ok := <-c:
				if !ok {
					break readChanges
				}
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
		}
		this.sys.Warnf("chan is closed and will rewatch")
	}
}
