package kafka

import (
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/liwei1dao/lego/sys/log"
)

type logger struct {
	log log.ILogger
}

func (this *logger) Print(v ...interface{}) {
	this.log.Panicln(v...)
}
func (this *logger) Printf(format string, v ...interface{}) {
	this.log.Printf(format, v...)
}
func (this *logger) Println(v ...interface{}) {
	this.log.Panicln(v...)
}

func newSys(options *Options) (sys *Kafka, err error) {
	sys = &Kafka{options: options}
	err = sys.init()
	return
}

type Kafka struct {
	options       *Options
	client        sarama.Client
	syncproducer  sarama.SyncProducer
	asyncproducer sarama.AsyncProducer
	consumer      IConsumer
}

func (this *Kafka) init() (err error) {
	var (
		version sarama.KafkaVersion
	)
	if this.options.Debug {
		sarama.Logger = &logger{log: this.options.Log}
	}
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

	if this.options.StartType == Client || this.options.StartType == All {
		if this.client, err = sarama.NewClient(this.options.Hosts, config); err != nil {
			return
		}
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
		if this.options.GroupId != "" {
			if this.consumer, err = newConsumerGroup(this, this.options.Log, this.options.Hosts, this.options.GroupId, this.options.Topics, config); err != nil {
				return
			}
		} else {
			if this.consumer, err = newConsumer(this, this.options.Log, this.options.Hosts, this.options.Topics[0], config); err != nil {
				return
			}
		}
	}
	return
}

func (this *Kafka) Controller() (*sarama.Broker, error) {
	return this.client.Controller()
}

func (this *Kafka) RefreshController() (*sarama.Broker, error) {
	return this.client.RefreshController()
}

func (this *Kafka) Brokers() []*sarama.Broker {
	return this.client.Brokers()
}

func (this *Kafka) Broker(brokerID int32) (*sarama.Broker, error) {
	return this.client.Broker(brokerID)
}

func (this *Kafka) Topics() ([]string, error) {
	return this.client.Topics()
}

func (this *Kafka) Partitions(topic string) ([]int32, error) {
	return this.client.Partitions(topic)
}

func (this *Kafka) WritablePartitions(topic string) ([]int32, error) {
	return this.client.WritablePartitions(topic)
}
func (this *Kafka) Leader(topic string, partitionID int32) (*sarama.Broker, error) {
	return this.client.Leader(topic, partitionID)
}
func (this *Kafka) Replicas(topic string, partitionID int32) ([]int32, error) {
	return this.client.Replicas(topic, partitionID)
}
func (this *Kafka) InSyncReplicas(topic string, partitionID int32) ([]int32, error) {
	return this.client.InSyncReplicas(topic, partitionID)
}
func (this *Kafka) OfflineReplicas(topic string, partitionID int32) ([]int32, error) {
	return this.client.OfflineReplicas(topic, partitionID)
}
func (this *Kafka) RefreshBrokers(addrs []string) error {
	return this.client.RefreshBrokers(addrs)
}
func (this *Kafka) RefreshMetadata(topics ...string) error {
	return this.client.RefreshMetadata(topics...)
}
func (this *Kafka) GetOffset(topic string, partitionID int32, time int64) (int64, error) {
	return this.client.GetOffset(topic, partitionID, time)
}
func (this *Kafka) Coordinator(consumerGroup string) (*sarama.Broker, error) {
	return this.client.Coordinator(consumerGroup)
}
func (this *Kafka) RefreshCoordinator(consumerGroup string) error {
	return this.client.RefreshCoordinator(consumerGroup)
}

func (this *Kafka) InitProducerID() (*sarama.InitProducerIDResponse, error) {
	return this.client.InitProducerID()
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
	return this.consumer.Consumer_Messages()
}

func (this *Kafka) Consumer_Errors() <-chan error {
	return this.consumer.Consumer_Errors()
}
func (this *Kafka) Consumer_Close() (err error) {
	return this.consumer.Consumer_Close()
}

func (this *Kafka) Close() (err error) {
	if this.client != nil {
		if err = this.client.Close(); err != nil {
			this.options.Log.Errorf("client Close err:%v", err)
		}
	}
	if this.syncproducer != nil {
		if err = this.syncproducer.Close(); err != nil {
			this.options.Log.Errorf("syncproducer Close err:%v", err)
		}
	}
	if this.asyncproducer != nil {
		if err = this.asyncproducer.Close(); err != nil {
			this.options.Log.Errorf("asyncproducer Close err:%v", err)
		}
	}
	if this.consumer != nil {
		if err = this.consumer.Consumer_Close(); err != nil {
			this.options.Log.Errorf("consumer Close err:%v", err)
		}
	}
	return
}
