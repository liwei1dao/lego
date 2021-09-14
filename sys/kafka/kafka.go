package kafka

import (
	"time"

	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
)

func newSys(options Options) (sys *Kafka, err error) {
	sys = &Kafka{options: options}
	err = sys.init()
	return
}

type Kafka struct {
	options       Options
	syncproducer  sarama.SyncProducer
	asyncproducer sarama.AsyncProducer
	consumer      *cluster.Consumer
}

func (this *Kafka) init() (err error) {
	config := sarama.NewConfig()
	config.ClientID = this.options.ClientID
	config.Producer.MaxMessageBytes = this.options.Producer_MaxMessageBytes
	config.Producer.RequiredAcks = this.options.Producer_RequiredAcks
	config.Producer.Return.Successes = this.options.Producer_Return_Successes
	config.Producer.Return.Errors = this.options.Producer_Return_Errors
	config.Producer.Retry.Max = this.options.Producer_Retry_Max
	config.Producer.Retry.Backoff = time.Duration(this.options.Producer_Retry_Backoff) * time.Millisecond
	config.Producer.Compression = this.options.Producer_Compression
	config.Producer.CompressionLevel = this.options.Producer_CompressionLevel
	config.Net.DialTimeout = this.options.Net_DialTimeout
	config.Net.ReadTimeout = this.options.Net_ReadTimeout
	config.Net.WriteTimeout = this.options.Net_WriteTimeout
	config.Net.KeepAlive = this.options.Net_KeepAlive
	if this.options.StartType == Syncproducer || this.options.StartType == All || this.options.StartType == SyncproducerAndConsumer {
		if this.syncproducer, err = sarama.NewSyncProducer(this.options.Hosts, config); err != nil {
			return
		}
	}
	if this.options.StartType == Asyncproducer || this.options.StartType == All || this.options.StartType == AsyncproducerAndConsumer {
		if this.asyncproducer, err = sarama.NewAsyncProducer(this.options.Hosts, config); err != nil {
			return
		}
	}
	if this.options.StartType == Consumer || this.options.StartType == All || this.options.StartType == SyncproducerAndConsumer || this.options.StartType == AsyncproducerAndConsumer {
		config := cluster.NewConfig()
		config.Consumer.Return.Errors = this.options.Consumer_Return_Errors
		config.Group.Return.Notifications = true
		config.Consumer.Offsets.Initial = this.options.Consumer_Offsets_Initial
		if this.consumer, err = cluster.NewConsumer(this.options.Hosts, this.options.GroupId, this.options.Topics, config); err != nil {
			return
		}
	}
	return
}

func (this *Kafka) Syncproducer_SendMessage(msg *sarama.ProducerMessage) (partition int32, offset int64, err error) {
	return this.syncproducer.SendMessage(msg)
}

func (this *Kafka) Syncproducer_SendMessages(msgs []*sarama.ProducerMessage) error {
	return this.syncproducer.SendMessages(msgs)
}

func (this *Kafka) Syncproducer_Close() error {
	return this.syncproducer.Close()
}

func (this *Kafka) Asyncproducer_Input() chan<- *sarama.ProducerMessage {
	return this.asyncproducer.Input()
}
func (this *Kafka) Asyncproducer_Successes() <-chan *sarama.ProducerMessage {
	return this.asyncproducer.Successes()
}
func (this *Kafka) Asyncproducer_Errors() <-chan *sarama.ProducerError {
	return this.asyncproducer.Errors()
}
func (this *Kafka) Asyncproducer_AsyncClose() {
	this.asyncproducer.AsyncClose()
}
func (this *Kafka) Asyncproducer_Close() error {
	return this.asyncproducer.Close()
}

func (this *Kafka) Consumer_Messages() <-chan *sarama.ConsumerMessage {
	return this.consumer.Messages()
}

func (this *Kafka) Consumer_Partitions() <-chan cluster.PartitionConsumer {
	return this.consumer.Partitions()
}

func (this *Kafka) Consumer_Notifications() <-chan *cluster.Notification {
	return this.consumer.Notifications()
}

func (this *Kafka) Consumer_Errors() <-chan error {
	return this.consumer.Errors()
}

func (this *Kafka) Consumer_MarkOffset(msg *sarama.ConsumerMessage, metadata string) {
	this.consumer.MarkOffset(msg, metadata)
}

func (this *Kafka) Consumer_MarkOffsets(s *cluster.OffsetStash) {
	this.consumer.MarkOffsets(s)
}

func (this *Kafka) Consumer_Close() (err error) {
	return this.consumer.Close()
}
func (this *Kafka) Close() (err error) {
	if this.syncproducer != nil {
		err = this.syncproducer.Close()
	}
	if this.asyncproducer != nil {
		err = this.asyncproducer.Close()
	}
	if this.consumer != nil {
		err = this.consumer.Close()
	}
	return
}
