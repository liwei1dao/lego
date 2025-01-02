package discovery

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc"
)

type (
	IDiscoveryServicePlugin interface {
		Start() error
		RegisterFunction(serviceName, fname string, fn interface{}, meta string) error
		Stop() error
	}
	IServiceDiscovery interface {
		GetServices() []*rpc.KV
		WatchService() chan []*rpc.KV
	}
	IDiscovery interface {
		Start(node core.IServiceNode) error
		RegisterFunction(serviceName, fname string, fn interface{}, meta string) error
		Stop() error
		GetServices() []*rpc.KV
		WatchService() chan []*rpc.KV
	}
)
