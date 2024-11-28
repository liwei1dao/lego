package discovery

import (
	"context"
	"net"
)

type KVPair struct {
	Key   string
	Value string
}

type DiscoveryFilter func(kvp *KVPair) bool
type IDiscovery interface {
	GetServices() []*KVPair
	WatchService() chan []*KVPair
	RemoveWatcher(ch chan []*KVPair)
	Clone(servicePath string) (IDiscovery, error)
	SetFilter(DiscoveryFilter)
	Close()
}

// 服务发现系统
type ISys interface {
	Start() error
	Stop() error
	Register(name string, rcvr interface{}, metadata string) error
	Unregister(name string) error
	RegisterFunction(serviceName, fname string, fn interface{}, metadata string) error
	HandleConnAccept(net.Conn) (net.Conn, bool)
	PreCall(ctx context.Context, serviceName, methodName string, args interface{}) (interface{}, error)
	NewFinder(basePath string, servicePath string) (IDiscovery, error)
}

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func Start() error {
	return defsys.Start()
}
func Stop() error {
	return defsys.Stop()
}
func Register(name string, rcvr interface{}, metadata string) error {
	return defsys.Register(name, rcvr, metadata)
}
func Unregister(name string) error {
	return defsys.Unregister(name)
}
func RegisterFunction(serviceName, fname string, fn interface{}, metadata string) error {
	return defsys.RegisterFunction(serviceName, fname, fn, metadata)
}
func HandleConnAccept(conn net.Conn) (net.Conn, bool) {
	return defsys.HandleConnAccept(conn)
}
func PreCall(ctx context.Context, serviceName, methodName string, args interface{}) (interface{}, error) {
	return defsys.PreCall(ctx, serviceName, methodName, args)
}
func NewFinder(basePath string, servicePath string) (IDiscovery, error) {
	return defsys.NewFinder(basePath, servicePath)
}
