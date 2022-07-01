package selector

import "github.com/liwei1dao/lego/sys/rpcl/core"

func NewSelector(selectMode core.SelectMode, servers map[string]string) core.ISelector {
	switch selectMode {
	case core.RandomSelect:
		return newRandomSelector(servers)
	case core.RoundRobin:
		return newRoundRobinSelector(servers)
	case core.WeightedRoundRobin:
		return newWeightedRoundRobinSelector(servers)
	case core.WeightedICMP:
		return newWeightedICMPSelector(servers)
	default:
		return newRandomSelector(servers)
	}
}
