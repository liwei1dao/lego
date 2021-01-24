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

func newSys(options Options) (sys *Consulregistry, err error) {
	sys = &Consulregistry{
		options: options,
	}
	config := api.DefaultConfig()
	if len(options.ConsulAddr) > 0 {
		config.Address = sys.options.ConsulAddr
	}
	if sys.options.Timeout > 0 {
		config.HttpClient.Timeout = sys.options.Timeout
	}
	if sys.client, err = api.NewClient(config); err != nil {
		return nil, err
	}
	if sys.snodeWP, err = watch.Parse(map[string]interface{}{"type": "services"}); err != nil {
		return nil, err
	}
	sys.snodeWP.Handler = sys.shandler
	sys.options.ConsulAddr = config.Address
	sys.shash = make(map[string]uint64)
	sys.services = make(map[string]*ServiceNode)
	sys.rpcsubs = make(map[core.Rpc_Key][]*ServiceNode)
	sys.watchers = make(map[string]*watch.Plan)
	sys.closeSig = make(chan bool)
	sys.isstart = false
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
		Tag:          this.options.Service.GetTag(),
		Id:           this.options.Service.GetId(),
		Type:         this.options.Service.GetType(),
		Category:     this.options.Service.GetCategory(),
		Version:      this.options.Service.GetVersion(),
		RpcId:        this.options.Service.GetRpcId(),
		PreWeight:    this.options.Service.GetPreWeight(),
		RpcSubscribe: this.getRpcInfo(),
	}); err != nil {
		return
	}
	go this.snodeWP.Run(this.options.ConsulAddr)
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
		Tags:  []string{this.options.Service.GetTag()},
		Check: check,
		Meta: map[string]string{
			"tag":          snode.Tag,
			"category":     string(snode.Category),
			"version":      fmt.Sprintf("%e", snode.Version),
			"rpcid":        snode.RpcId,
			"preweight":    fmt.Sprintf("%e", snode.PreWeight),
			"rpcsubscribe": string(rpcsubscribe),
		},
	}
	if err := this.client.Agent().ServiceRegister(asr); err != nil {
		return err
	}
	return this.client.Agent().PassTTL("service:"+snode.Id, "")
}

func (this *Consulregistry) deregisterSNode() (err error) {
	return this.client.Agent().ServiceDeregister(this.options.Service.GetId())
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
			Tag:          this.options.Service.GetTag(),
			Id:           this.options.Service.GetId(),
			Type:         this.options.Service.GetType(),
			Category:     this.options.Service.GetCategory(),
			Version:      this.options.Service.GetVersion(),
			RpcId:        this.options.Service.GetRpcId(),
			PreWeight:    this.options.Service.GetPreWeight(),
			RpcSubscribe: this.getRpcInfo(),
		})
	}
	return
}

func (this *Consulregistry) getRpcInfo() (rfs []core.Rpc_Key) {
	rpcSubscribe := rpc.GetRpcInfo()
	a := sort.StringSlice{}
	for _, v := range rpcSubscribe {
		a = append(a, string(v))
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
func (this *Consulregistry) GetServiceByType(sType string) (n []*ServiceNode) {
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

func (this *Consulregistry) GetAllServices() (n []*ServiceNode) {
	this.slock.RLock()
	defer this.slock.RUnlock()
	n = make([]*ServiceNode, 0)
	for _, v := range this.services {
		n = append(n, v)
	}
	return
}

func (this *Consulregistry) GetServiceByCategory(category core.S_Category) (n []*ServiceNode) {
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
			if v[0] == this.options.Service.GetTag() {
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
					go wp.Run(this.options.ConsulAddr)
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
			if v.Tags[0] == this.options.Service.GetTag() { //自能读取相同标签的服务
				this.addandupdataServiceNode(v)
			}
		}
	}
	return err
}
func (this *Consulregistry) addandupdataServiceNode(as *api.AgentService) (sn *ServiceNode, err error) {
	tag := as.Meta["tag"]
	category := as.Meta["category"]
	version, err := strconv.ParseFloat(as.Meta["version"], 32)
	if err != nil {
		log.Errorf("registry 读取服务节点异常:%s", as.Meta["version"])
		return
	}
	rpcid := as.Meta["rpcid"]
	preweight, err := strconv.ParseFloat(as.Meta["preweight"], 64)
	if err != nil {
		log.Errorf("registry 读取服务节点异常:%s err:%v", as.Meta["preweight"], err)
		return
	}

	snode := &ServiceNode{
		Tag:          tag,
		Type:         as.Service,
		Category:     core.S_Category(category),
		Id:           as.ID,
		Version:      float32(version),
		RpcId:        rpcid,
		PreWeight:    preweight,
		RpcSubscribe: make([]core.Rpc_Key, 0),
	}
	json.Unmarshal([]byte(as.Meta["rpcsubscribe"]), &snode.RpcSubscribe)
	if h, err := hash.Hash(snode, nil); err == nil {
		_, ok := this.services[snode.Id]
		if !ok {
			log.Infof("发现新的服务【%s】", snode.Id)
			this.services[snode.Id] = snode
			this.shash[snode.Id] = h
			if this.options.Listener != nil { //异步通知
				go this.options.Listener.FindServiceHandlefunc(*snode)
			}
		} else {
			if v, ok := this.shash[snode.Id]; !ok || v != h { //校验不一致
				// log.Debugf("更新服务【%s】", snode.Id)
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
	log.Infof("丢失服务【%s】", sId)
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
