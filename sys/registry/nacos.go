package registry

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strconv"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc"
	hash "github.com/mitchellh/hashstructure"
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

func newNacos(options Options) (sys *Nacos_Registry, err error) {
	ctx, cancel := context.WithCancel(context.TODO())
	sys = &Nacos_Registry{
		options:   options,
		ctx:       ctx,
		cancel:    cancel,
		shash:     make(map[string]uint64),
		services:  make(map[string]*ServiceNode),
		subscribe: make(map[string]*vo.SubscribeParam),
		rpcsubs:   make(map[core.Rpc_Key][]*ServiceNode),
	}
	// 创建clientConfig
	clientConfig := constant.ClientConfig{
		NamespaceId:         options.Nacos_NamespaceId,
		TimeoutMs:           options.Nacos_TimeoutMs,
		BeatInterval:        options.Nacos_BeatInterval,
		NotLoadCacheAtStart: true,
		LogDir:              "nacos/log",
		CacheDir:            "nacos/cache",
		RotateTime:          "1h",
		MaxAge:              3,
		LogLevel:            "error",
	}
	// 至少一个ServerConfig
	serverConfigs := []constant.ServerConfig{
		{
			IpAddr:      options.Nacos_NacosAddr,
			ContextPath: "/nacos",
			Port:        options.Nacos_Port,
			Scheme:      "http",
		},
	}

	// 创建服务发现客户端的另一种方式 (推荐)
	sys.client, err = clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig:  &clientConfig,
			ServerConfigs: serverConfigs,
		},
	)
	return
}

type Nacos_Registry struct {
	options   Options
	client    naming_client.INamingClient
	isstart   bool
	ctx       context.Context
	cancel    context.CancelFunc
	slock     sync.RWMutex
	shash     map[string]uint64
	subscribe map[string]*vo.SubscribeParam
	services  map[string]*ServiceNode
	rlock     sync.RWMutex
	rpcsubs   map[core.Rpc_Key][]*ServiceNode
}

func (this *Nacos_Registry) Start() (err error) {
	if err = this.getServices(); err != nil && err.Error() != "instance list is empty!" {
		return
	}
	if err = this.registerSNode(&ServiceNode{
		Tag:          this.options.Service.GetTag(),
		Id:           this.options.Service.GetId(),
		IP:           this.options.Service.GetIp(),
		Type:         this.options.Service.GetType(),
		Category:     this.options.Service.GetCategory(),
		Version:      this.options.Service.GetVersion(),
		RpcId:        this.options.Service.GetRpcId(),
		PreWeight:    this.options.Service.GetPreWeight(),
		RpcSubscribe: this.getRpcInfo(),
	}); err != nil {
		return
	}
	ctx1, _ := context.WithCancel(this.ctx)
	go this.run(ctx1)
	ctx2, _ := context.WithCancel(this.ctx)
	go this.listener(ctx2)
	return
}

func (this *Nacos_Registry) Stop() (err error) {
	this.cancel()
	this.isstart = false
	return
}

func (this *Nacos_Registry) PushServiceInfo() (err error) {
	if this.isstart {
		err = this.registerSNode(&ServiceNode{
			Tag:          this.options.Service.GetTag(),
			Id:           this.options.Service.GetId(),
			IP:           this.options.Service.GetIp(),
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
func (this *Nacos_Registry) GetServiceById(sId string) (n *ServiceNode, err error) {
	this.slock.RLock()
	n, ok := this.services[sId]
	this.slock.RUnlock()
	if !ok {
		return nil, fmt.Errorf("No Found %s", sId)
	}
	return
}
func (this *Nacos_Registry) GetServiceByType(sType string) (n []*ServiceNode) {
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
func (this *Nacos_Registry) GetAllServices() (n []*ServiceNode) {
	this.slock.RLock()
	defer this.slock.RUnlock()
	n = make([]*ServiceNode, 0)
	for _, v := range this.services {
		n = append(n, v)
	}
	return
}
func (this *Nacos_Registry) GetServiceByCategory(category core.S_Category) (n []*ServiceNode) {
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
func (this *Nacos_Registry) GetRpcSubById(rId core.Rpc_Key) (n []*ServiceNode) {
	n = make([]*ServiceNode, 0)
	this.rlock.RLock()
	d, ok := this.rpcsubs[rId]
	this.rlock.RUnlock()
	if ok {
		n = d
	}
	return
}

//内部函数-------------------------------------------------------------------------------------------------------------------------
func (this *Nacos_Registry) run(ctx context.Context) {
	t := time.NewTicker(time.Duration(this.options.Consul_RegisterInterval) * time.Second)
	defer t.Stop()
locp:
	for {
		select {
		case <-t.C:
			err := this.PushServiceInfo()
			if err != nil {
				log.Warnf("service run Server.Register error: ", err)
			}
		case <-ctx.Done():
			break locp
		}
	}
	log.Infof("Sys Registry is succ exit run")
}

func (this *Nacos_Registry) listener(ctx context.Context) {
	var (
		slist     model.ServiceList
		instances []model.Instance
		err       error
	)
	t := time.NewTicker(time.Duration(this.options.Nacos_RegisterTTL) * time.Second)
	defer t.Stop()
locp:
	for {
		select {
		case <-t.C:
			if slist, err = this.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
				NameSpace: this.options.Nacos_NamespaceId,
				GroupName: this.options.Service.GetTag(),
				PageNo:    1,
				PageSize:  20,
			}); err == nil {
				var service = make(map[string]struct{})
				for _, v := range slist.Doms {
					service[v] = struct{}{}
					this.rlock.RLock()
					_, ok := this.services[v]
					this.rlock.RUnlock()
					if !ok {
						if instances, err = this.client.SelectInstances(vo.SelectInstancesParam{
							ServiceName: v,
							GroupName:   this.options.Service.GetTag(),
							HealthyOnly: true,
						}); err == nil {
							for _, v1 := range instances {
								if v1.Enable && v1.Healthy {
									this.addandupdataServiceNode(v1)
								} else {
									this.removeServiceNode(v)
								}
							}
						}
					}
				}
				temp := make(map[string]*ServiceNode)
				this.rlock.RLock()
				for k, v := range this.services {
					temp[k] = v
				}
				this.rlock.RUnlock()
				for k, v := range temp {
					if _, ok := service[k]; !ok { //不存在了 移除
						this.removeServiceNode(v.Id)
					}
				}
			}
		case <-ctx.Done():
			break locp
		}
	}
	log.Infof("Sys Registry is succ exit listener")
}

func (this *Nacos_Registry) getServices() (err error) {
	var (
		slist     model.ServiceList
		instances []model.Instance
	)
	if slist, err = this.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		NameSpace: this.options.Nacos_NamespaceId,
		GroupName: this.options.Service.GetTag(),
		PageNo:    1,
		PageSize:  20,
	}); err == nil {
		for _, v := range slist.Doms {
			this.rlock.RLock()
			_, ok := this.services[v]
			this.rlock.RUnlock()
			if !ok {
				if instances, err = this.client.SelectInstances(vo.SelectInstancesParam{
					ServiceName: v,
					GroupName:   this.options.Service.GetTag(),
					HealthyOnly: true,
				}); err == nil {
					for _, v1 := range instances {
						if v1.Enable && v1.Healthy {
							this.addandupdataServiceNode(v1)
						}
					}
				}
			}
		}
	}
	return err
}
func (this *Nacos_Registry) registerSNode(snode *ServiceNode) (err error) {
	h, err := hash.Hash(snode, nil)
	if err != nil {
		return err
	}
	if v, ok := this.shash[snode.Id]; ok && v != 0 && v == h { //没变化不用注册
		return
	}
	var (
	// success bool
	)
	rpcsubscribe, _ := json.Marshal(snode.RpcSubscribe)
	_, err = this.client.RegisterInstance(vo.RegisterInstanceParam{
		Ip:          snode.IP,
		Port:        8848,
		Weight:      snode.PreWeight,
		GroupName:   snode.Tag,
		ServiceName: snode.Id,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata: map[string]string{
			"id":           snode.Id,
			"tag":          snode.Tag,
			"type":         string(snode.Type),
			"category":     string(snode.Category),
			"version":      fmt.Sprintf("%v", snode.Version),
			"rpcid":        snode.RpcId,
			"rpcsubscribe": string(rpcsubscribe),
		},
	})
	return
}
func (this *Nacos_Registry) deregisterSNode() (err error) {
	_, err = this.client.DeregisterInstance(vo.DeregisterInstanceParam{
		Ip:          this.options.Service.GetIp(),
		Port:        8848,
		ServiceName: this.options.Service.GetId(),
		GroupName:   this.options.Service.GetTag(),
	})
	return
}
func (this *Nacos_Registry) addandupdataServiceNode(as model.Instance) (sn *ServiceNode, err error) {
	var (
		version float64
		rpcid   string
		snode   *ServiceNode
		h       uint64
	)
	version, err = strconv.ParseFloat(as.Metadata["version"], 32)
	if err != nil {
		log.Errorf("registry 读取服务节点异常:%s err:%v", as.Metadata["version"], err)
		return
	}
	rpcid = as.Metadata["rpcid"]
	snode = &ServiceNode{
		Tag:          as.Metadata["tag"],
		Type:         as.Metadata["type"],
		Category:     core.S_Category(as.Metadata["category"]),
		Id:           as.Metadata["id"],
		Version:      float32(version),
		RpcId:        rpcid,
		PreWeight:    as.Weight,
		RpcSubscribe: make([]core.Rpc_Key, 0),
	}
	json.Unmarshal([]byte(as.Metadata["rpcsubscribe"]), &snode.RpcSubscribe)
	if h, err = hash.Hash(snode, nil); err == nil {
		_, ok := this.services[snode.Id]
		if !ok {
			log.Infof("发现新的服务【%s】", snode.Id)
			sub := &vo.SubscribeParam{
				ServiceName:       snode.Id,
				GroupName:         snode.Tag,
				SubscribeCallback: this.subscribecallback,
			}
			this.rlock.Lock()
			this.services[snode.Id] = snode
			this.shash[snode.Id] = h
			this.subscribe[snode.Id] = sub
			this.rlock.Unlock()
			if this.options.Listener != nil { //异步通知
				go this.options.Listener.FindServiceHandlefunc(*snode)
			}
			err = this.client.Subscribe(sub)
		} else {
			this.rlock.RLock()
			v, ok := this.shash[snode.Id]
			this.rlock.RUnlock()
			if !ok || v != h { //校验不一致
				// log.Debugf("更新服务【%s】", snode.Id)
				this.rlock.Lock()
				this.services[snode.Id] = snode
				this.shash[snode.Id] = h
				this.rlock.Unlock()
				if this.options.Listener != nil { //异步通知
					go this.options.Listener.UpDataServiceHandlefunc(*snode)
				}
			}
		}
	} else {
		log.Errorf("registry 校验服务 hash 值异常:%s", err.Error())
	}
	this.syncRpcInfo()
	return
}
func (this *Nacos_Registry) removeServiceNode(sId string) {
	this.rlock.RLock()
	_, ok := this.services[sId]
	sub, _ := this.subscribe[sId]
	this.rlock.RUnlock()
	if !ok {
		return
	}
	log.Infof("丢失服务【%s】", sId)
	this.rlock.Lock()
	delete(this.services, sId)
	delete(this.subscribe, sId)
	delete(this.shash, sId)
	this.rlock.Unlock()
	this.syncRpcInfo()
	this.client.Unsubscribe(sub)
	if this.options.Listener != nil { //异步通知
		go this.options.Listener.LoseServiceHandlefunc(sId)
	}
}
func (this *Nacos_Registry) syncRpcInfo() {
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
func (this *Nacos_Registry) getRpcInfo() (rfs []core.Rpc_Key) {
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

func (this *Nacos_Registry) subscribecallback(services []model.SubscribeService, err error) {
	for _, v := range services {
		log.Debugf("subscribecallback:%+v", v)
		this.addandupdataServiceNode(model.Instance{
			Valid:       v.Valid,
			Enable:      v.Enable,
			ServiceName: v.ServiceName,
			Ip:          v.Ip,
			Port:        v.Port,
			Metadata:    v.Metadata,
			Weight:      v.Weight,
		})
	}
}
