package pulsar

import (
	"context"

	"github.com/apache/pulsar-client-go/pulsar"
)

type (
	ProducerError struct {
		Msg *pulsar.ProducerMessage
		Err error
	}
	ISys interface {
		Producer_Errors() <-chan *ProducerError
		Producer_SendAsync() chan<- *pulsar.ProducerMessage
		Producer_Send(msg *pulsar.ProducerMessage) (pulsar.MessageID, error)
		Consumer_Receive(ctx context.Context) (pulsar.Message, error)
		Consumer_Messages() <-chan pulsar.ConsumerMessage
		Close()
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

func Producer_Errors() <-chan *ProducerError {
	return defsys.Producer_Errors()
}
func Producer_SendAsync() chan<- *pulsar.ProducerMessage {
	return defsys.Producer_SendAsync()
}

func Producer_Send(msg *pulsar.ProducerMessage) (pulsar.MessageID, error) {
	return defsys.Producer_Send(msg)
}
func Consumer_Receive(ctx context.Context) (pulsar.Message, error) {
	return defsys.Consumer_Receive(ctx)
}
func Consumer_Messages() <-chan pulsar.ConsumerMessage {
	return defsys.Consumer_Messages()
}
func Close() {
	defsys.Close()
}
