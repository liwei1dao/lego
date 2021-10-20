package kafka

import (
	"context"
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/liwei1dao/lego/sys/log"
)

// type MessagChan <-chan *sarama.ConsumerMessage

func newConsumerGroup(brokers []string, group string, topics []string, config *sarama.Config) (kcgroup *KafkaConsumerGroup, err error) {
	var (
		cgroup sarama.ConsumerGroup
	)
	if cgroup, err = sarama.NewConsumerGroup(brokers, group, config); err != nil {
		err = fmt.Errorf("newConsumerGroup err:%v\n", err)
		return
	}
	ctx, cancel := context.WithCancel(context.Background())
	kcgroup = &KafkaConsumerGroup{
		cgroup:   cgroup,
		ready:    make(chan bool),
		topics:   topics,
		ctx:      ctx,
		cancel:   cancel,
		wg:       new(sync.WaitGroup),
		messages: make(chan *sarama.ConsumerMessage, 10),
	}
	go kcgroup.run()
	<-kcgroup.ready // Await till the consumer has been set up
	return
}

type KafkaConsumerGroup struct {
	ready    chan bool
	topics   []string
	cgroup   sarama.ConsumerGroup
	wg       *sync.WaitGroup
	ctx      context.Context
	cancel   context.CancelFunc
	messages chan *sarama.ConsumerMessage
}

func (this *KafkaConsumerGroup) Consumer_Messages() <-chan *sarama.ConsumerMessage {
	return this.messages
}

func (this *KafkaConsumerGroup) Consumer_Close() (err error) {
	this.cancel()
	this.wg.Wait()
	err = this.cgroup.Close()
	close(this.messages)
	return
}

func (this *KafkaConsumerGroup) run() {
	this.wg.Add(1)
	defer this.wg.Done()
locp:
	for {
		if err := this.cgroup.Consume(this.ctx, this.topics, this); err != nil {
			log.Errorf("Error from consumer: %v", err)
		}
		if this.ctx.Err() != nil {
			break locp
		}
		this.ready = make(chan bool)
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim
func (this *KafkaConsumerGroup) Setup(sarama.ConsumerGroupSession) error {
	close(this.ready)
	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited
func (this *KafkaConsumerGroup) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (this *KafkaConsumerGroup) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	log.Debugf("Sarama consumer ConsumeClaim%s...1", session.MemberID())
	for message := range claim.Messages() {
		// log.Printf("Message claimed: value = %s, timestamp = %v, topic = %s", string(message.Value), message.Timestamp, message.Topic)
		this.messages <- message
		session.MarkMessage(message, "")
	}
	log.Debugf("Sarama consumer ConsumeClaim%s...2", session.MemberID())
	return nil
}
