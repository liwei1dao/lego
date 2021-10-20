package kafka

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/Shopify/sarama"
)

func Test_sys(t *testing.T) {
	if err := OnInit(map[string]interface{}{
		"StartType":                 Consumer,
		"Hosts":                     []string{"172.20.27.126:9092", "172.20.27.127:9092", " 172.20.27.128:9092"},
		"Topics":                    []string{"ETL-IN-CX20211013641233861142316500871138900239"},
		"GroupId":                   "liwei3dao",
		"ClientID":                  "test",
		"Net_KeepAlive":             0,
		"Producer_MaxMessageBytes":  201231293891,
		"Net_DialTimeout":           5 * time.Second,
		"Producer_Compression":      sarama.CompressionGZIP,
		"Producer_CompressionLevel": 1,
		"Producer_Retry_Max":        3,
	}); err != nil {
		fmt.Printf("start sys err:%v", err)
	} else {
		fmt.Printf("start sys succ")
		// Syncproducer_SendMessage(&sarama.ProducerMessage{
		// 	Topic: "send",
		// 	Value: sarama.StringEncoder("asdasd"),
		// })
	}
}

func Test_newConsumerGroup(t *testing.T) {
	config := sarama.NewConfig()
	version, err := sarama.ParseKafkaVersion("2.1.1")
	if err != nil {
		fmt.Printf("Error parsing Kafka version: %v\n", err)
		return
	}
	config.Version = version
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	cgroup, err := newConsumerGroup([]string{
		"172.20.27.126:9092", "172.20.27.127:9092", "172.20.27.128:9092"},
		"liwei4dao",
		[]string{"ETL-IN-CX20211013641233861142316500871138900249"},
		config)
	if err != nil {
		fmt.Printf(" newConsumerGroup err: %v\n", err)
		return
	}
	go func() {
		for v := range cgroup.Consumer_Messages() {
			fmt.Printf("采集数据 Partition:%d Offset:%d\n", v.Partition, v.Offset)
		}
	}()

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		fmt.Printf("terminating: via signal\n")
	}
	if err = cgroup.Consumer_Close(); err != nil {
		fmt.Printf("Consumer_Close err: %v\n", err)
		return
	} else {
		fmt.Printf("Consumer_Close end: %v\n", err)
	}
}
func Test_ConsumerGroup(t *testing.T) {
	if err := OnInit(map[string]interface{}{
		"StartType":                 AsyncproducerAndConsumer,
		"Hosts":                     []string{"172.20.27.126:9092", "172.20.27.127:9092", " 172.20.27.128:9092"},
		"Topics":                    []string{"ETL-IN-CX20211013641233861142316500871138900249"},
		"GroupId":                   "liwei4dao",
		"ClientID":                  "test",
		"Net_KeepAlive":             0,
		"Producer_MaxMessageBytes":  201231293891,
		"Net_DialTimeout":           5 * time.Second,
		"Producer_Compression":      sarama.CompressionGZIP,
		"Producer_CompressionLevel": 1,
		"Producer_Retry_Max":        3,
	}); err != nil {
		fmt.Printf("start sys err:%v\n", err)
	} else {
		fmt.Printf("start sys succ\n")
		go func() {
			for v := range Consumer_Messages() {
				fmt.Printf("Consumer_Messages:%d\n", v.Partition)
			}
		}()
	}
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		fmt.Printf("terminating: via signal\n")
	}
	if err := Consumer_Close(); err != nil {
		fmt.Printf("Consumer_Close err: %v\n", err)
		return
	} else {
		fmt.Printf("Consumer_Close end: %v\n", err)
	}
}
