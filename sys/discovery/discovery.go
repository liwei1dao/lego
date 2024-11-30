package discovery

import (
	"context"
	"net"

	"github.com/liwei1dao/lego/sys/discovery/dcore"
)

func newSys(options Options) (sys *Discovery, err error) {
	sys = &Discovery{options: options}
	return
}

// 服务发现
type Discovery struct {
	options   Options
	sPlugin   dcore.ServicePlugin
	newFinder func(basePath string, servicePath string) (IDiscovery, error)
}

func (this *Discovery) Start() error {
	return this.sPlugin.Start()
}

func (this *Discovery) Stop() error {
	return this.sPlugin.Stop()
}

func (this *Discovery) Register(name string, rcvr interface{}, metadata string) error {
	return this.sPlugin.Register(name, rcvr, metadata)
}

func (this *Discovery) Unregister(name string) error {
	return this.sPlugin.Unregister(name)
}
func (this *Discovery) RegisterFunction(serviceName, fname string, fn interface{}, metadata string) error {
	return this.sPlugin.RegisterFunction(serviceName, fname, fn, metadata)
}
func (this *Discovery) HandleConnAccept(conn net.Conn) (net.Conn, bool) {
	return this.sPlugin.HandleConnAccept(conn)
}
func (this *Discovery) PreCall(ctx context.Context, serviceName, methodName string, args interface{}) (interface{}, error) {
	return this.sPlugin.PreCall(ctx, serviceName, methodName, args)
}

func (this *Discovery) NewFinder(basePath string, servicePath string) (IDiscovery, error) {
	return this.newFinder(basePath, servicePath)
}
