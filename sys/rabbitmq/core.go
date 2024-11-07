package rabbitmq

import "github.com/streadway/amqp"

type (
	ISys interface {
		Producer_SendAsync(msg amqp.Publishing) (err error)
		Producer_Send(msg amqp.Publishing) (err error)
		Consumer_Messages() (output <-chan amqp.Delivery, err error)
	}
	IChannel interface {
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

func Producer_SendAsync(msg amqp.Publishing) (err error) {
	err = defsys.Producer_SendAsync(msg)
	return
}
func Producer_Send(msg amqp.Publishing) (err error) {
	err = defsys.Producer_Send(msg)
	return
}
func Consumer_Messages() (output <-chan amqp.Delivery, err error) {
	output, err = defsys.Consumer_Messages()
	return
}
