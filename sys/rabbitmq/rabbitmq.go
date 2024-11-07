package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

func newSys(options Options) (sys *Rabbitmq, err error) {
	sys = &Rabbitmq{options: options}
	err = sys.init()
	return
}

type Rabbitmq struct {
	options  Options
	conn     *amqp.Connection
	channel  *amqp.Channel
	queue    amqp.Queue
	confirms chan amqp.Confirmation
}

func (this *Rabbitmq) init() (err error) {
	if this.conn, err = amqp.Dial(this.options.RabbitmqUrl); err != nil {
		return
	}
	// 创建一个 channel（信道）
	if this.channel, err = this.conn.Channel(); err != nil {
		return
	}
	// 声明队列（如果队列不存在，则会创建）
	if this.queue, err = this.channel.QueueDeclare(
		this.options.ChannelName, // 队列名
		true,                     // 是否持久化
		false,                    // 是否自动删除
		false,                    // 是否排他（用于临时队列）
		false,                    // 是否等待
		nil,                      // 附加参数
	); err != nil {
		return
	}
	if this.options.StartType != Consumer && this.options.ConfirmNotif { //开启生产确认通知
		// 启用发布确认模式
		if err := this.channel.Confirm(false); err != nil {
			log.Fatalf("无法启用发布确认模式: %s", err)
		}
		this.confirms = this.channel.NotifyPublish(make(chan amqp.Confirmation, 2))
	}
	return
}

// 生产确认队列
func (this *Rabbitmq) Producer_Confirm() <-chan amqp.Confirmation {
	return this.confirms
}

// 异步模式生产
func (this *Rabbitmq) Producer_SendAsync(msg amqp.Publishing) (err error) {
	err = this.channel.Publish(
		"", // 默认交换机
		this.queue.Name,
		false, // 是否强制交换机
		false, // 是否等待确认
		msg,
	)
	return
}

// 同步模式生产
func (this *Rabbitmq) Producer_Send(msg amqp.Publishing) (err error) {
	err = this.channel.Publish(
		"", // 默认交换机
		this.queue.Name,
		false, // 是否强制交换机
		true,  // 是否等待确认
		msg,
	)
	return
}

// 异步生产者
func (this *Rabbitmq) Consumer_Messages() (output <-chan amqp.Delivery, err error) {
	// 从队列中获取消息
	output, err = this.channel.Consume(
		this.queue.Name, // 队列名
		"",              // 消费者标识符
		true,            // 是否自动确认消息
		false,           // 是否独占队列
		false,           // 是否阻塞（消费者无法关闭）
		false,           // 是否等待
		nil,             // 附加参数
	)
	return
}
