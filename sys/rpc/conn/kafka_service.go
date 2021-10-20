package conn

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/liwei1dao/lego"
	"github.com/liwei1dao/lego/sys/kafka"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc/core"
)

func NewKafkaService(serviceId string, kafkahost []string, rpcId string) (kafkaService *KafkaService, err error) {
	var (
		kfk kafka.IKafka
	)
	if kfk, err = kafka.NewSys(
		kafka.SetStartType(kafka.AsyncproducerAndConsumer),
		kafka.SetHosts(kafkahost),
		kafka.SetClientID(serviceId),
		kafka.SetTopics([]string{rpcId}),
		kafka.SetGroupId(serviceId),
		kafka.SetConsumer_Offsets_Initial(-1),
	); err != nil {
		err = fmt.Errorf("RPC NewKafkaService kafka.NewSys err:%v", err)
		return
	}
	kafkaService = &KafkaService{
		kafka: kfk,
	}
	return
}

type KafkaService struct {
	service core.IRpcServer
	kafka   kafka.IKafka
}

func (this *KafkaService) Start(service core.IRpcServer) (err error) {
	this.service = service
	go this.on_request_handle()
	return
}

func (this *KafkaService) Stop() (err error) {
	err = this.kafka.Close()
	return
}

func (this *KafkaService) Callback(callinfo core.CallInfo) error {
	body, _ := this.MarshalResult(callinfo.Result)
	reply_to := callinfo.Props["reply_to"].(string)
	// log.Debugf("RPC KafkaService Callback reply_to:%v", reply_to)
	this.kafka.Asyncproducer_Input() <- &sarama.ProducerMessage{
		Topic: reply_to,
		Value: sarama.ByteEncoder(body),
	}
	return nil
}

func (this *KafkaService) on_request_handle() {
	defer lego.Recover("RPC KafkaService")
	for v := range this.kafka.Consumer_Messages() {
		// log.Debugf("RPC KafkaService Receive: %+v", v)
		rpcInfo, err := this.Unmarshal(v.Value)
		if err == nil {
			callInfo := &core.CallInfo{
				RpcInfo: *rpcInfo,
			}
			callInfo.Props = map[string]interface{}{
				"reply_to": rpcInfo.ReplyTo,
			}
			callInfo.Agent = this //设置代理为NatsServer
			this.service.Call(*callInfo)
		} else {
			fmt.Println("error ", err)
		}
	}
	log.Debugf("RPC KafkaService on_request_handle exit")
}

func (s *KafkaService) Unmarshal(data []byte) (*core.RPCInfo, error) {
	var rpcInfo core.RPCInfo
	err := proto.Unmarshal(data, &rpcInfo)
	if err != nil {
		return nil, err
	} else {
		return &rpcInfo, err
	}
}

func (s *KafkaService) MarshalResult(resultInfo core.ResultInfo) ([]byte, error) {
	b, err := proto.Marshal(&resultInfo)
	return b, err
}
