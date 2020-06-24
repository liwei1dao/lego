package registry

import (
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	hash "github.com/mitchellh/hashstructure"

	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc"
)

func newConsulregistry(opts ...Option) (c *Consulregistry, err error) {
	c = new(Consulregistry)
	config := api.DefaultConfig()
	c.options = newOptions(opts...)
	if len(c.options.Address) > 0 {
		config.Address = c.options.Address
	}
	if c.options.Timeout > 0 {
		config.HttpClient.Timeout = c.options.Timeout
	}
	if c.client, err = api.NewClient(config); err != nil {
		return nil, err
	}
	if c.snodeWP, err = watch.Parse(map[string]interface{}{"type": "services"}); err != nil {
		return nil, err
	}
	c.snodeWP.Handler = c.shandler
	c.options.Address = config.Address
	c.shash = make(map[string]uint64)
	c.services = make(map[string]*ServiceNode)
	c.rpcsubs = make(map[core.Rpc_Key][]*ServiceNode)
	c.watchers = make(map[string]*watch.Plan)
	c.closeSig = make(chan bool)
	c.isstart = false
	return
}

type Consulregistry struct {
	options  Options
	client   *api.Client
	shash    map[string]uint64
	isstart  bool
	closeSig chan bool
	snodeWP  *watch.Plan
	slock    sync.RWMutex
	services map[string]*ServiceNode
	rlock    sync.RWMutex
	rpcsubs  map[core.Rpc_Key][]*ServiceNode
	watchers map[string]*watch.Plan
}

//启动
func (this *Consulregistry) Start() (err error) {
	if err = this.getServices(); err != nil {
		return
	}
	if err = this.registerSNode(&ServiceNode{
		Tag:          service.GetTag(),
		Id:           service.GetId(),
		Type:         service.GetType(),
		Category:     service.GetCategory(),
		Version:      service.GetVersion(),
		RpcId:        service.GetRpcId(),
		PreWeight:    service.GetPreWeight(),
		RpcSubscribe: this.getRpcInfo(),
	}); err != nil {
		return
	}
	go this.snodeWP.Run(this.options.Address)
	go this.run()
	this.isstart = true
	return
}

//停止
func (this *Consulregistry) Stop() (err error) {
	this.closeSig <- true
	this.isstart = false
	this.snodeWP.Stop()
	this.deregisterSNode()
	return
}

func (this *Consulregistry) registerSNode(snode *ServiceNode) (err error) {
	h, err := hash.Hash(snode, nil)
	if err != nil {
		return err
	}

	if v, ok := this.shash[snode.Id]; ok && v != 0 && v == h {
		if err := this.client.Agent().PassTTL("service:"+snode.Id, ""); err == nil {
			return nil
		}
	}
	deregTTL := getDeregisterTTL(this.options.RegisterTTL)
	check := &api.AgentServiceCheck{
		TTL:                            fmt.Sprintf("%v", this.options.RegisterTTL),
		DeregisterCriticalServiceAfter: fmt.Sprintf("%v", deregTTL),
	}
	rpcsubscribe, _ := json.Marshal(snode.RpcSubscribe)
	asr := &api.AgentServiceRegistration{
		ID:    snode.Id,
		Name:  snode.Type,
		Tags:  []string{this.options.Tag},
		Check: check,
		Meta: map[string]string{
			"tag":          snode.Tag,
			"category":     string(snode.Category),
			"version":      fmt.Sprintf("%d", snode.Version),
			"rpcid":        snode.RpcId,
			"preweight":    fmt.Sprintf("%d", snode.PreWeight),
			"rpcsubscribe": string(rpcsubscribe),
		},
	}
	if err := this.client.Agent().ServiceRegister(asr); err != nil {
		return err
	}
	return this.client.Agent().PassTTL("service:"+snode.Id, "")
}

func (this *Consulregistry) deregisterSNode() (err error) {
	return this.client.Agent().ServiceDeregister(service.GetId())
}

func (this *Consulregistry) run() {
	t := time.NewTicker(this.options.RegisterInterval)
	for {
		select {
		case <-t.C:
			err := this.PushServiceInfo()
			if err != nil {
				log.Warnf("service run Server.Register error: ", err)
			}
		case <-this.closeSig:
			return
		}
	}
	log.Debugf("退出 Consulregistry Run 函数")
}

func (this *Consulregistry) PushServiceInfo() (err error) {
	if this.isstart {
		err = this.registerSNode(&ServiceNode{
			Tag:          service.GetTag(),
			Id:           service.GetId(),
			Type:         service.GetType(),
			Category:     service.GetCategory(),
			Version:      service.GetVersion(),
			RpcId:        service.GetRpcId(),
			PreWeight:    service.GetPreWeight(),
			RpcSubscribe: this.getRpcInfo(),
		})
	}
	return
}

func (this *Consulregistry) getRpcInfo() (rfs []core.Rpc_Key) {
	rpcSubscribe := rpc.GetRpcInfo()
	a := sort.StringSlice{}
	for k, _ := range rpcSubscribe {
		a = append(a, string(k))
	}
	sort.Sort(a)
	rfs = make([]core.Rpc_Key, len(a))
	for i, k := range a {
		rfs[i] = core.Rpc_Key(k)
	}
	return
}

func (this *Consulregistry) GetServiceById(sId string) (n *ServiceNode, err error) {
	this.slock.RLock()
	n, ok := this.services[sId]
	this.slock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("No Found %s", sId)
	}
	return
}
func (this *Consulregistry) GetServiceByType(sType string) (n []*ServiceNode, err error) {
	this.slock.RLock()
	defer this.slock.RUnlock()
	n = make([]*ServiceNode, 0)
	for _, v := range this.services {
		if v.Type == sType {
			n = append(n, v)
		}
	}
	return
}
func (this *Consulregistry) GetServiceByCategory(category core.S_Category) (n []*ServiceNode, err error) {
	this.slock.RLock()
	defer this.slock.RUnlock()
	n = make([]*ServiceNode, 0)
	for _, v := range this.services {
		if v.Category == category {
			n = append(n, v)
		}
	}
	return
}

func (this *Consulregistry) GetRpcSubById(rId core.Rpc_Key) (n []*ServiceNode) {
	n = make([]*ServiceNode, 0)
	this.rlock.RLock()
	d, ok := this.rpcsubs[rId]
	this.rlock.RUnlock()
	if ok {
		n = d
	}
	return
}

func getDeregisterTTL(t time.Duration) time.Duration {
	// splay slightly for the watcher?
	splay := time.Second * 5
	deregTTL := t + splay
	// consul has a minimum timeout on deregistration of 1 minute.
	if t < time.Minute {
		deregTTL = time.Minute + splay
	}
	return deregTTL
}
func (this *Consulregistry) shandler(idx uint64, data interface{}) {
	services, ok := data.(map[string][]string)
	if !ok {
		return
	}
	for k, v := range services {
		if len(v) > 0 {
			if v[0] == service.GetTag() {
				_, ok := this.watchers[k]
				if ok {
					continue
				}
				wp, err := watch.Parse(map[string]interface{}{
					"type":    "service",
					"service": k,
				})
				if err == nil {
					wp.Handler = this.snodehandler
					go wp.Run(this.options.Address)
					this.watchers[k] = wp
				}
			}
		}
	}
	this.slock.Lock()
	defer this.slock.Unlock()
	for k, v := range this.watchers {
		if _, ok := services[k]; !ok { //不存在了
			v.Stop()
			delete(this.watchers, k)

			for k1, v1 := range this.services {
				if v1.Type == k {
					this.removeServiceNode(k1)
				}
			}
		}
	}
}
func (this *Consulregistry) snodehandler(idx uint64, data interface{}) {
	entries, ok := data.([]*api.ServiceEntry)
	if !ok {
		return
	}
	this.slock.Lock()
	defer this.slock.Unlock()
	stype := ""
	serviceMap := map[string]struct{}{}
	for _, v := range entries {
		isdele := false
		for _, check := range v.Checks {
			if check.Status == "critical" {
				this.removeServiceNode(v.Service.ID)
				isdele = true
				break
			}
		}
		if !isdele {
			this.addandupdataServiceNode(v.Service)
			stype = v.Service.Service
			serviceMap[v.Service.ID] = struct{}{}
		}
	}
	for _, v := range this.services {
		if v.Type == stype {
			if _, ok := serviceMap[v.Id]; !ok {
				this.removeServiceNode(v.Id)
			}
		}
	}
}
func (this *Consulregistry) getServices() (err error) {
	services, err := this.client.Agent().Services()
	if err == nil {
		for _, v := range services {
			if v.Tags[0] == service.GetTag() { //自能读取相同标签的服务
				this.addandupdataServiceNode(v)
			}
		}
	}
	return err
}
func (this *Consulregistry) addandupdataServiceNode(as *api.AgentService) (sn *ServiceNode, err error) {
	tag := as.Meta["tag"]
	category := as.Meta["category"]
	version, err := strconv.Atoi(as.Meta["version"])
	if err != nil {
		log.Errorf("registry 读取服务节点异常:%v", as.Meta["version"])
		return
	}
	rpcid := as.Meta["rpcid"]
	preweight, err := strconv.Atoi(as.Meta["preweight"])
	if err != nil {
		log.Errorf("registry 读取服务节点异常:%v", as.Meta["preweight"])
		return
	}

	snode := &ServiceNode{
		Tag:          tag,
		Type:         as.Service,
		Category:     core.S_Category(category),
		Id:           as.ID,
		Version:      int32(version),
		RpcId:        rpcid,
		PreWeight:    int32(preweight),
		RpcSubscribe: make([]core.Rpc_Key, 0),
	}
	json.Unmarshal([]byte(as.Meta["rpcsubscribe"]), &snode.RpcSubscribe)
	if h, err := hash.Hash(snode, nil); err == nil {
		_, ok := this.services[snode.Id]
		if !ok {
			log.Debugf("发现新的服务【%s】", snode.Id)
			this.services[snode.Id] = snode
			this.shash[snode.Id] = h
			if this.options.Listener != nil { //异步通知
				go this.options.Listener.FindServiceHandlefunc(*snode)
			}
		} else {
			if v, ok := this.shash[snode.Id]; !ok || v != h { //校验不一致
				log.Debugf("更新服务【%s】", snode.Id)
				this.services[snode.Id] = snode
				this.shash[snode.Id] = h
				if this.options.Listener != nil { //异步通知
					go this.options.Listener.UpDataServiceHandlefunc(*snode)
				}
			}
		}
	} else {
		log.Errorf("registry 校验服务 hash 值异常:%s", err.Error())
	}
	this.SyncRpcInfo()
	return
}
func (this *Consulregistry) removeServiceNode(sId string) {
	_, ok := this.services[sId]
	if !ok {
		return
	}
	log.Debugf("丢失服务【%s】", sId)
	delete(this.services, sId)
	this.SyncRpcInfo()
	if this.options.Listener != nil { //异步通知
		go this.options.Listener.LoseServiceHandlefunc(sId)
	}
}

func (this *Consulregistry) SyncRpcInfo() {
	this.rlock.Lock()
	defer this.rlock.Unlock()
	this.rpcsubs = make(map[core.Rpc_Key][]*ServiceNode)
	for _, v := range this.services {
		for _, v1 := range v.RpcSubscribe {
			if _, ok := this.rpcsubs[v1]; !ok {
				this.rpcsubs[v1] = make([]*ServiceNode, 0)
			}
			this.rpcsubs[v1] = append(this.rpcsubs[v1], v)
		}
	}
}
