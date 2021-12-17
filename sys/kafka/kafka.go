package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
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
	consumerGroup *KafkaConsumerGroup
}

func (this *Kafka) init() (err error) {
	var (
		version sarama.KafkaVersion
	)
	config := sarama.NewConfig()
	version, err = sarama.ParseKafkaVersion(this.options.Version)
	if err != nil {
		err = fmt.Errorf("Kafka version:%s err:%v", this.options.Version, err)
		return
	}
	config.Version = version
	config.ClientID = this.options.ClientID
	config.Net.SASL.Enable = this.options.Sasl_Enable
	config.Net.SASL.Mechanism = this.options.Sasl_Mechanism
	config.Net.SASL.GSSAPI = this.options.Sasl_GSSAPI
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
	config.Consumer.Offsets.Initial = this.options.Consumer_Offsets_Initial
	config.Consumer.Return.Errors = this.options.Consumer_Return_Errors
	switch this.options.Consumer_Assignor {
	case "sticky":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategySticky
	case "roundrobin":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRoundRobin
	case "range":
		config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	default:
		err = fmt.Errorf("Unrecognized consumer group partition assignor: %s", this.options.Consumer_Assignor)
		return
	}
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
		if this.consumerGroup, err = newConsumerGroup(this.options.Hosts, this.options.GroupId, this.options.Topics, config); err != nil {
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
	return this.consumerGroup.Consumer_Messages()
}

func (this *Kafka) Consumer_Errors() <-chan error {
	return this.consumerGroup.Consumer_Errors()
}
func (this *Kafka) Consumer_Close() (err error) {
	return this.consumerGroup.Consumer_Close()
}

func (this *Kafka) Close() (err error) {
	if this.syncproducer != nil {
		err = this.syncproducer.Close()
	}
	if this.asyncproducer != nil {
		err = this.asyncproducer.Close()
	}
	if this.consumerGroup != nil {
		err = this.consumerGroup.Consumer_Close()
	}
	return
}
