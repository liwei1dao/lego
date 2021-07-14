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
	config.Net.DialTimeout = time.Duration(this.options.Net_DialTimeout) * time.Second
	config.Net.ReadTimeout = time.Duration(this.options.Net_ReadTimeout) * time.Second
	config.Net.WriteTimeout = time.Duration(this.options.Net_WriteTimeout) * time.Second
	config.Net.KeepAlive = time.Duration(this.options.Net_KeepAlive) * time.Second
	if this.options.StartType == Syncproducer || this.options.StartType == All {
		if this.syncproducer, err = sarama.NewSyncProducer(this.options.Hosts, config); err != nil {
			return
		}
	}
	if this.options.StartType == Asyncproducer || this.options.StartType == All {
		if this.asyncproducer, err = sarama.NewAsyncProducer(this.options.Hosts, config); err == nil {
			return
		}
	}
	if this.options.StartType == Consumer || this.options.StartType == All {
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

func (this *Kafka) Consumer_Notifications() <-chan *cluster.Notification {
	return this.consumer.Notifications()
}

func (this *Kafka) Consumer_Errors() <-chan error {
	return this.consumer.Errors()
}
