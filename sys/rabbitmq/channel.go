package rabbitmq

import "github.com/streadway/amqp"

type Channel struct {
	channel *amqp.Channel
	queue   amqp.Queue
}

/*
方法描述:异步模式生产
mandatory:是否强制交换机
immediate: 是否等待确认
msg:消息
*/
func (this *Channel) Producer_SendAsync(mandatory, immediate bool, msg amqp.Publishing) (err error) {
	err = this.channel.Publish(
		"", // 默认交换机
		this.queue.Name,
		mandatory, // 是否强制交换机
		immediate, // 是否等待确认
		msg,
	)
	return
}

/*
方法描述:同步模式生产
mandatory:是否强制交换机
immediate: 是否等待确认
msg:消息
*/
func (this *Channel) Producer_Send(mandatory, immediate bool, msg amqp.Publishing) (err error) {
	err = this.channel.Publish(
		"", // 默认交换机
		this.queue.Name,
		mandatory, // 是否强制交换机
		immediate, // 是否等待确认
		msg,
	)
	return
}

/*
方法描述:异步生产者
autoAck:是否自动确认消息
exclusive: 是否独占队列
noLocal: 是否阻塞（消费者无法关闭）
noWait: 是否等待
msg:消息
*/
func (this *Channel) Consumer_Messages(autoAck, exclusive, noLocal, noWait bool) (output <-chan amqp.Delivery, err error) {
	// 从队列中获取消息
	output, err = this.channel.Consume(
		this.queue.Name, // 队列名
		"",              // 消费者标识符
		autoAck,         // 是否自动确认消息
		exclusive,       // 是否独占队列
		noLocal,         // 是否阻塞（消费者无法关闭）
		noWait,          // 是否等待
		nil,             // 附加参数
	)
	return
}
