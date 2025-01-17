package rabbitmq

import "github.com/streadway/amqp"

type (
	IChannel interface {
		Producer_SendAsync(mandatory, immediate bool, msg amqp.Publishing) (err error)
		Producer_Send(mandatory, immediate bool, msg amqp.Publishing) (err error)
		Consumer_Messages(autoAck, exclusive, noLocal, noWait bool) (output <-chan amqp.Delivery, err error)
	}
	ISys interface {
		IChannel
		NewChannel(channelName string) (channel IChannel, err error)
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	if defsys, err = newSys(newOptions(config, option...)); err == nil {
	}
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	if sys, err = newSys(newOptionsByOption(option...)); err == nil {
	}
	return
}

func Producer_SendAsync(mandatory, immediate bool, msg amqp.Publishing) (err error) {
	err = defsys.Producer_SendAsync(mandatory, immediate, msg)
	return
}
func Producer_Send(mandatory, immediate bool, msg amqp.Publishing) (err error) {
	err = defsys.Producer_Send(mandatory, immediate, msg)
	return
}
func Consumer_Messages(autoAck, exclusive, noLocal, noWait bool) (output <-chan amqp.Delivery, err error) {
	output, err = defsys.Consumer_Messages(autoAck, exclusive, noLocal, noWait)
	return
}
