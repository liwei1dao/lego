package etcd2

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
	"github.com/rpcxio/rpcx-etcd/store/etcd"
	"github.com/smallnest/rpcx/log"
)

// EtcdRegisterPlugin implements etcd registry.
type EtcdRegisterPlugin struct {
	// service address, for example, tcp@127.0.0.1:8972, quic@127.0.0.1:1234
	ServiceAddress string
	// etcd addresses
	EtcdServers []string
	// base path for rpcx server, for example com/example/rpcx
	BasePath string
	Metrics  metrics.Registry
	// Registered services
	Services       []string
	metasLock      sync.RWMutex
	metas          map[string]string
	UpdateInterval time.Duration
	Expired        time.Duration

	Options *Options
	kv      discovery.IStore

	dying chan struct{}
	done  chan struct{}
}

// Start starts to connect etcd cluster
func (this *EtcdRegisterPlugin) Start() error {
	if this.Expired == 0 {
		this.Expired = this.UpdateInterval
	}

	if this.done == nil {
		this.done = make(chan struct{})
	}
	if this.dying == nil {
		this.dying = make(chan struct{})
	}

	if this.kv == nil {
		kv, err := NewStore(this.Options)
		if err != nil {
			log.Errorf("cannot create etcd registry: %v", err)
			return err
		}
		this.kv = kv
	}

	err := this.kv.Put(this.BasePath, []byte("rpcx_path"), &discovery.WriteOptions{IsDir: true, TTL: this.UpdateInterval + this.Expired})
	if err != nil && !strings.Contains(err.Error(), "Not a file") {
		log.Errorf("cannot create etcd path %s: %v", this.BasePath, err)
		return err
	}

	if this.UpdateInterval > 0 {
		ticker := time.NewTicker(this.UpdateInterval)
		go func() {
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
						nodePath := fmt.Sprintf("%s/%s/%s", this.BasePath, name, this.ServiceAddress)
						kvPair, err := this.kv.Get(nodePath)
						if err != nil {
							log.Infof("can't get data of node: %s, because of %v", nodePath, err.Error())

							this.metasLock.RLock()
							meta := this.metas[name]
							this.metasLock.RUnlock()

							err = this.kv.Put(nodePath, []byte(meta), &discovery.WriteOptions{TTL: this.UpdateInterval + this.Expired})
							if err != nil {
								log.Errorf("cannot re-create etcd path %s: %v", nodePath, err)
							}

						} else {
							v, _ := url.ParseQuery(string(kvPair.Value))
							for key, value := range extra {
								v.Set(key, value)
							}
							this.kv.Put(nodePath, []byte(v.Encode()), &discovery.WriteOptions{TTL: this.UpdateInterval + this.Expired})
						}
					}
				}
			}
		}()
	}

	return nil
}

// Stop unregister all services.
func (this *EtcdRegisterPlugin) Stop() error {
	if this.kv == nil {
		kv, err := NewStore(this.Options)
		if err != nil {
			log.Errorf("cannot create etcd registry: %v", err)
			return err
		}
		this.kv = kv
	}

	for _, name := range this.Services {
		nodePath := fmt.Sprintf("%s/%s/%s", this.BasePath, name, this.ServiceAddress)
		exist, err := this.kv.Exists(nodePath)
		if err != nil {
			log.Errorf("cannot delete path %s: %v", nodePath, err)
			continue
		}
		if exist {
			this.kv.Delete(nodePath)
			log.Infof("delete path %s", nodePath, err)
		}
	}

	close(this.dying)
	<-this.done
	return nil
}

// HandleConnAccept handles connections from clients
func (this *EtcdRegisterPlugin) HandleConnAccept(conn net.Conn) (net.Conn, bool) {
	if this.Metrics != nil {
		metrics.GetOrRegisterMeter("connections", this.Metrics).Mark(1)
	}
	return conn, true
}

// PreCall handles rpc call from clients
func (this *EtcdRegisterPlugin) PreCall(_ context.Context, _, _ string, args interface{}) (interface{}, error) {
	if this.Metrics != nil {
		metrics.GetOrRegisterMeter("calls", this.Metrics).Mark(1)
	}
	return args, nil
}

// Register handles registering event.
// this service is registered at BASE/serviceName/thisIpAddress node
func (this *EtcdRegisterPlugin) Register(name string, rcvr interface{}, metadata string) (err error) {
	if strings.TrimSpace(name) == "" {
		err = errors.New("Register service `name` can't be empty")
		return
	}

	if this.kv == nil {
		etcd.Register()
		kv, err := NewStore(this.Options)
		if err != nil {
			log.Errorf("cannot create etcd registry: %v", err)
			return err
		}
		this.kv = kv
	}

	err = this.kv.Put(this.BasePath, []byte("rpcx_path"), &discovery.WriteOptions{IsDir: true})
	if err != nil && !strings.Contains(err.Error(), "Not a file") {
		log.Errorf("cannot create etcd path %s: %v", this.BasePath, err)
		return err
	}

	nodePath := fmt.Sprintf("%s/%s", this.BasePath, name)
	err = this.kv.Put(nodePath, []byte(name), &discovery.WriteOptions{IsDir: true})
	if err != nil && !strings.Contains(err.Error(), "Not a file") {
		log.Errorf("cannot create etcd path %s: %v", nodePath, err)
		return err
	}

	nodePath = fmt.Sprintf("%s/%s/%s", this.BasePath, name, this.ServiceAddress)
	err = this.kv.Put(nodePath, []byte(metadata), &discovery.WriteOptions{TTL: this.UpdateInterval + this.Expired})
	if err != nil {
		log.Errorf("cannot create etcd path %s: %v", nodePath, err)
		return err
	}

	services := make(map[string]struct{})
	for _, v := range this.Services {
		services[v] = struct{}{}
	}

	if _, ok := services[name]; !ok {
		this.Services = append(this.Services, name)
	}

	this.metasLock.Lock()
	if this.metas == nil {
		this.metas = make(map[string]string)
	}
	this.metas[name] = metadata
	this.metasLock.Unlock()
	return
}

func (this *EtcdRegisterPlugin) RegisterFunction(serviceName, fname string, fn interface{}, metadata string) error {
	return this.Register(serviceName, fn, metadata)
}

func (this *EtcdRegisterPlugin) Unregister(name string) (err error) {
	if len(this.Services) == 0 {
		return nil
	}

	if strings.TrimSpace(name) == "" {
		err = errors.New("Register service `name` can't be empty")
		return
	}

	if this.kv == nil {
		etcd.Register()
		kv, err := NewStore(this.Options)
		if err != nil {
			log.Errorf("cannot create etcd registry: %v", err)
			return err
		}
		this.kv = kv
	}

	err = this.kv.Put(this.BasePath, []byte("rpcx_path"), &discovery.WriteOptions{IsDir: true})
	if err != nil && !strings.Contains(err.Error(), "Not a file") {
		log.Errorf("cannot create etcd path %s: %v", this.BasePath, err)
		return err
	}

	nodePath := fmt.Sprintf("%s/%s", this.BasePath, name)
	err = this.kv.Put(nodePath, []byte(name), &discovery.WriteOptions{IsDir: true})
	if err != nil && !strings.Contains(err.Error(), "Not a file") {
		log.Errorf("cannot create etcd path %s: %v", nodePath, err)
		return err
	}

	nodePath = fmt.Sprintf("%s/%s/%s", this.BasePath, name, this.ServiceAddress)

	err = this.kv.Delete(nodePath)
	if err != nil {
		log.Errorf("cannot create consul path %s: %v", nodePath, err)
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
	return
}
