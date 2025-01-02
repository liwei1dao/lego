package nats

import (
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc"
)

var (
	defsys rpc.IConnPool
)

func OnInit(config map[string]interface{}, opts ...Option) (err error) {
	var option *Options
	if option, err = newOptions(config, opts...); err != nil {
		return
	}
	defsys, err = newSys(option)
	return
}

func NewSys(opts ...Option) (sys rpc.IConnPool, err error) {
	var option *Options
	if option, err = newOptionsByOption(opts...); err != nil {
		return
	}
	sys, err = newSys(option)
	return
}

func Start(host rpc.IBodyHost) (err error) {
	return defsys.Start(host)
}

func Close() {
	defsys.Close()
}

func GetClient(node core.IServiceNode) (client rpc.IConnClient, err error) {
	return defsys.GetClient(node)
}
func AddClient(client rpc.IConnClient, node core.IServiceNode) (err error) {
	return defsys.AddClient(client, node)
}
