package rabbitmq

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type Option func(*Options)
type Options struct {
	RabbitmqUrl string //Rabbitmq连接地址 amqp://guest:guest@localhost:5672/
	ChannelName string //管道名称
	Durable     bool   //是否持久化
	AutoDelete  bool   //是否自动删除
	Wait        bool   //是都等待
}

func SetRabbitmqUrl(v string) Option {
	return func(o *Options) {
		o.RabbitmqUrl = v
	}
}
func SetChannelName(v string) Option {
	return func(o *Options) {
		o.ChannelName = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		RabbitmqUrl: "amqp://guest:guest@localhost:5672/",
	}
	if config != nil {
		mapstructure.Decode(config, &options)
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}

func newOptionsByOption(opts ...Option) Options {
	options := Options{
		RabbitmqUrl: "amqp://guest:guest@localhost:5672/",
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
