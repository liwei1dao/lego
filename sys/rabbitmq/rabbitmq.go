package rabbitmq

import (
	"github.com/streadway/amqp"
)

func newSys(options Options) (sys *Rabbitmq, err error) {
	sys = &Rabbitmq{options: options}
	err = sys.init()
	return
}

type Rabbitmq struct {
	options Options
	conn    *amqp.Connection
	channel IChannel
}

func (this *Rabbitmq) init() (err error) {
	if this.conn, err = amqp.Dial(this.options.RabbitmqUrl); err != nil {
		return
	}
	if this.options.ChannelName != "" {
		this.channel, err = this.NewChannel(this.options.ChannelName)
	}
	return
}

// 新建 channel
func (this *Rabbitmq) NewChannel(channelName string) (channel IChannel, err error) {
	var (
		_channel *amqp.Channel
		_queue   amqp.Queue
	)

	// 创建一个 channel（信道）
	if _channel, err = this.conn.Channel(); err != nil {
		return
	}
	// 声明队列（如果队列不存在，则会创建）
	if _queue, err = _channel.QueueDeclare(
		this.options.ChannelName, // 队列名
		true,                     // 是否持久化
		false,                    // 是否自动删除
		false,                    // 是否排他（用于临时队列）
		false,                    // 是否等待
		nil,                      // 附加参数
	); err != nil {
		return
	}
	channel = &Channel{
		channel: _channel,
		queue:   _queue,
	}
	return
}

// 异步模式生产
func (this *Rabbitmq) Producer_SendAsync(mandatory, immediate bool, msg amqp.Publishing) (err error) {
	err = this.channel.Producer_SendAsync(mandatory, immediate, msg)
	return
}

// 同步模式生产
func (this *Rabbitmq) Producer_Send(mandatory, immediate bool, msg amqp.Publishing) (err error) {
	err = this.channel.Producer_Send(mandatory, immediate, msg)
	return
}

// 异步生产者
func (this *Rabbitmq) Consumer_Messages(autoAck, exclusive, noLocal, noWait bool) (output <-chan amqp.Delivery, err error) {
	output, err = this.channel.Consumer_Messages(autoAck, exclusive, noLocal, noWait)
	return
}
