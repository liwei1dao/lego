package pulsar

import (
	"context"
	"fmt"
	"time"

	"github.com/apache/pulsar-client-go/pulsar"
)

func newSys(options Options) (sys *Pulsar, err error) {
	sys = &Pulsar{options: options, ctx: context.Background()}
	err = sys.init()
	return
}

type Pulsar struct {
	options  Options
	client   pulsar.Client
	producer pulsar.Producer
	consumer pulsar.Consumer
	ctx      context.Context
	errors   chan *ProducerError
	input    chan *pulsar.ProducerMessage
}

func (this *Pulsar) init() (err error) {
	if this.client, err = pulsar.NewClient(
		pulsar.ClientOptions{
			URL:               this.options.PulsarUrl,
			OperationTimeout:  30 * time.Second,
			ConnectionTimeout: 30 * time.Second},
	); err != nil {
		err = fmt.Errorf("Sys:Pulsar connect err:%v", err)
		return
	}
	if this.options.StartType == Producer || this.options.StartType == All {
		if this.producer, err = this.client.CreateProducer(pulsar.ProducerOptions{
			Topic:            this.options.Topics[0],
			BatchingMaxSize:  this.options.Producer_BatchingMaxSize,
			CompressionType:  this.options.Producer_CompressionType,
			CompressionLevel: this.options.Producer_CompressionLevel,
		}); err != nil {
			err = fmt.Errorf("Sys:Pulsar CreateProducer err:%v", err)
			return
		}
		this.input = make(chan *pulsar.ProducerMessage)
		if this.options.Producer_Return_Errors {
			this.errors = make(chan *ProducerError)
		}
	}
	if this.options.StartType == Consumer || this.options.StartType == All {
		if this.consumer, err = this.client.Subscribe(pulsar.ConsumerOptions{
			Topics:           this.options.Topics,
			SubscriptionName: this.options.ConsumerGroupId,
			Type:             pulsar.Shared,
		}); err != nil {
			err = fmt.Errorf("Sys:Pulsar Subscribe err:%v", err)
			return
		}
	}
	this.run()
	return
}

func (this *Pulsar) Producer_Errors() <-chan *ProducerError {
	return this.errors
}
func (this *Pulsar) Producer_SendAsync() chan<- *pulsar.ProducerMessage {
	return this.input
}

func (this *Pulsar) Producer_Send(msg *pulsar.ProducerMessage) (pulsar.MessageID, error) {
	return this.producer.Send(this.ctx, msg)
}

func (this *Pulsar) Consumer_Receive(ctx context.Context) (pulsar.Message, error) {
	return this.consumer.Receive(ctx)
}

func (this *Pulsar) Consumer_Messages() <-chan pulsar.ConsumerMessage {
	return this.consumer.Chan()
}

func (this *Pulsar) Close() {
	if this.producer != nil {
		close(this.input)
		this.producer.Close()
	}
	if this.consumer != nil {
		this.consumer.Close()
	}
	return
}

func (this *Pulsar) run() {
	if this.producer != nil {
		go func() {
			for v := range this.input {
				this.producer.SendAsync(this.ctx, v, this.producer_SendAsync_Call)
			}
		}()
	}
}

func (this *Pulsar) producer_SendAsync_Call(mi pulsar.MessageID, pm *pulsar.ProducerMessage, err error) {
	if this.options.Producer_Return_Errors && err != nil {
		this.errors <- &ProducerError{Msg: pm, Err: err}
	}
}
