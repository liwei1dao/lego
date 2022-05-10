package conn

import (
	"fmt"

	"github.com/Shopify/sarama"
	"github.com/golang/protobuf/proto"
	"github.com/liwei1dao/lego/sys/kafka"
	"github.com/liwei1dao/lego/sys/rpc/core"
	"github.com/liwei1dao/lego/utils/container"
)

func NewKafkaClient(sys core.ISys, serviceId string, kafkaversion string, kafkahost []string, pushrpcId, receiveId string) (kafkaClient *KafkaClient, err error) {
	var (
		kfk kafka.IKafka
	)
	if kfk, err = kafka.NewSys(
		kafka.SetVersion(kafkaversion),
		kafka.SetStartType(kafka.AsyncproducerAndConsumer),
		kafka.SetClientID(serviceId),
		kafka.SetHosts(kafkahost),
		kafka.SetTopics([]string{receiveId}),
		kafka.SetConsumer_Offsets_Initial(-1),
	); err != nil {
		err = fmt.Errorf("RPC NewKafkaClient kafka.NewSys err:%v", err)
		return
	}
	kafkaClient = &KafkaClient{
		sys:               sys,
		callinfos:         container.NewBeeMap(),
		callbackqueueName: receiveId,
		rpcId:             pushrpcId,
		kafka:             kfk,
	}
	go kafkaClient.on_request_handle()
	return
}

type KafkaClient struct {
	sys               core.ISys
	callinfos         *container.BeeMap
	callbackqueueName string
	rpcId             string
	kafka             kafka.IKafka
}

func (this *KafkaClient) Stop() (err error) {
	err = this.kafka.Close()
	if this.callinfos != nil {
		//清理 callinfos 列表
		for key, clinetCallInfo := range this.callinfos.Items() {
			if clinetCallInfo != nil {
				//关闭管道
				close(clinetCallInfo.(core.ClinetCallInfo).Call)
				//从Map中删除
				this.callinfos.Delete(key)
			}
		}
		this.callinfos = nil
	} else {
		this.sys.Errorf("RCP NatsClient callinfos 异常为空")
	}
	return
}

func (this *KafkaClient) Delete(key string) (err error) {
	this.callinfos.Delete(key)
	return
}

/**
消息请求
*/
func (this *KafkaClient) Call(callInfo core.CallInfo, callback chan core.ResultInfo) (err error) {
	var (
		body []byte
	)
	if this.callinfos == nil {
		return fmt.Errorf("RPC NatsClient callinfos nil")
	}
	callInfo.RpcInfo.ReplyTo = this.callbackqueueName
	var correlation_id = callInfo.RpcInfo.Cid

	clinetCallInfo := &core.ClinetCallInfo{
		Correlation_id: correlation_id,
		Call:           callback,
		Timeout:        callInfo.RpcInfo.Expired,
	}
	this.callinfos.Set(correlation_id, *clinetCallInfo)
	body, err = this.Marshal(&callInfo.RpcInfo)
	if err != nil {
		return err
	}
	this.kafka.Asyncproducer_Input() <- &sarama.ProducerMessage{
		Topic: this.rpcId,
		Value: sarama.ByteEncoder(body),
	}
	return
}

/**
消息请求 不需要回复
*/
func (this *KafkaClient) CallNR(callInfo core.CallInfo) (err error) {
	var (
		body []byte
	)
	body, err = this.Marshal(&callInfo.RpcInfo)
	if err != nil {
		return err
	}
	this.kafka.Asyncproducer_Input() <- &sarama.ProducerMessage{
		Topic: this.rpcId,
		Value: sarama.ByteEncoder(body),
	}
	return
}

func (this *KafkaClient) on_request_handle() {
	defer func() {
		if r := recover(); r != nil {
			this.sys.Errorf("KafkaClient panic:", r)
		}
	}()
	for v := range this.kafka.Consumer_Messages() {
		// log.Debugf("RPC KafkaClient Receive: %+v", v)
		resultInfo, err := this.UnmarshalResult(v.Value)
		if err != nil {
			this.sys.Errorf("RPC NatsClient Unmarshal faild", err)
		} else {
			correlation_id := resultInfo.Cid
			clinetCallInfo := this.callinfos.Get(correlation_id)
			//删除
			this.callinfos.Delete(correlation_id)
			if clinetCallInfo != nil {
				clinetCallInfo.(core.ClinetCallInfo).Call <- *resultInfo
				close(clinetCallInfo.(core.ClinetCallInfo).Call)
			} else {
				//可能客户端已超时了，但服务端处理完还给回调了
				this.sys.Warnf("rpc callback no found : [%s]", correlation_id)
			}
		}
	}
	this.sys.Debugf("RPC KafkaClient on_request_handle exit")
}

func (this *KafkaClient) UnmarshalResult(data []byte) (*core.ResultInfo, error) {
	//fmt.Println(msg)
	//保存解码后的数据，Value可以为任意数据类型
	var resultInfo core.ResultInfo
	err := proto.Unmarshal(data, &resultInfo)
	if err != nil {
		return nil, err
	} else {
		return &resultInfo, err
	}
}
func (this *KafkaClient) Unmarshal(data []byte) (*core.RPCInfo, error) {
	//fmt.Println(msg)
	//保存解码后的数据，Value可以为任意数据类型
	var rpcInfo core.RPCInfo
	err := proto.Unmarshal(data, &rpcInfo)
	if err != nil {
		return nil, err
	} else {
		return &rpcInfo, err
	}
}
func (this *KafkaClient) Marshal(rpcInfo *core.RPCInfo) ([]byte, error) {
	b, err := proto.Marshal(rpcInfo)
	return b, err
}
