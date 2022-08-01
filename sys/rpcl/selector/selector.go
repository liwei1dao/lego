package selector

import (
	"context"
	"regexp"
	"strings"
	"sync"

	"github.com/liwei1dao/lego/core"
	lcore "github.com/liwei1dao/lego/sys/rpcl/core"
)

var rex_nogather = regexp.MustCompile(`\!\[([^)]+)\]`)
var rex_noid = regexp.MustCompile(`\!([^)]+)`)
var rex_gather = regexp.MustCompile(`\[([^)]+)\]`)

func NewSelector(ervers []*core.ServiceNode) (selector lcore.ISelector, err error) {
	if ervers == nil {
		ervers = make([]*core.ServiceNode, 0)
	}
	selector = &Selector{
		servers: ervers,
	}
	return
}

type Selector struct {
	mutex   sync.RWMutex
	servers []*core.ServiceNode
}

///servicePath = (worker)/(worker/worker_1)/(worker/!worker_1)/(worker/[worker_1,worker_2])/(worker/![worker_1,worker_2])
func (this *Selector) Select(ctx context.Context, servicePath string) (result []*core.ServiceNode) {
	result = make([]*core.ServiceNode, 0)
	service := strings.Split(servicePath, "/")
	leng := len(service)
	this.mutex.RLock()
	if leng == 1 {
		for _, v := range this.servers {
			if v.Type == service[0] {
				result = append(result, v)
			}
		}
	} else if leng == 2 {
		result = this.ParseRoutRules(service[1])
	}
	this.mutex.RUnlock()
	return
}

func (this *Selector) UpdateServer(servers []*core.ServiceNode) {
	if servers == nil {
		return
	}
	this.mutex.Lock()
	this.servers = servers
	this.mutex.Unlock()
}

//路由规则解析
func (this *Selector) ParseRoutRules(rules string) (result []*core.ServiceNode) {
	result = make([]*core.ServiceNode, 0)

	//解析 ![sid,sid] 格式规则
	if out := rex_nogather.FindAllStringSubmatch(rules, -1); len(out) == 1 && len(out[0]) == 2 {
		if nogather := strings.Split(out[0][1], ","); len(nogather) > 0 {
			for _, n := range this.servers {
				iskeep := false
				for _, v := range nogather {
					if n.Id == v {
						iskeep = true
						break
					}
				}
				if !iskeep {
					result = append(result, n)
				}
			}
			return
		}
	}
	//解析 !sid 格式规则
	if out := rex_noid.FindAllStringSubmatch(rules, -1); len(out) == 1 && len(out[0]) == 2 {
		for _, n := range this.servers {
			iskeep := false
			if n.Id == out[0][1] {
				iskeep = true
				break
			}
			if !iskeep {
				result = append(result, n)
			}
		}
		return
	}
	//解析 [sid,sid] 格式规则
	if out := rex_gather.FindAllStringSubmatch(rules, -1); len(out) == 1 && len(out[0]) == 2 {
		if nogather := strings.Split(out[0][1], ","); len(nogather) > 0 {
			for _, n := range this.servers {
				iskeep := false
				for _, v := range nogather {
					if n.Id == v {
						iskeep = true
						break
					}
				}
				if iskeep {
					result = append(result, n)
				}
			}
			return
		}
	}
	for _, n := range this.servers {
		if n.Id == rules {
			result = append(result, n)
		}
	}
	return
}
