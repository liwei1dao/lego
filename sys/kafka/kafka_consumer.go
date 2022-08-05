package kafka

import (
	"fmt"
	"sync"

	"github.com/Shopify/sarama"
	"github.com/liwei1dao/lego/sys/log"
)

func newConsumer(sys ISys, log log.ILogger, brokers []string, topic string, config *sarama.Config) (consumer *KafkaConsumer, err error) {
	consumer = &KafkaConsumer{
		sys:      sys,
		log:      log,
		topic:    topic,
		messages: make(chan *sarama.ConsumerMessage, 10),
		errors:   make(chan error, 0),
	}
	//创建消费者
	consumer.consumer, err = sarama.NewConsumer(brokers, config)
	if err != nil {
		log.Errorf("fail to start consumer, err:%v\n", err)
		return
	}
	//获取主题分区
	consumer.partitionList, err = consumer.consumer.Partitions(topic) // 根据topic取到所有的分区
	if err != nil {
		log.Errorf("fail to get list of partition:err%v\n", err)
		return
	}
	consumer.pc = make([]sarama.PartitionConsumer, len(consumer.partitionList))
	go consumer.run()
	return
}

type KafkaConsumer struct {
	sys           ISys
	log           log.ILogger
	consumer      sarama.Consumer
	wg            sync.WaitGroup
	topic         string
	partitionList []int32
	pc            []sarama.PartitionConsumer
	messages      chan *sarama.ConsumerMessage
	errors        chan error
}

func (this *KafkaConsumer) Consumer_Messages() <-chan *sarama.ConsumerMessage {
	return this.messages
}
func (this *KafkaConsumer) Consumer_Errors() <-chan error {
	return this.errors
}
func (this *KafkaConsumer) Consumer_Close() (err error) {
	err = this.consumer.Close()
	for _, v := range this.pc {
		if v != nil {
			v.Close()
		}
	}
	this.wg.Wait()
	close(this.messages)
	close(this.errors)
	return
}

func (this *KafkaConsumer) run() {
	for i, partition := range this.partitionList {
		//针对每个分区创建一个对应的分区消费者
		pc, err := this.consumer.ConsumePartition(this.topic, int32(partition), sarama.OffsetNewest)
		//pc, err := consumer.ConsumePartition("web_log", int32(partition), 90)
		if err != nil {
			this.log.Errorf("failed to start consumer for partition %d,err:%v\n", partition, err)
			return
		}
		this.pc[i] = pc
		// 异步从每个分区消费信息
		this.wg.Add(2) //+1
		go func(sarama.PartitionConsumer) {
			defer this.wg.Done() //-1
			for msg := range pc.Messages() {
				//fmt.Printf("Partition:%d Offset:%d Key:%v Value:%s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
				this.messages <- msg
			}
		}(pc)
		go func(sarama.PartitionConsumer) {
			defer this.wg.Done() //-1
			for err := range pc.Errors() {
				//fmt.Printf("Partition:%d Offset:%d Key:%v Value:%s\n", msg.Partition, msg.Offset, msg.Key, msg.Value)
				this.errors <- fmt.Errorf("Topic:%s Partition:%d err:%v", err.Topic, err.Partition, err.Err)
			}
		}(pc)

	}
	this.wg.Wait() //
}
