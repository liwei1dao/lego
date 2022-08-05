package kafka

import (
	"github.com/Shopify/sarama"
)

type (
	IConsumer interface {
		Consumer_Messages() <-chan *sarama.ConsumerMessage
		Consumer_Errors() <-chan error
		Consumer_Close() error
	}
	ISys interface {
		Topics() ([]string, error)
		Partitions(topic string) ([]int32, error)
		GetOffset(topic string, partitionID int32, time int64) (int64, error)
		Syncproducer_SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
		Syncproducer_SendMessages(msgs []*sarama.ProducerMessage) error
		Syncproducer_Close() error
		Asyncproducer_Input() chan<- *sarama.ProducerMessage
		Asyncproducer_Successes() <-chan *sarama.ProducerMessage
		Asyncproducer_Errors() <-chan *sarama.ProducerError
		Asyncproducer_AsyncClose()
		Asyncproducer_Close() error
		Consumer_Messages() <-chan *sarama.ConsumerMessage
		Consumer_Errors() <-chan error
		Consumer_Close() error
		Close() (err error)
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, opt ...Option) (err error) {
	var option *Options
	if option, err = newOptions(config, opt...); err != nil {
		return
	}
	defsys, err = newSys(option)
	return
}

func NewSys(opt ...Option) (sys ISys, err error) {
	var option *Options
	if option, err = newOptionsByOption(opt...); err != nil {
		return
	}
	sys, err = newSys(option)
	return
}

func Topics() ([]string, error) {
	return defsys.Topics()
}
func Partitions(topic string) ([]int32, error) {
	return defsys.Partitions(topic)
}
func GetOffset(topic string, partitionID int32, time int64) (int64, error) {
	return defsys.GetOffset(topic, partitionID, time)
}

func Syncproducer_SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return defsys.Syncproducer_SendMessage(msg)
}

func Syncproducer_SendMessages(msgs []*sarama.ProducerMessage) error {
	return defsys.Syncproducer_SendMessages(msgs)
}
func Syncproducer_Close() error {
	return defsys.Syncproducer_Close()
}
func Asyncproducer_Input() chan<- *sarama.ProducerMessage {
	return defsys.Asyncproducer_Input()
}

func Asyncproducer_Successes() <-chan *sarama.ProducerMessage {
	return defsys.Asyncproducer_Successes()
}
func Asyncproducer_Errors() <-chan *sarama.ProducerError {
	return defsys.Asyncproducer_Errors()
}
func Asyncproducer_AsyncClose() {
	defsys.Asyncproducer_AsyncClose()
}

func Asyncproducer_Close() error {
	return defsys.Asyncproducer_Close()
}

func Consumer_Messages() <-chan *sarama.ConsumerMessage {
	return defsys.Consumer_Messages()
}
func Consumer_Errors() <-chan error {
	return defsys.Consumer_Errors()
}
func Consumer_Close() error {
	return defsys.Consumer_Close()
}
func Close() (err error) {
	return defsys.Close()
}
