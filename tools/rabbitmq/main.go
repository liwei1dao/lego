package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rabbitmq"
	"github.com/streadway/amqp"
)

/*
	kafka 测试工具
*/

var (
	_addr    = flag.String("a", "amqp://root:li13451234@localhost:5672/", "rabbitmq服务地址")
	_channel = flag.String("c", "test001", "数据管道")
)
var (
	sys rabbitmq.ISys
	err error
)

func main() {
	flag.Parse()
	if err = log.OnInit(nil, log.SetFileName("./rabbitmq.log")); err != nil {
		fmt.Printf("log init err:%v", err)
		return
	}
	log.Debugf("rabbitmq init addr:%s channel:%s", *_addr, *_channel)
	if sys, err = rabbitmq.NewSys(
		rabbitmq.SetRabbitmqUrl(*_addr),
		rabbitmq.SetChannelName(*_channel),
	); err != nil {
		log.Errorf("rabbitmq init err:%v", err)
		return
	}
	go consumer()
	go producer()
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		log.Debugf("test end signal")
	}
}

// 生产
func producer() {
	c := time.NewTicker(time.Second)
	defer c.Stop()
	for {
		select {
		case <-c.C:
			sys.Producer_SendAsync(false, false, amqp.Publishing{
				Body: []byte(fmt.Sprintf("liwei1dao:%d", time.Now().Unix())),
			})
			break
		}
	}
}

// 消费
func consumer() {
	output, err := sys.Consumer_Messages(true, false, false, false)
	if err != nil {
		log.Errorln("启动消费异常", err)
	}
	for v := range output {
		log.Debugf("Consumer:%s", string(v.Body))
	}
}
