package rpcx

import (
	"context"
	"regexp"
	"strings"
	"sync"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/smallnest/rpcx/share"
	"github.com/valyala/fastrand"
)

var rex_nogather = regexp.MustCompile(`\!\[([^)]+)\]`)
var rex_noid = regexp.MustCompile(`\!([^)]+)`)
var rex_gather = regexp.MustCompile(`\[([^)]+)\]`)

func newSelector(log log.ILogger, stag string, fn func(map[string]*ServiceNode)) *Selector {
	return &Selector{
		log:               log,
		stag:              stag,
		updateServerEvent: fn,
		servers:           make(map[string]*ServiceNode),
		serversType:       make(map[string][]*ServiceNode),
		i:                 make(map[string]int),
	}
}

type ServiceNode struct {
	ServiceTag  string `json:"stag"`    //服务集群标签
	ServiceId   string `json:"sid"`     //服务id
	ServiceType string `json:"stype"`   //服务类型
	Version     string `json:"version"` //服务版本
	ServiceAddr string `json:"addr"`    //服务地址
}

type Selector struct {
	log               log.ILogger
	stag              string
	updateServerEvent func(map[string]*ServiceNode)
	servers           map[string]*ServiceNode
	serversType       map[string][]*ServiceNode
	lock              sync.RWMutex
	i                 map[string]int
}

// /servicePath = (worker)/(worker/worker_1)/(worker/!worker_1)/(worker/[worker_1,worker_2])/(worker/![worker_1,worker_2])
func (this *Selector) Select(ctx context.Context, servicePath, serviceMethod string, args interface{}) string {
	routrules := ctx.Value(share.ReqMetaDataKey).(map[string]string)[CallRoutRulesKey]
	service := strings.Split(routrules, "/")
	leng := len(service)
	if leng == 1 {
		if nodes, ok := this.serversType[service[0]]; ok {
			this.lock.RLock()
			i, ok := this.i[service[0]]
			this.lock.RUnlock()
			if !ok {
				i = 0
			}
			i = i % len(nodes)
			this.lock.Lock()
			this.i[service[0]] = i + 1
			this.lock.Unlock()
			return nodes[i].ServiceAddr
		}
	} else if leng == 2 {
		result := this.ParseRoutRules(service[1])
		if len(result) == 0 {
			this.log.Error("Select no found any node",
				log.Field{Key: "stag", Value: this.stag},
				log.Field{Key: "servicePath", Value: servicePath},
				log.Field{Key: "serviceMethod", Value: serviceMethod},
				log.Field{Key: "routrules", Value: routrules},
			)
			return ""
		}
		i := fastrand.Uint32n(uint32(len(result)))
		if node, ok := this.servers[result[i]]; ok {
			return node.ServiceAddr
		}
	}
	// this.log.Error("Select no found any node", log.Field{"stag", this.stag}, log.Field{"servicePath", servicePath}, log.Field{"serviceMethod", serviceMethod}, log.Field{"routrules", routrules})
	return ""
}

// 找到同类型节点信息
func (this *Selector) Find(ctx context.Context, servicePath, serviceMethod string, args interface{}) []string {
	if nodes, ok := this.serversType[servicePath]; ok {
		addrs := make([]string, len(nodes))
		for i, v := range nodes {
			addrs[i] = v.ServiceAddr
		}
		return addrs
	}
	return nil
}

// 更新服务列表
func (this *Selector) UpdateServer(servers map[string]string) {
	ss := make(map[string]*ServiceNode)
	sst := make(map[string][]*ServiceNode)
	for _, v := range servers {
		if node, err := smetaToServiceNode(v); err != nil {
			this.log.Errorf("smetaToServiceNode:%s err:%v", v, err)
			continue
		} else {
			ss[node.ServiceId] = node
			if _, ok := sst[node.ServiceType]; !ok {
				sst[node.ServiceType] = make([]*ServiceNode, 0)
				sst[node.ServiceType] = append(sst[node.ServiceType], node)
			} else {
				sst[node.ServiceType] = append(sst[node.ServiceType], node)
			}
		}

	}
	this.servers = ss
	this.serversType = sst
	if this.updateServerEvent != nil {
		go this.updateServerEvent(ss)
	}
}

// 路由规则解析
func (this *Selector) ParseRoutRules(rules string) (result []string) {
	result = make([]string, 0)

	//解析 ![sid,sid] 格式规则
	if out := rex_nogather.FindAllStringSubmatch(rules, -1); len(out) == 1 && len(out[0]) == 2 {
		if nogather := strings.Split(out[0][1], ","); len(nogather) > 0 {
			for k, _ := range this.servers {
				iskeep := false
				for _, v := range nogather {
					if k == v {
						iskeep = true
						break
					}
				}
				if !iskeep {
					result = append(result, k)
				}
			}
			return
		}
	}
	//解析 !sid 格式规则
	if out := rex_noid.FindAllStringSubmatch(rules, -1); len(out) == 1 && len(out[0]) == 2 {
		for k, _ := range this.servers {
			iskeep := false
			if k == out[0][1] {
				iskeep = true
				break
			}
			if !iskeep {
				result = append(result, k)
			}
		}
		return
	}
	//解析 [sid,sid] 格式规则
	if out := rex_gather.FindAllStringSubmatch(rules, -1); len(out) == 1 && len(out[0]) == 2 {
		if nogather := strings.Split(out[0][1], ","); len(nogather) > 0 {
			for k, _ := range this.servers {
				iskeep := false
				for _, v := range nogather {
					if k == v {
						iskeep = true
						break
					}
				}
				if iskeep {
					result = append(result, k)
				}
			}
			return
		}
	}
	if _, ok := this.servers[rules]; ok {
		result = append(result, rules)
	}
	return
}
