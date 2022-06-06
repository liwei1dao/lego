package base

import (
	"context"
	"fmt"
	"sort"
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
	Select(ctx context.Context, stype, sip string, args interface{}) (session IServiceSession, err error) // SelectFunc
	UpdateServer(servers map[string]IServiceSession)
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
func NewSelector(selectMode SelectMode, servers map[string]IServiceSession) ISelector {
	switch selectMode {
	case RandomSelect:
		return newRandomSelector(servers)
	case RoundRobin:
		return newRoundRobinSelector(servers)
	case WeightedRoundRobin:
		return newWeightedRoundRobinSelector(servers)
	case WeightedICMP:
		return newWeightedICMPSelector(servers)
	default:
		return newWeightedRoundRobinSelector(servers)
	}
}

///随机选择
func newRandomSelector(servers map[string]IServiceSession) ISelector {
	ss := make([]IServiceSession, 0, len(servers))
	for _, v := range servers {
		ss = append(ss, v)
	}

	return &randomSelector{servers: ss}
}

type randomSelector struct {
	servers ServiceSessionSlice
}

func (s randomSelector) Select(ctx context.Context, stype, sip string, args interface{}) (session IServiceSession, err error) {
	var ss ServiceSessionSlice
	if ss, err = s.servers.Filter(stype, sip); err != nil {
		return
	}
	i := fastrand.Uint32n(uint32(len(ss)))
	session = ss[i]
	return
}

func (s *randomSelector) UpdateServer(servers map[string]IServiceSession) {
	ss := make(ServiceSessionSlice, 0, len(servers))
	for _, v := range servers {
		ss = append(ss, v)
	}
	s.servers = ss
}

///轮询选择
func newRoundRobinSelector(servers map[string]IServiceSession) ISelector {
	ss := make(ServiceSessionSlice, 0, len(servers))
	for _, v := range servers {
		ss = append(ss, v)
	}

	return &roundRobinSelector{servers: ss}
}

type roundRobinSelector struct {
	servers ServiceSessionSlice
	i       int
}

func (s *roundRobinSelector) Select(ctx context.Context, stype, sip string, args interface{}) (session IServiceSession, err error) {
	var ss ServiceSessionSlice
	if ss, err = s.servers.Filter(stype, sip); err != nil {
		return
	}
	i := s.i
	i = i % len(ss)
	s.i = i + 1
	session = ss[i]
	return
}

func (s *roundRobinSelector) UpdateServer(servers map[string]IServiceSession) {
	ss := make([]IServiceSession, 0, len(servers))
	for _, v := range servers {
		ss = append(ss, v)
	}

	s.servers = ss
}

// 权重选择器
func newWeightedRoundRobinSelector(servers map[string]IServiceSession) ISelector {
	ss := make(ServiceSessionSlice, 0, len(servers))
	for _, v := range servers {
		ss = append(ss, v)
	}

	return &roundRobinSelector{servers: ss}
}

type weightedRoundRobinSelector struct {
	servers ServiceSessionSlice
}

func (s *weightedRoundRobinSelector) Select(ctx context.Context, stype, sip string, args interface{}) (session IServiceSession, err error) {
	var ss ServiceSessionSlice
	if ss, err = s.servers.Filter(stype, sip); err != nil {
		return
	}
	//排序找到最优服务
	sort.Sort(ss)
	session = ss[0]
	return
}

func (s *weightedRoundRobinSelector) UpdateServer(servers map[string]IServiceSession) {
	ss := make([]IServiceSession, 0, len(servers))
	for _, v := range servers {
		ss = append(ss, v)
	}
	s.servers = ss
}

// 网络质量选择器
func newWeightedICMPSelector(servers map[string]IServiceSession) ISelector {
	ss := make(ServiceSessionSlice, 0, len(servers))
	for _, v := range servers {
		rtt, _ := Ping(v.GetIp())
		rtt = CalculateWeight(rtt)
		v.SetPreWeight(float64(rtt))
		ss = append(ss, v)
	}
	return &weightedICMPSelector{servers: ss}
}

// 权重选择器
type weightedICMPSelector struct {
	servers ServiceSessionSlice
}

func (s *weightedICMPSelector) Select(ctx context.Context, stype, sip string, args interface{}) (session IServiceSession, err error) {
	var ss ServiceSessionSlice
	if ss, err = s.servers.Filter(stype, sip); err != nil {
		return
	}
	//排序找到最优服务
	sort.Sort(ss)
	session = ss[0]
	return
}

func (s *weightedICMPSelector) UpdateServer(servers map[string]IServiceSession) {
	ss := make([]IServiceSession, 0, len(servers))
	for _, v := range servers {
		rtt, _ := Ping(v.GetIp())
		rtt = CalculateWeight(rtt)
		v.SetPreWeight(float64(rtt))
		ss = append(ss, v)
	}
	s.servers = ss
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
