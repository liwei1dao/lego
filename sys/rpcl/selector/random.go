package selector

import (
	"context"

	"github.com/liwei1dao/lego/sys/rpcl/core"
	"github.com/valyala/fastrand"
)

// randomSelector selects randomly.
type randomSelector struct {
	servers []string
}

func newRandomSelector(servers map[string]string) core.ISelector {
	ss := make([]string, 0, len(servers))
	for k := range servers {
		ss = append(ss, k)
	}

	return &randomSelector{servers: ss}
}

func (s randomSelector) Select(ctx context.Context, route core.IRoute, serviceMethod string) string {
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
