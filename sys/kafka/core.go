package kafka

import (
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

type (
	IKafka interface {
		Syncproducer_SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error)
		Syncproducer_SendMessages(msgs []*sarama.ProducerMessage) error
		Syncproducer_Close() error
		Asyncproducer_Input() chan<- *sarama.ProducerMessage
		Asyncproducer_Successes() <-chan *sarama.ProducerMessage
		Asyncproducer_Errors() <-chan *sarama.ProducerError
		Asyncproducer_AsyncClose()
		Asyncproducer_Close() error
		Consumer_Messages() <-chan *sarama.ConsumerMessage
		Consumer_Partitions() <-chan cluster.PartitionConsumer
		Consumer_Notifications() <-chan *cluster.Notification
		Consumer_Errors() <-chan error
		Consumer_MarkOffset(msg *sarama.ConsumerMessage, metadata string)
		Consumer_MarkOffsets(s *cluster.OffsetStash)
		Consumer_Close() (err error)
	}
)

var (
	defsys IKafka
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	if defsys, err = newSys(newOptions(config, option...)); err == nil {
	}
	return
}

func NewSys(option ...Option) (sys IKafka, err error) {
	if sys, err = newSys(newOptionsByOption(option...)); err == nil {
	}
	return
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

func Consumer_Partitions() <-chan cluster.PartitionConsumer {
	return defsys.Consumer_Partitions()
}
func Consumer_Notifications() <-chan *cluster.Notification {
	return defsys.Consumer_Notifications()
}

func Consumer_Errors() <-chan error {
	return defsys.Consumer_Errors()
}

func Consumer_MarkOffset(msg *sarama.ConsumerMessage, metadata string) {
	defsys.Consumer_MarkOffset(msg, metadata)
}

func Consumer_MarkOffsets(s *cluster.OffsetStash) {
	defsys.Consumer_MarkOffsets(s)
}

func Consumer_Close() (err error) {
	return defsys.Consumer_Close()
}
