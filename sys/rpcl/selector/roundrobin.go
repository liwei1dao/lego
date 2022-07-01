package selector

import (
	"context"

	"github.com/liwei1dao/lego/sys/rpcl/core"
)

//轮询选择器
type roundRobinSelector struct {
	servers []string
	i       int
}

func newRoundRobinSelector(servers map[string]string) core.ISelector {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}

	return &roundRobinSelector{servers: ss}
}

func (s *roundRobinSelector) Select(ctx context.Context, route core.IRoute, serviceMethod string) string {
	ss := s.servers
	if len(ss) == 0 {
		return ""
	}
	i := s.i
	i = i % len(ss)
	s.i = i + 1

	return ss[i]
}

func (s *roundRobinSelector) UpdateServer(servers map[string]string) {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}

	s.servers = ss
}
