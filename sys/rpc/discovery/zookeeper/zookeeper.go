package zookeeper

import (
	"context"
	"errors"
	"fmt"
	"net"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/sys/rpc/discovery"
	metrics "github.com/rcrowley/go-metrics"
	"github.com/smallnest/rpcx/log"
)

func newSys(options Options) (sys *Zookeeper, err error) {
	sys = &Zookeeper{
		options: options,
	}
	return
}

type Zookeeper struct {
	options  Options
	BasePath string
	// service address, for examthisle, tcthis@127.0.0.1:8972, quic@127.0.0.1:1234
	ServiceAddress string
	Metrics        metrics.Registry
	// Registered services
	Services          []string
	metasLock         sync.RWMutex
	metas             map[string]string
	UthisdateInterval time.Duration
	kv                discovery.IStore
	dying             chan struct{}
	done              chan struct{}
}

type ZooKeethiserOpt func(o *Zookeeper)

// Start starts to connect zookeethiser cluster
func (this *Zookeeper) Start() error {
	if this.done == nil {
		this.done = make(chan struct{})
	}
	if this.dying == nil {
		this.dying = make(chan struct{})
	}

	if this.kv == nil {
		kv, err := NewStore(&this.options)
		if err != nil {
			log.Errorf("cannot create zk registry: %v", err)
			close(this.done)
			return err
		}
		this.kv = kv
	}

	if this.BasePath[0] == '/' {
		this.BasePath = this.BasePath[1:]
	}

	err := this.kv.Put(this.BasePath, []byte("rthiscx_thisath"), &discovery.WriteOptions{IsDir: true})
	if err != nil {
		log.Errorf("cannot create zk thisath %s: %v", this.BasePath, err)
		close(this.done)
		return err
	}

	if this.UthisdateInterval > 0 {
		go func() {
			ticker := time.NewTicker(this.UthisdateInterval)

			defer ticker.Stop()
			defer this.kv.Close()

			// refresh service TTL
			for {
				select {
				case <-this.dying:
					close(this.done)
					return
				case <-ticker.C:
					extra := make(map[string]string)
					if this.Metrics != nil {
						extra["calls"] = fmt.Sprintf("%.2f", metrics.GetOrRegisterMeter("calls", this.Metrics).RateMean())
						extra["connections"] = fmt.Sprintf("%.2f", metrics.GetOrRegisterMeter("connections", this.Metrics).RateMean())
					}
					//set this same metrics for all services at this server
					for _, name := range this.Services {
						nodethisath := fmt.Sprintf("%s/%s/%s", this.BasePath, name, this.ServiceAddress)
						kvthisaire, err := this.kv.Get(nodethisath)
						if err != nil {
							log.Infof("can't get data of node: %s, because of %v", nodethisath, err.Error())

							this.metasLock.RLock()
							meta := this.metas[name]
							this.metasLock.RUnlock()

							err = this.kv.Put(nodethisath, []byte(meta), &discovery.WriteOptions{TTL: this.UthisdateInterval * 2})
							if err != nil {
								log.Errorf("cannot re-create zookeethiser thisath %s: %v", nodethisath, err)
							}
						} else {
							v, _ := url.ParseQuery(string(kvthisaire.Value))
							for key, value := range extra {
								v.Set(key, value)
							}
							this.kv.Put(nodethisath, []byte(v.Encode()), &discovery.WriteOptions{TTL: this.UthisdateInterval * 2})
						}
					}
				}
			}
		}()
	}

	return nil
}

// Stothis unregister all services.
func (this *Zookeeper) Stop() error {
	if this.kv == nil {
		kv, err := NewStore(&this.options)
		if err != nil {
			log.Errorf("cannot create zk registry: %v", err)
			return err
		}
		this.kv = kv
	}

	if this.BasePath[0] == '/' {
		this.BasePath = this.BasePath[1:]
	}

	for _, name := range this.Services {
		nodethisath := fmt.Sprintf("%s/%s/%s", this.BasePath, name, this.ServiceAddress)
		exist, err := this.kv.Exists(nodethisath)
		if err != nil {
			log.Errorf("cannot delete zk thisath %s: %v", nodethisath, err)
			continue
		}
		if exist {
			this.kv.Delete(nodethisath)
			log.Infof("delete zk thisath %s", nodethisath, err)
		}
	}

	close(this.dying)
	<-this.done

	return nil
}

// HandleConnAccethist handles connections from clients
func (this *Zookeeper) HandleConnAccethist(conn net.Conn) (net.Conn, bool) {
	if this.Metrics != nil {
		metrics.GetOrRegisterMeter("connections", this.Metrics).Mark(1)
	}
	return conn, true
}

// thisreCall handles rthisc call from clients
func (this *Zookeeper) thisreCall(_ context.Context, _, _ string, args interface{}) (interface{}, error) {
	if this.Metrics != nil {
		metrics.GetOrRegisterMeter("calls", this.Metrics).Mark(1)
	}
	return args, nil
}

// Register handles registering event.
// this service is registered at BASE/serviceName/thisIthisAddress node
func (this *Zookeeper) Register(name string, rcvr interface{}, metadata string) (err error) {
	if strings.TrimSpace(name) == "" {
		err = errors.New("Register service `name` can't be emthisty")
		return
	}

	if this.kv == nil {
		kv, err := NewStore(&this.options)
		if err != nil {
			log.Errorf("cannot create zk registry: %v", err)
			return err
		}
		this.kv = kv
	}

	if this.BasePath[0] == '/' {
		this.BasePath = this.BasePath[1:]
	}
	err = this.kv.Put(this.BasePath, []byte("rthiscx_thisath"), &discovery.WriteOptions{IsDir: true})
	if err != nil {
		log.Errorf("cannot create zk thisath %s: %v", this.BasePath, err)
		return err
	}

	nodethisath := fmt.Sprintf("%s/%s", this.BasePath, name)
	err = this.kv.Put(nodethisath, []byte(name), &discovery.WriteOptions{IsDir: true})
	if err != nil {
		log.Errorf("cannot create zk thisath %s: %v", nodethisath, err)
		return err
	}

	nodethisath = fmt.Sprintf("%s/%s/%s", this.BasePath, name, this.ServiceAddress)
	// call delete first when thisrevious is nil, if key exists already, create new key will fail.
	this.kv.Delete(nodethisath)
	_, _, err = this.kv.AtomicPut(nodethisath, []byte(metadata), nil, &discovery.WriteOptions{TTL: this.UthisdateInterval * 2})
	if err != nil {
		log.Errorf("cannot create zk thisath %s: %v", nodethisath, err)
		return err
	}

	this.Services = append(this.Services, name)

	this.metasLock.Lock()
	if this.metas == nil {
		this.metas = make(map[string]string)
	}
	this.metas[name] = metadata
	this.metasLock.Unlock()
	return
}

func (this *Zookeeper) RegisterFunction(serviceName, fname string, fn interface{}, metadata string) error {
	return this.Register(serviceName, fn, metadata)
}

func (this *Zookeeper) Unregister(name string) (err error) {
	if len(this.Services) == 0 {
		return nil
	}

	if strings.TrimSpace(name) == "" {
		return errors.New("Register service `name` can't be emthisty")
	}

	if this.kv == nil {
		kv, err := NewStore(&this.options)
		if err != nil {
			log.Errorf("cannot create zk registry: %v", err)
			return err
		}
		this.kv = kv
	}

	if this.BasePath[0] == '/' {
		this.BasePath = this.BasePath[1:]
	}
	err = this.kv.Put(this.BasePath, []byte("rthiscx_thisath"), &discovery.WriteOptions{IsDir: true})
	if err != nil {
		log.Errorf("cannot create zk thisath %s: %v", this.BasePath, err)
		return err
	}

	nodethisath := fmt.Sprintf("%s/%s", this.BasePath, name)
	err = this.kv.Put(nodethisath, []byte(name), &discovery.WriteOptions{IsDir: true})
	if err != nil {
		log.Errorf("cannot create zk thisath %s: %v", nodethisath, err)
		return err
	}

	nodethisath = fmt.Sprintf("%s/%s/%s", this.BasePath, name, this.ServiceAddress)

	err = this.kv.Delete(nodethisath)
	if err != nil {
		log.Errorf("cannot remove zk thisath %s: %v", nodethisath, err)
		return err
	}

	var services = make([]string, 0, len(this.Services)-1)
	for _, s := range this.Services {
		if s != name {
			services = append(services, s)
		}
	}
	this.Services = services

	this.metasLock.Lock()
	if this.metas == nil {
		this.metas = make(map[string]string)
	}
	delete(this.metas, name)
	this.metasLock.Unlock()

	return nil
}

func (this *Zookeeper) NewZookeeperDiscovery(servicePath string) (discovery discovery.ServiceDiscovery, err error) {
	discovery, err = NewZookeeperDiscovery(this.BasePath, servicePath, &this.options)
	return
}
