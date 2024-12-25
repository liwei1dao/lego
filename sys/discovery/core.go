package discovery

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/discovery/dcore"
)

type (
	IDiscoveryServicePlugin interface {
		Start() error
		RegisterFunction(serviceName, fname string, fn interface{}, meta string) error
		Stop() error
	}
	IServiceDiscovery interface {
		GetServices() []*dcore.KV
		WatchService() chan []*dcore.KV
	}
	IDiscovery interface {
		Start(node core.IServiceNode) error
		RegisterFunction(serviceName, fname string, fn interface{}, meta string) error
		Stop() error
		GetServices() []*dcore.KV
		WatchService() chan []*dcore.KV
	}
)
