package lgrpc

import (
	"context"
	"regexp"
)

/*
系统描述:rpc服务通信系统，参考xrpc设计思路
*/

type SelectMode int

const (
	//随机选择
	RandomSelect SelectMode = iota
	//轮询选择
	RoundRobin
	权重轮询选择器
	WeightedRoundRobin
	//权重ICMP选择器
	WeightedICMP
	//一致性哈希
	ConsistentHash
	//关闭选择器
	Closest

	// SelectByUser is selecting by implementation of users
	SelectByUser = 1000
)

type (
	// 选择器
	ISelector interface {
		Select(ctx context.Context, service string, args interface{}) string
		UpdateServer(servers map[string]string)
	}
)

func newSelector(selectMode SelectMode, servers map[string]string) ISelector {
	switch selectMode {
	case RandomSelect:
		return newRandomSelector(servers)
	case RoundRobin:
		return newRoundRobinSelector(servers)
	case WeightedRoundRobin:
		return newWeightedRoundRobinSelector(servers)
	case WeightedICMP:
		return newWeightedICMPSelector(servers)
	case ConsistentHash:
		return newConsistentHashSelector(servers)
	case SelectByUser:
		return nil
	default:
		return newRandomSelector(servers)
	}
}

var rex_nogather = regexp.MustCompile(`\!\[([^)]+)\]`)
var rex_noid = regexp.MustCompile(`\!([^)]+)`)
var rex_gather = regexp.MustCompile(`\[([^)]+)\]`)

// func NewSelector(ervers []core.IServiceNode) (selector rpccore.ISelector, err error) {
// 	if ervers == nil {
// 		ervers = make([]core.IServiceNode, 0)
// 	}
// 	selector = &Selector{
// 		servers: ervers,
// 	}
// 	return
// }

// type Selector struct {
// 	mutex   sync.RWMutex
// 	servers []core.IServiceNode
// }

// // /servicePath = (stype)|(stype/sid)|(stype/!sid)|(stype/[sid1,sid2])|(stype/![sid1,sid2])
// func (this *Selector) Select(ctx context.Context) (result []core.IServiceNode) {
// 	result = make([]core.IServiceNode, 0)
// 	service := strings.Split(servicePath, "/")
// 	leng := len(service)
// 	this.mutex.RLock()
// 	if leng == 1 {
// 		for _, v := range this.servers {
// 			if v.Type() == service[0] {
// 				result = append(result, v)
// 			}
// 		}
// 	} else if leng == 2 {

// 		result = this.ParseRoutRules(service[1])
// 	}
// 	this.mutex.RUnlock()
// 	return
// }

// func (this *Selector) UpdateServer(servers map[string]string) (add, del, change []core.IServiceNode) {
// 	if servers == nil {
// 		log.Error("UpdateServer 传参错误!")
// 		return
// 	}
// 	var (
// 		iskeep bool
// 	)
// 	add = make([]core.IServiceNode, 0)
// 	change = make([]core.IServiceNode, 0)
// 	this.mutex.RLock()
// 	del = make([]core.IServiceNode, len(this.servers))
// 	for i, v := range this.servers {
// 		del[i] = v
// 	}
// 	this.mutex.RUnlock()
// 	for _, v1 := range servers {
// 		iskeep = false
// 		for i, v2 := range del {
// 			if v1.Path() == v2.Path() {
// 				iskeep = true
// 				if !v1.Equal(v2) { //有变化
// 					change = append(change, v1)
// 				}
// 				del = append(del[0:i], del[i+1:]...) //移除存在的节点 过滤出被销毁的节点
// 				break
// 			}
// 		}
// 		if !iskeep {
// 			add = append(add, v1)
// 		}
// 	}

// 	this.mutex.Lock()
// 	this.servers = servers
// 	this.mutex.Unlock()
// 	return
// }

// // 路由规则解析
// func (this *Selector) ParseRoutRules(rules string) (result []core.IServiceNode) {
// 	if rules == "" {
// 		return
// 	}

// 	result = make([]core.IServiceNode, 0)

// 	//解析 ![sid,sid] 格式规则
// 	if out := rex_nogather.FindAllStringSubmatch(rules, -1); len(out) == 1 && len(out[0]) == 2 {
// 		if nogather := strings.Split(out[0][1], ","); len(nogather) > 0 {
// 			for _, n := range this.servers {
// 				iskeep := false
// 				for _, v := range nogather {
// 					if n.Id() == v {
// 						iskeep = true
// 						break
// 					}
// 				}
// 				if !iskeep {
// 					result = append(result, n)
// 				}
// 			}
// 			return
// 		}
// 	}
// 	//解析 !sid 格式规则
// 	if out := rex_noid.FindAllStringSubmatch(rules, -1); len(out) == 1 && len(out[0]) == 2 {
// 		for _, n := range this.servers {
// 			iskeep := false
// 			if n.Id() == out[0][1] {
// 				iskeep = true
// 				break
// 			}
// 			if !iskeep {
// 				result = append(result, n)
// 			}
// 		}
// 		return
// 	}
// 	//解析 [sid,sid] 格式规则
// 	if out := rex_gather.FindAllStringSubmatch(rules, -1); len(out) == 1 && len(out[0]) == 2 {
// 		if nogather := strings.Split(out[0][1], ","); len(nogather) > 0 {
// 			for _, n := range this.servers {
// 				iskeep := false
// 				for _, v := range nogather {
// 					if n.Id() == v {
// 						iskeep = true
// 						break
// 					}
// 				}
// 				if iskeep {
// 					result = append(result, n)
// 				}
// 			}
// 			return
// 		}
// 	}
// 	for _, n := range this.servers {
// 		if n.Id() == rules {
// 			result = append(result, n)
// 		}
// 	}
// 	return
// }
