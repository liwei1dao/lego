package kafka

import (
	"time"

	"github.com/Shopify/sarama"
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type KafkaStartType int8 //kafka启动类型
const (
	Syncproducer             KafkaStartType = iota ///同步生产者
	Asyncproducer                                  ///异步生产者
	Consumer                                       ///消费者
	SyncproducerAndConsumer                        ///同步生产者和消费
	AsyncproducerAndConsumer                       ///异步生产和消费
	All                                            ///全部启动
)

type Option func(*Options)
type Options struct {
	StartType                 KafkaStartType          //kafka启动类型
	Hosts                     []string                //集群地址
	Topics                    []string                //消费者配置 Topics
	GroupId                   string                  //消费者配置 GroupId
	ClientID                  string                  //用户提供的字符串随每个请求发送到代理进行日志记录，调试和审计目的。默认为“sarama”，但你应该将其设置为特定于您的应用程序的内容。
	Producer_RequiredAcks     sarama.RequiredAcks     //设置生产者 消息 回复等级 0 1 all
	Producer_MaxMessageBytes  int                     //消息的最大允许大小（默认为 1000000）。应该设置等于或小于代理的 `message.max.bytes`
	Producer_Return_Successes bool                    //如果启用 成功发送的消息将在成功通道中
	Producer_Return_Errors    bool                    //如果启用，发送失败的消息将在错误通道中
	Producer_Retry_Max        int                     //重试发送消息的总次数（默认为 3）
	Producer_Retry_Backoff    int                     //在重试之间等待集群稳定的时间（默认为 100 毫秒）
	Producer_Compression      sarama.CompressionCodec //用于消息的压缩类型（默认为无压缩）
	Producer_CompressionLevel int                     //用于消息的压缩级别。意义取决于在实际使用的压缩类型上，默认为默认压缩编解码器的级别。
	Net_DialTimeout           time.Duration           //默认30 秒
	Net_ReadTimeout           time.Duration           //默认30 秒
	Net_WriteTimeout          time.Duration           //默认30 秒
	Net_KeepAlive             time.Duration           //KeepAlive 指定活动网络连接的保持活动期。如果为零，则禁用保活。 （默认为 0：禁用）单位 秒
	Consumer_Offsets_Initial  int64                   //OffsetNewest -1 or OffsetOldest -2 默认 OffsetOldest
	Consumer_Return_Errors    bool                    //如果启用，则在消费时发生的任何错误都将返回错误通道（默认禁用）
}

///kafka启动类型
func SetStartType(v KafkaStartType) Option {
	return func(o *Options) {
		o.StartType = v
	}
}

///设置kafka集群地址
func SetHosts(v []string) Option {
	return func(o *Options) {
		o.Hosts = v
	}
}

///设置kafka消费者的消费Topics
func SetTopics(v []string) Option {
	return func(o *Options) {
		o.Topics = v
	}
}

///设置消费组 GroupId
func SetGroupId(v string) Option {
	return func(o *Options) {
		o.GroupId = v
	}
}

///用户提供的字符串随每个请求发送到代理进行日志记录，调试和审计目的。默认为“sarama”，但你应该将其设置为特定于您的应用程序的内容。
func SetClientID(v string) Option {
	return func(o *Options) {
		o.ClientID = v
	}
}

///设置生产者 消息 回复等级 0 1 all
func SetProducer_RequiredAcks(v sarama.RequiredAcks) Option {
	return func(o *Options) {
		o.Producer_RequiredAcks = v
	}
}

///消息的最大允许大小（默认为 1000000）。应该设置等于或小于代理的 `message.max.bytes`
func SetProducer_MaxMessageBytes(v int) Option {
	return func(o *Options) {
		o.Producer_MaxMessageBytes = v
	}
}

///如果启用 成功发送的消息将在成功通道中
func SetProducer_Return_Successes(v bool) Option {
	return func(o *Options) {
		o.Producer_Return_Successes = v
	}
}

///如果启用，发送失败的消息将在错误通道中
func SetProducer_Return_Errors(v bool) Option {
	return func(o *Options) {
		o.Producer_Return_Errors = v
	}
}

///重试发送消息的总次数（默认为 3）
func SetProducer_Retry_Max(v int) Option {
	return func(o *Options) {
		o.Producer_Retry_Max = v
	}
}

///重试发送消息的总次数（默认为 3）
func SetProducer_Retry_Backoff(v int) Option {
	return func(o *Options) {
		o.Producer_Retry_Backoff = v
	}
}

///用于消息的压缩类型（默认为无压缩）
func SetProducer_Compression(v sarama.CompressionCodec) Option {
	return func(o *Options) {
		o.Producer_Compression = v
	}
}

///用于消息的压缩级别。意义取决于在实际使用的压缩类型上，默认为默认压缩编解码器的级别。
func SetProducer_CompressionLevel(v int) Option {
	return func(o *Options) {
		o.Producer_CompressionLevel = v
	}
}

///默认5 秒
func SetNet_DialTimeout(v time.Duration) Option {
	return func(o *Options) {
		o.Net_DialTimeout = v
	}
}

///默认60 秒
func SetNet_ReadTimeout(v time.Duration) Option {
	return func(o *Options) {
		o.Net_ReadTimeout = v
	}
}

///默认60 秒
func SetNet_WriteTimeout(v time.Duration) Option {
	return func(o *Options) {
		o.Net_WriteTimeout = v
	}
}

///KeepAlive 指定活动网络连接的保持活动期。如果为零，则禁用保活。 （默认为 0：禁用）单位 秒
func SetNet_KeepAlive(v time.Duration) Option {
	return func(o *Options) {
		o.Net_KeepAlive = v
	}
}

///OffsetNewest -1 or OffsetOldest -2 默认 OffsetOldest
func SetConsumer_Offsets_Initial(v int64) Option {
	return func(o *Options) {
		o.Consumer_Offsets_Initial = v
	}
}

///如果启用，则在消费时发生的任何错误都将返回错误通道（默认禁用）
func SetConsumer_Return_Errors(v bool) Option {
	return func(o *Options) {
		o.Consumer_Return_Errors = v
	}
}

func newOptions(config map[string]interface{}, opts ...Option) Options {
	options := Options{
		Producer_Return_Successes: true,
		Producer_Return_Errors:    true,
		Producer_RequiredAcks:     sarama.NoResponse,
		Producer_MaxMessageBytes:  1000000,
		Producer_Retry_Max:        3,
		Producer_Retry_Backoff:    100,
		Producer_Compression:      sarama.CompressionNone,
		Net_DialTimeout:           time.Second * 5,
		Net_ReadTimeout:           time.Second * 60,
		Net_WriteTimeout:          time.Second * 60,
		Net_KeepAlive:             0,
		Consumer_Offsets_Initial:  -1,
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
		Producer_RequiredAcks:    sarama.NoResponse,
		Producer_MaxMessageBytes: 1000000,
		Producer_Retry_Max:       3,
		Producer_Retry_Backoff:   100,
		Producer_Compression:     sarama.CompressionNone,
		Net_DialTimeout:          time.Second * 5,
		Net_ReadTimeout:          time.Second * 60,
		Net_WriteTimeout:         time.Second * 60,
		Net_KeepAlive:            0,
		Consumer_Offsets_Initial: -1,
	}
	for _, o := range opts {
		o(&options)
	}
	return options
}
