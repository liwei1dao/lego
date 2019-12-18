package registry

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/api/watch"
	hash "github.com/mitchellh/hashstructure"
	"lego/core"
	"lego/sys/log"
	cont "lego/utils/concurrent"
	"strconv"
	"strings"
	"time"
)

func newConsulregistry(opts ...Option) (c *Consulregistry, err error) {
	c = new(Consulregistry)
	config := api.DefaultConfig()
	c.opts = newOptions(opts...)
	c.watchers = cont.NewBeeMap()
	c.snodes = cont.NewBeeMap()
	c.hash = 0
	c.rpcs = cont.NewBeeMap()
	if len(c.opts.Address) > 0 {
		config.Address = c.opts.Address
	}
	if c.opts.Timeout > 0 {
		config.HttpClient.Timeout = c.opts.Timeout
	}
	if c.client, err = api.NewClient(config); err != nil {
		return nil, err
	}
	if c.snodeWP, err = watch.Parse(map[string]interface{}{"type": "services"}); err != nil {
		return nil, err
	}
	if c.rpclegoP, err = watch.Parse(map[string]interface{}{
		"type":   "keyprefix",
		"prefix": c.opts.Tag + "Rpc/",
	}); err != nil {
		return nil, err
	}
	c.snodeWP.Handler = c.SHandler
	c.rpclegoP.Handler = c.RpcHandler
	c.opts.Address = config.Address
	c.getServices()
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

type Consulregistry struct {
	opts     Options
	client   *api.Client
	snodeWP  *watch.Plan
	rpclegoP *watch.Plan
	watchers *cont.BeeMap
	hash     uint64
	//shash    *cont.BeeMap
	snodes *cont.BeeMap
	rpcs   *cont.BeeMap
}

func (this *Consulregistry) Options() Options {
	return this.opts
}

func (this *Consulregistry) getServices() {
	services, err := this.client.Agent().Services()
	if err == nil {
		for _, v := range services {
			if v.Tags[0] == service.GetTag() { //自能读取相同标签的服务
				this.addandupdataServiceNode(v)
			}
		}
	}
}

func (this *Consulregistry) GetServiceById(sId string) (n *ServiceNode, err error) {
	if v := this.snodes.Get(sId); v != nil {
		return v.(*ServiceNode), nil
	} else {
		return nil, fmt.Errorf("No Found %s", sId)
	}
}
func (this *Consulregistry) GetServiceByType(sType string) (n []*ServiceNode, err error) {
	n = make([]*ServiceNode, 0)
	for _, v := range this.snodes.Items() {
		item := v.(*ServiceNode)
		if item.Type == sType {
			n = append(n, item)
		}
	}
	return
}
func (this *Consulregistry) GetServiceByCategory(category core.S_Category) (n []*ServiceNode, err error) {
	n = make([]*ServiceNode, 0)
	for _, v := range this.snodes.Items() {
		item := v.(*ServiceNode)
		if item.Category == category {
			n = append(n, item)
		}
	}
	return
}

func (this *Consulregistry) GetRpcSubById(rId core.Rpc_Key) (rf *RpcFuncInfo, err error) {
	if v := this.rpcs.Get(string(rId)); v != nil {
		return v.(*RpcFuncInfo), nil
	} else {
		return nil, fmt.Errorf("No Found Rpcid [%s]!", rId)
	}
}

//获取服务节点
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
	snode := ServiceNode{
		Tag:       tag,
		Type:      as.Service,
		Category:  core.S_Category(category),
		Id:        as.ID,
		Version:   int32(version),
		RpcId:     rpcid,
		PreWeight: int32(preweight),
	}
	sn = &snode
	if v := this.snodes.Get(snode.Id); v == nil {
		this.snodes.Set(snode.Id, sn)
		if this.opts.Find != nil { //异步通知
			go this.opts.Find(snode)
		}
	} else {
		this.snodes.Set(snode.Id, sn)
		if this.opts.UpData != nil { //异步通知
			go this.opts.UpData(snode)
		}
	}
	return
}

func (this *Consulregistry) removeServiceNode(sId string) {
	if this.snodes.Get(sId) == nil {
		return
	}
	this.snodes.Delete(sId)
	if this.opts.Lose != nil { //异步通知
		go this.opts.Lose(sId)
	}
}

//注册服务节点
func (this *Consulregistry) RegisterSNode(s *ServiceNode) (err error) {
	h, err := hash.Hash(s, nil)
	if err != nil {
		return err
	}

	if this.hash != 0 && this.hash == h {
		if err := this.client.Agent().PassTTL("service:"+s.Id, ""); err == nil {
			return nil
		}
	}

	deregTTL := getDeregisterTTL(this.opts.RegisterTTL)
	check := &api.AgentServiceCheck{
		TTL:                            fmt.Sprintf("%v", this.opts.RegisterTTL),
		DeregisterCriticalServiceAfter: fmt.Sprintf("%v", deregTTL),
	}

	asr := &api.AgentServiceRegistration{
		ID:    s.Id,
		Name:  s.Type,
		Tags:  []string{this.opts.Tag},
		Check: check,
		Meta: map[string]string{
			"tag":       s.Tag,
			"category":  string(s.Category),
			"version":   fmt.Sprintf("%d", s.Version),
			"rpcid":     s.RpcId,
			"preweight": fmt.Sprintf("%d", s.PreWeight),
		},
	}
	if err := this.client.Agent().ServiceRegister(asr); err != nil {
		return err
	}
	this.hash = h
	go this.snodeWP.Run(this.opts.Address)
	go this.rpclegoP.Run(this.opts.Address)
	return this.client.Agent().PassTTL("service:"+s.Id, "")
}

//注销服务节点
func (this *Consulregistry) DeregisterSNode(sId string) (err error) {
	this.hash = 0
	this.snodes.Delete(sId)
	for _, v := range this.rpcs.Items() {
		item := v.(*RpcFuncInfo)
		if _, ok := item.SubServiceIds[sId]; ok {
			if v, ok := item.SubServiceIds[sId]; ok {
				v.SubNum--
				if v.SubNum == 0 {
					delete(item.SubServiceIds, sId)
				}
				key := this.opts.Tag + "Rpc/" + string(item.Id)
				if len(item.SubServiceIds) > 0 {
					if value, err := json.Marshal(item); err != nil {
						return err
					} else {
						kv := &api.KVPair{
							Key:   key,
							Flags: 0,
							Value: value,
						}
						_, err = this.client.KV().Put(kv, nil)
						return err
					}
				} else {
					this.client.KV().Delete(key, nil)
				}
			}
		}
	}
	this.snodeWP.Stop()
	this.rpclegoP.Stop()
	return this.client.Agent().ServiceDeregister(sId)
}

//注册rpc订阅信息
func (this *Consulregistry) RegisterRpcSub(rId core.Rpc_Key, sId string) (err error) {
	rpcinfo := &RpcFuncInfo{Id: string(rId), SubServiceIds: map[string]*ServiceRpcSub{}}
	if rf := this.rpcs.Get(rId); rf != nil {
		rpcinfo = rf.(*RpcFuncInfo)
	}

	if v, ok := rpcinfo.SubServiceIds[sId]; ok {
		v.SubNum++
	} else {
		rpcinfo.SubServiceIds[sId] = &ServiceRpcSub{
			Id:     sId,
			SubNum: 1,
		}
	}
	if value, err := json.Marshal(rpcinfo); err != nil {
		return err
	} else {
		key := this.opts.Tag + "Rpc/" + string(rId)
		kv := &api.KVPair{
			Key:   key,
			Flags: 0,
			Value: value,
		}
		_, err = this.client.KV().Put(kv, nil)
		return err
	}
}

func (this *Consulregistry) UnRegisterRpcSub(rId core.Rpc_Key, sId string) (err error) {
	rpcinfo := &RpcFuncInfo{Id: string(rId), SubServiceIds: map[string]*ServiceRpcSub{}}
	if rf := this.rpcs.Get(rId); rf != nil {
		rpcinfo = rf.(*RpcFuncInfo)
	}

	if v, ok := rpcinfo.SubServiceIds[sId]; ok {
		v.SubNum--
		if v.SubNum == 0 {
			delete(rpcinfo.SubServiceIds, sId)
		}
		key := this.opts.Tag + "Rpc/" + string(rId)
		if len(rpcinfo.SubServiceIds) > 0 {
			if value, err := json.Marshal(rpcinfo); err != nil {
				return err
			} else {
				kv := &api.KVPair{
					Key:   key,
					Flags: 0,
					Value: value,
				}
				_, err = this.client.KV().Put(kv, nil)
				return err
			}
		} else {
			this.client.KV().Delete(key, nil)
		}
	}
	return
}

//监控服务列表变化
func (this *Consulregistry) SHandler(idx uint64, data interface{}) {
	services, ok := data.(map[string][]string)
	if !ok {
		return
	}
	for k, v := range services {
		if len(v) > 0 {
			if v[0] == service.GetTag() {
				if swh := this.watchers.Get(k); swh != nil {
					continue
				}
				wp, err := watch.Parse(map[string]interface{}{
					"type":    "service",
					"service": k,
				})
				if err == nil {
					wp.Handler = this.SNodeHandler
					go wp.Run(this.opts.Address)
					this.watchers.Set(k, wp)
				}
			}
		}
	}
	// remove unknown services from watchers
	for service, w := range this.watchers.Items() {
		if _, ok := services[service.(string)]; !ok {
			item := w.(*watch.Plan)
			item.Stop()
			this.watchers.Delete(service)
			for k, v := range this.snodes.Items() {
				item := v.(*ServiceNode)
				if item.Type == service {
					this.removeServiceNode(k.(string))
				}
			}
		}
	}
}

//监控服务列表变化
func (this *Consulregistry) SNodeHandler(idx uint64, data interface{}) {
	entries, ok := data.([]*api.ServiceEntry)
	if !ok {
		return
	}
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

	for _, v := range this.snodes.Items() {
		item := v.(*ServiceNode)
		if item.Type == stype {
			if _, ok := serviceMap[item.Id]; !ok {
				this.removeServiceNode(item.Id)
			}
		}
	}
}

//监控服务列表变化
func (this *Consulregistry) RpcHandler(idx uint64, data interface{}) {
	srpc, ok := data.(api.KVPairs)
	if !ok {
		return
	}
	if srpc != nil && len(srpc) > 0 {
		for _, v := range srpc {
			rId := v.Key[strings.LastIndex(v.Key, "/")+1:]
			rpcinfo := &RpcFuncInfo{Id: rId}
			json.Unmarshal(v.Value, &rpcinfo)
			this.rpcs.Set(rId, rpcinfo)
		}
	}
}
