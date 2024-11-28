package zookeeper

import (
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/discovery/dcore"
	"github.com/smallnest/rpcx/client"
	"github.com/smallnest/rpcx/log"
)

func init() {
	Register()
}

// ZookeeperDiscovery is a zoopkeer service discovery.
// It always returns the registered servers in zookeeper.
type ZookeeperDiscovery struct {
	basePath string
	kv       dcore.IStore
	pairsMu  sync.RWMutex
	pairs    []*client.KVPair
	chans    []chan []*client.KVPair
	mu       sync.Mutex

	// -1 means it always retry to watch until zookeeper is ok, 0 means no retry.
	RetriesAfterWatchFailed int

	filter client.ServiceDiscoveryFilter

	stopCh chan struct{}
}

// NewZookeeperDiscovery returns a new ZookeeperDiscovery.
func NewZookeeperDiscovery(basePath string, servicePath string, zkAddr []string, options *dcore.Config) (*ZookeeperDiscovery, error) {
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	if len(basePath) > 1 && strings.HasSuffix(basePath, "/") {
		basePath = basePath[:len(basePath)-1]
	}

	kv, err := dcore.NewStore(dcore.ZK, zkAddr, options)
	if err != nil {
		log.Infof("cannot create store: %v", err)
		return nil, err
	}

	return NewZookeeperDiscoveryWithStore(basePath+"/"+servicePath, kv)
}

// NewZookeeperDiscoveryWithStore returns a new ZookeeperDiscovery with specified store.
func NewZookeeperDiscoveryWithStore(basePath string, kv dcore.IStore) (*ZookeeperDiscovery, error) {
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

	pairs := make([]*client.KVPair, 0, len(ps))
	for _, p := range ps {
		pair := &client.KVPair{Key: p.Key, Value: string(p.Value)}
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
func NewZookeeperDiscoveryTemplate(basePath string, zkAddr []string, options *dcore.Config) (*ZookeeperDiscovery, error) {
	if basePath[0] == '/' {
		basePath = basePath[1:]
	}

	if len(basePath) > 1 && strings.HasSuffix(basePath, "/") {
		basePath = basePath[:len(basePath)-1]
	}

	kv, err := dcore.NewStore(dcore.ZK, zkAddr, options)
	if err != nil {
		log.Infof("cannot create store: %v", err)
		return nil, err
	}

	return NewZookeeperDiscoveryWithStore(basePath, kv)
}

// Clone clones this ServiceDiscovery with new servicePath.
func (d *ZookeeperDiscovery) Clone(servicePath string) (client.ServiceDiscovery, error) {
	return NewZookeeperDiscoveryWithStore(d.basePath+"/"+servicePath, d.kv)
}

// SetFilter sets the filer.
func (d *ZookeeperDiscovery) SetFilter(filter client.ServiceDiscoveryFilter) {
	d.filter = filter
}

// GetServices returns the servers
func (d *ZookeeperDiscovery) GetServices() []*client.KVPair {
	d.pairsMu.RLock()
	defer d.pairsMu.RUnlock()

	return d.pairs
}

// WatchService returns a nil chan.
func (d *ZookeeperDiscovery) WatchService() chan []*client.KVPair {
	d.mu.Lock()
	defer d.mu.Unlock()

	ch := make(chan []*client.KVPair, 10)
	d.chans = append(d.chans, ch)
	return ch
}

func (d *ZookeeperDiscovery) RemoveWatcher(ch chan []*client.KVPair) {
	d.mu.Lock()
	defer d.mu.Unlock()

	var chans []chan []*client.KVPair
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
		var c <-chan []*dcore.KVPair
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
				log.Info("discovery has been closed")
				return
			case ps, ok := <-c:
				if !ok {
					break readChanges
				}
				var pairs []*client.KVPair // latest servers
				if ps == nil {
					d.pairsMu.Lock()
					d.pairs = pairs
					d.pairsMu.Unlock()
					continue
				}
				for _, p := range ps {
					pair := &client.KVPair{Key: p.Key, Value: string(p.Value)}
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
