package selector

import (
	"context"

	"github.com/liwei1dao/lego/sys/rpcl/core"
)

//权重加轮询选择器
type weightedRoundRobinSelector struct {
}

func newWeightedRoundRobinSelector(servers map[string]string) core.ISelector {
	return &weightedRoundRobinSelector{}
}

func (s *weightedRoundRobinSelector) Select(ctx context.Context, route core.IRoute, serviceMethod string) string {
	return ""
}

func (s *weightedRoundRobinSelector) UpdateServer(servers map[string]string) {

}
