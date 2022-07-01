package kafka

import (
	"context"
	"sync"

	"github.com/liwei1dao/lego/sys/rpcl/core"
)

type RPC struct {
	sys         core.ISys
	service     *Service
	clientMapMu sync.RWMutex
	clients     map[string]*Client
}

func (this *RPC) Go(ctx context.Context, servicePath, serviceMethod string, args interface{}, reply interface{}, done chan *core.Call) (call *core.Call, err error) {
	this.clientMapMu.RLock()
	client, ok := this.clients[servicePath]
	this.clientMapMu.RUnlock()
	if !ok {
		if client, err = newClient(this.sys); err != nil {
			return
		}
	}
	call = client.Go(ctx, serviceMethod, args, reply, done)
	return
}
