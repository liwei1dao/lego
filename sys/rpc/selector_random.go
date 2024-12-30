package rpc

import (
	"context"

	"github.com/valyala/fastrand"
)

// 随机选择器
type randomSelector struct {
	servers []string
}

func newRandomSelector(servers map[string]string) ISelector {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}

	return &randomSelector{servers: ss}
}

func (s randomSelector) Select(ctx context.Context, service string, args interface{}) string {
	ss := s.servers
	if len(ss) == 0 {
		return ""
	}
	i := fastrand.Uint32n(uint32(len(ss)))
	return ss[i]
}

func (s *randomSelector) UpdateServer(servers map[string]string) {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}

	s.servers = ss
}
