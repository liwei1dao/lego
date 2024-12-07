package discovery

import (
	"errors"

	"github.com/liwei1dao/lego/sys/discovery/consul"
	"github.com/liwei1dao/lego/sys/discovery/etcd2"
	"github.com/liwei1dao/lego/sys/discovery/etcd3"
	"github.com/liwei1dao/lego/sys/discovery/redis"
)

func newSys(options *Options) (d *Discovery, err error) {
	d = &Discovery{
		options: options,
	}
	switch options.StoreType {
	case StoreConsul:
		d.service = &consul.ConsulRegisterPlugin{
			ConsulServers:  options.Endpoints,
			UpdateInterval: options.UpdateInterval,
			Expired:        options.UpdateInterval * 2,
		}
		break
	case StoreZookeeper:

		break
	case StoreEtcd2:
		d.service = &etcd2.EtcdRegisterPlugin{
			EtcdServers:    options.Endpoints,
			UpdateInterval: options.UpdateInterval,
			Expired:        options.UpdateInterval * 2,
		}
		break
	case StoreEtcd3:
		d.service = &etcd3.EtcdV3RegisterPlugin{
			EtcdServers:    options.Endpoints,
			UpdateInterval: options.UpdateInterval,
			Expired:        options.UpdateInterval * 2,
		}
		break
	case StoreRedis:
		d.service = &redis.RedisRegisterPlugin{
			RedisServers:   options.Endpoints,
			UpdateInterval: options.UpdateInterval,
			Options:        options.Config,
		}
		break
	default:
		err = errors.New("StoreType Not Found")
	}
	return
}

type Discovery struct {
	options *Options
	service IServicePlugin
	client  IDiscoveryClient
}
