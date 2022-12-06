package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/Shopify/sarama"
	"github.com/liwei1dao/lego/sys/kafka"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/codec"
)

var (
	_starttype = flag.Int("s", 2, "kafka使用模式 0 生产者 1消费者 2同时生产和消费")
	_clientid  = flag.String("c", "test1", "kafka客户端id")
	_hosts     = flag.String("h", "120.79.200.115:10092", "kafka客户端id")
	_groupId   = flag.String("g", "test002", "kafka客户端id")
	_topics    = flag.String("t", "topicstest001", "数据管道")
)
var (
	err       error
	starttype kafka.KafkaStartType
	clientid  string
	hosts     []string
	groupId   string
	topics    []string
	sys       kafka.ISys
)

func main() {
	flag.Parse()
	if err = log.OnInit(nil, log.SetFileName("./kafka.log")); err != nil {
		fmt.Printf("log init err:%v", err)
		return
	}
	if *_starttype == 0 {
		starttype = kafka.Asyncproducer
	} else if *_starttype == 1 {
		starttype = kafka.Consumer
	} else {
		starttype = kafka.AsyncproducerAndConsumer
	}

	clientid = *_clientid
	hosts = strings.Split(*_hosts, ",")
	groupId = *_groupId
	topics = []string{*_topics}
	log.Debugf("kafka init hosts:%v topics:%v clientid:%s groupId:%s ", hosts, topics, clientid, groupId)
	if sys, err = kafka.NewSys(
		kafka.SetStartType(starttype),
		kafka.SetClientID(clientid),
		kafka.SetHosts(hosts),
		kafka.SetGroupId(groupId),
		kafka.SetTopics(topics),
		kafka.SetConsumer_Offsets_Initial(-1),
		kafka.SetProducer_Compression(sarama.CompressionGZIP),
		kafka.SetProducer_CompressionLevel(1),
		kafka.SetVersion("2.1.0"),
		kafka.SetDebug(true),
	); err != nil {
		log.Errorf("kafka init err:%v", err)
		return
	}
	if starttype == kafka.Consumer || starttype == kafka.AsyncproducerAndConsumer {
		log.Debugf("kafka init consumer:")
		go consumer()
		go consumererror()
	}
	if starttype == kafka.Asyncproducer || starttype == kafka.AsyncproducerAndConsumer {
		log.Debugf("kafka init producer:")
		go producer()
	}
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		log.Debugf("test end signal")
	}
}

//生产
func producer() {
	c := time.NewTicker(time.Second)
	defer c.Stop()
	for {
		select {
		case <-c.C:
			sys.Asyncproducer_Input() <- &sarama.ProducerMessage{
				Topic: topics[0],
				Value: sarama.StringEncoder(fmt.Sprintf("%d", time.Now().UnixMilli())),
			}
			break
		}
	}
}

//消费
func consumer() {
	for v := range sys.Consumer_Messages() {
		send := codec.StringToInt64(string(v.Value))
		log.Debugf("Consumer_Messages:send:%d receive:%d delay:%d", time.Now().UnixMilli(), send, time.Now().UnixMilli()-send)
	}
}

//消费
func consumererror() {
	for v := range sys.Consumer_Errors() {
		log.Debugf("Consumer_err:%v", v)
	}
}
