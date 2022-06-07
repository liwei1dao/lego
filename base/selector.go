package base

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/go-ping/ping"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/utils/container/version"
	"github.com/valyala/fastrand"
)

//服务选择器类型
type SelectMode int

const (
	//随机选择
	RandomSelect SelectMode = iota
	//轮询选择
	RoundRobin
	//权重选择
	WeightedRoundRobin
	//网络延迟选择
	WeightedICMP
	//地理位置选择
	ConsistentHash
)

type ISelector interface {
	Load(sId string) (session IServiceSession, ok bool)
	Select(ctx context.Context, stype, sip string) (session IServiceSession, err error) // SelectFunc
	IsHave(sId string) bool
	AddServer(server IServiceSession)
	RemoveServer(sId string)
	UpdateServer(node core.ServiceNode)
}

//服务会话切片
type ServiceSessionSlice []IServiceSession

func (this ServiceSessionSlice) Len() int { return len(this) }
func (this ServiceSessionSlice) Less(i, j int) bool {
	if iscompare := version.CompareStrVer(this[i].GetVersion(), this[j].GetVersion()); iscompare != 0 {
		return iscompare > 0
	} else {
		if this[i].GetPreWeight() < this[j].GetPreWeight() {
			return true
		} else if this[i].GetPreWeight() > this[j].GetPreWeight() {
			return false
		} else {
			return true
		}
	}
}
func (this ServiceSessionSlice) Swap(i, j int) { this[i], this[j] = this[j], this[i] }

func (this ServiceSessionSlice) Filter(stype, sip string) (ss ServiceSessionSlice, err error) {
	ss = make([]IServiceSession, 0, len(this))
	for _, v := range this {
		if (sip == core.AutoIp || sip == v.GetIp()) && v.GetType() == stype {
			ss = append(ss, v)
		}
	}
	if len(ss) == 0 {
		err = fmt.Errorf("on found services type[%s]ip [%s] ", stype, sip)
	}
	return
}

//创建选择器
func NewSelector(selectMode SelectMode) ISelector {
	switch selectMode {
	case RandomSelect:
		return newRandomSelector()
	case RoundRobin:
		return newRoundRobinSelector()
	case WeightedRoundRobin:
		return newWeightedRoundRobinSelector()
	case WeightedICMP:
		return newWeightedICMPSelector()
	default:
		return newWeightedRoundRobinSelector()
	}
}

type baseSelector struct {
	servers ServiceSessionSlice
	lock    sync.RWMutex
}

func (this *baseSelector) Load(sId string) (session IServiceSession, ok bool) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, v := range this.servers {
		if v.GetId() == sId {
			return v, true
		}
	}
	return nil, false
}

func (this *baseSelector) IsHave(sId string) bool {
	this.lock.RLock()
	defer this.lock.RUnlock()
	for _, v := range this.servers {
		if v.GetId() == sId {
			return true
		}
	}
	return false
}

func (this *baseSelector) AddServer(server IServiceSession) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, v := range this.servers {
		if v.GetId() == server.GetId() {
			server.Done()
			return
		}
	}
	this.servers = append(this.servers, server)
}
func (this *baseSelector) RemoveServer(sId string) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for i, v := range this.servers {
		if v.GetId() == sId {
			this.servers = append(this.servers[0:i], this.servers[i+1:]...)
			v.Done()
			return
		}
	}
}
func (this *baseSelector) UpdateServer(node core.ServiceNode) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, v := range this.servers {
		if v.GetId() == node.Id {
			v.SetPreWeight(node.PreWeight)
			v.SetVersion(node.Version)
		}
	}
}

///随机选择
func newRandomSelector() ISelector {
	return &randomSelector{}
}

type randomSelector struct {
	baseSelector
}

func (this *randomSelector) Select(ctx context.Context, stype, sip string) (session IServiceSession, err error) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	var ss ServiceSessionSlice
	if ss, err = this.servers.Filter(stype, sip); err != nil {
		return
	}
	i := fastrand.Uint32n(uint32(len(ss)))
	session = ss[i]
	return
}

///轮询选择
func newRoundRobinSelector() ISelector {
	return &roundRobinSelector{}
}

type roundRobinSelector struct {
	baseSelector
	i int
}

func (this *roundRobinSelector) Select(ctx context.Context, stype, sip string) (session IServiceSession, err error) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	var ss ServiceSessionSlice
	if ss, err = this.servers.Filter(stype, sip); err != nil {
		return
	}
	i := this.i
	i = i % len(ss)
	this.i = i + 1
	session = ss[i]
	return
}

// 权重选择器
func newWeightedRoundRobinSelector() ISelector {
	return &roundRobinSelector{}
}

type weightedRoundRobinSelector struct {
	baseSelector
}

func (this *weightedRoundRobinSelector) Select(ctx context.Context, stype, sip string, args interface{}) (session IServiceSession, err error) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	var ss ServiceSessionSlice
	if ss, err = this.servers.Filter(stype, sip); err != nil {
		return
	}
	//排序找到最优服务
	sort.Sort(ss)
	session = ss[0]
	return
}

// 网络质量选择器
func newWeightedICMPSelector() ISelector {
	return &weightedICMPSelector{}
}

// 权重选择器
type weightedICMPSelector struct {
	baseSelector
}

func (this *weightedICMPSelector) Select(ctx context.Context, stype, sip string) (session IServiceSession, err error) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	var ss ServiceSessionSlice
	if ss, err = this.servers.Filter(stype, sip); err != nil {
		return
	}
	//排序找到最优服务
	sort.Sort(ss)
	session = ss[0]
	return
}

func (this *weightedICMPSelector) AddServer(server IServiceSession) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, v := range this.servers {
		if v.GetId() == server.GetId() {
			server.Done()
			return
		}
	}
	rtt, _ := Ping(server.GetIp())
	rtt = CalculateWeight(rtt)
	server.SetPreWeight(float64(rtt))
	this.servers = append(this.servers, server)
}

func (this *weightedICMPSelector) UpdateServer(node core.ServiceNode) {
	this.lock.Lock()
	defer this.lock.Unlock()
	for _, v := range this.servers {
		if v.GetId() == node.Id {
			rtt, _ := Ping(node.IP)
			rtt = CalculateWeight(rtt)
			v.SetPreWeight(float64(rtt))
			v.SetVersion(node.Version)
		}
	}
}

// 获取目标主机的网络延迟数
func Ping(host string) (rtt int, err error) {
	rtt = 1000 //default and timeout is 1000 ms

	pinger, err := ping.NewPinger(host)
	if err != nil {
		return rtt, err
	}
	pinger.Count = 3
	pinger.Timeout = 3 * time.Second
	err = pinger.Run()
	if err != nil {
		return rtt, err
	}
	stats := pinger.Statistics()
	// ping failed
	if len(stats.Rtts) == 0 {
		return rtt, err
	}
	rtt = int(stats.AvgRtt) / 1e6

	return rtt, err
}

// 将延迟数转换成权重值
func CalculateWeight(rtt int) int {
	switch {
	case rtt >= 0 && rtt <= 10:
		return 191
	case rtt > 10 && rtt <= 200:
		return 201 - rtt
	case rtt > 100 && rtt < 1000:
		return 1
	default:
		return 0
	}
}
