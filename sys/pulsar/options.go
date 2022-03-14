package pulsar

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type PulsarStartType int8 //Pulsar启动类型
const (
	Producer PulsarStartType = iota //生产者
	Consumer                        //消费者
	All                             //同时支持生成消费
)

type Option func(*Options)
type Options struct {
	StartType              PulsarStartType //kafka启动类型
	PulsarUrl              string          //Pulsar连接地址 pulsar://localhost:6550,localhost:6651,localhost:6652
	Topics                 []string        //生产消费Topic 生产默认第一个Topic
	ConsumerGroupId        string          //消费者配置 GroupId
	Producer_Return_Errors bool            //是否接受错误信息
}

func SetStartType(v PulsarStartType) Option {
	return func(o *Options) {
		o.StartType = v
	}
}
func SetPulsarUrl(v string) Option {
	return func(o *Options) {
		o.PulsarUrl = v
	}
}
func SetTopics(v []string) Option {
	return func(o *Options) {
		o.Topics = v
	}
}
func SetConsumerGroupId(v string) Option {
	return func(o *Options) {
		o.ConsumerGroupId = v
	}
}
func SetProducer_Return_Errors(v bool) Option {
	return func(o *Options) {
		o.Producer_Return_Errors = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		PulsarUrl: "pulsar://127.0.0.1:6550",
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
	options := Options{}
	for _, o := range opts {
		o(&options)
	}
	return options
}
