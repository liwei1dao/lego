package zookeeper

import (
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/rpc"
	"github.com/liwei1dao/lego/sys/rpc/discovery"
	"github.com/smallnest/rpcx/log"
)

// ZookeeperDiscovery is a zoopkeer service dcore.
// It always returns the registered servers in zookeeper.
type ZookeeperDiscovery struct {
	options  *Options
	basePath string
	kv       discovery.IStore
	pairsMu  sync.RWMutex
	pairs    []*rpc.KV
	chans    []chan []*rpc.KV
	mu       sync.Mutex

	// -1 means it always retry to watch until zookeeper is ok, 0 means no retry.
	RetriesAfterWatchFailed int

	filter discovery.ServiceDiscoveryFilter

	stopCh chan struct{}
}

// NewZookeeperDiscovery returns a new ZookeeperDiscovery.
func NewZookeeperDiscovery(basePath string, servicePath string, options *Options) (*ZookeeperDiscovery, error) {
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	if len(basePath) > 1 && strings.HasSuffix(basePath, "/") {
		basePath = basePath[:len(basePath)-1]
	}

	kv, err := NewStore(options)
	if err != nil {
		log.Infof("cannot create store: %v", err)
		return nil, err
	}

	return NewZookeeperDiscoveryWithStore(basePath+"/"+servicePath, kv)
}

// NewZookeeperDiscoveryWithStore returns a new ZookeeperDiscovery with specified store.
func NewZookeeperDiscoveryWithStore(basePath string, kv discovery.IStore) (*ZookeeperDiscovery, error) {
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}
	d := &ZookeeperDiscovery{basePath: basePath, kv: kv}
	d.stopCh = make(chan struct{})

	ps, err := kv.List(basePath)
	if err != nil {
		log.Infof("cannot get services of %s from registry: %v, err: %v", basePath, err)
		return nil, err
	}

	pairs := make([]*rpc.KV, 0, len(ps))
	for _, p := range ps {
		pair := &rpc.KV{Key: p.Key, Value: string(p.Value)}
		if d.filter != nil && !d.filter(pair) {
			continue
		}
		pairs = append(pairs, pair)
	}
	d.pairsMu.Lock()
	d.pairs = pairs
	d.pairsMu.Unlock()
	d.RetriesAfterWatchFailed = -1
	go d.watch()

	return d, nil
}

// NewZookeeperDiscoveryTemplate returns a new ZookeeperDiscovery template.
func NewZookeeperDiscoveryTemplate(basePath string, options *Options) (*ZookeeperDiscovery, error) {
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	if len(basePath) > 1 && strings.HasSuffix(basePath, "/") {
		basePath = basePath[:len(basePath)-1]
	}

	kv, err := NewStore(options)
	if err != nil {
		log.Infof("cannot create store: %v", err)
		return nil, err
	}

	return NewZookeeperDiscoveryWithStore(basePath, kv)
}

// Clone clones this ServiceDiscovery with new servicePath.
func (d *ZookeeperDiscovery) Clone(servicePath string) (discovery.ServiceDiscovery, error) {
	return NewZookeeperDiscoveryWithStore(d.basePath+"/"+servicePath, d.kv)
}

// SetFilter sets the filer.
func (d *ZookeeperDiscovery) SetFilter(filter discovery.ServiceDiscoveryFilter) {
	d.filter = filter
}

// GetServices returns the servers
func (d *ZookeeperDiscovery) GetServices() []*rpc.KV {
	d.pairsMu.RLock()
	defer d.pairsMu.RUnlock()

	return d.pairs
}

// WatchService returns a nil chan.
func (d *ZookeeperDiscovery) WatchService() chan []*rpc.KV {
	d.mu.Lock()
	defer d.mu.Unlock()

	ch := make(chan []*rpc.KV, 10)
	d.chans = append(d.chans, ch)
	return ch
}

func (d *ZookeeperDiscovery) RemoveWatcher(ch chan []*rpc.KV) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var chans []chan []*rpc.KV
	for _, c := range d.chans {
		if c == ch {
			continue
		}

		chans = append(chans, c)
	}

	d.chans = chans
}

func (d *ZookeeperDiscovery) watch() {
	defer func() {
		d.kv.Close()
	}()

	for {
		var err error
		var c <-chan []*discovery.KVPair
		var tempDelay time.Duration

		retry := d.RetriesAfterWatchFailed
		for d.RetriesAfterWatchFailed < 0 || retry >= 0 {
			c, err = d.kv.WatchTree(d.basePath, nil)
			if err != nil {
				if d.RetriesAfterWatchFailed > 0 {
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
				log.Warnf("can not watchtree (with retry %d, sleep %v): %s: %v", retry, tempDelay, d.basePath, err)
				time.Sleep(tempDelay)
				continue
			}
			break
		}

		if err != nil {
			log.Errorf("can't watch %s: %v", d.basePath, err)
			return
		}

		if err != nil {
			log.Fatalf("can not watchtree: %s: %v", d.basePath, err)
		}

	readChanges:
		for {
			select {
			case <-d.stopCh:
				log.Info("dcore has been closed")
				return
			case ps, ok := <-c:
				if !ok {
					break readChanges
				}
				var pairs []*rpc.KV // latest servers
				if ps == nil {
					d.pairsMu.Lock()
					d.pairs = pairs
					d.pairsMu.Unlock()
					continue
				}
				for _, p := range ps {
					pair := &rpc.KV{Key: p.Key, Value: string(p.Value)}
					if d.filter != nil && !d.filter(pair) {
						continue
					}
					pairs = append(pairs, pair)
				}
				d.pairsMu.Lock()
				d.pairs = pairs
				d.pairsMu.Unlock()

				d.mu.Lock()
				for _, ch := range d.chans {
					ch := ch
					go func() {
						defer func() {
							recover()
						}()
						select {
						case ch <- pairs:
						case <-time.After(time.Minute):
							log.Warn("chan is full and new change has been dropped")
						}
					}()
				}
				d.mu.Unlock()
			}
		}

		log.Warn("chan is closed and will rewatch")
	}
}

func (d *ZookeeperDiscovery) Close() {
	close(d.stopCh)
}
