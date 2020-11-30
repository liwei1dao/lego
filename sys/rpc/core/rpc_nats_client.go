package core

import (
	"fmt"
	"github.com/liwei1dao/lego/sys/log"
	cont "github.com/liwei1dao/lego/utils/concurrent"
	"runtime"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

func NewNatsClient(rId string, nts *nats.Conn) (natsClient *NatsClient, err error) {
	natsClient = new(NatsClient)
	natsClient.callinfos = cont.NewBeeMap()
	natsClient.rpcId = rId
	natsClient.nats = nts
	natsClient.callbackqueueName = nats.NewInbox()
	if subs, err := natsClient.nats.SubscribeSync(natsClient.callbackqueueName); err != nil {
		return nil, fmt.Errorf("rpc NatsClient 连接错误 err:%s", err.Error())
	} else {
		natsClient.subs = subs
		go natsClient.on_request_handle()
		return natsClient, nil
	}
}

type NatsClient struct {
	callinfos         *cont.BeeMap
	callbackqueueName string
	rpcId             string
	nats              *nats.Conn
	subs              *nats.Subscription
	done              chan error
}

func (c *NatsClient) Done() (err error) {
	if c.callinfos != nil {
		//清理 callinfos 列表
		for key, clinetCallInfo := range c.callinfos.Items() {
			if clinetCallInfo != nil {
				//关闭管道
				close(clinetCallInfo.(ClinetCallInfo).call)
				//从Map中删除
				c.callinfos.Delete(key)
			}
		}
		c.callinfos = nil
	} else {
		log.Errorf("rpc NatsClient callinfos 异常为空")
	}
	c.done <- nil
	return
}

func (c *NatsClient) Delete(key string) (err error) {
	c.callinfos.Delete(key)
	return
}

/**
消息请求
*/
func (this *NatsClient) Call(callInfo CallInfo, callback chan ResultInfo) error {
	//var err error
	if this.callinfos == nil {
		return fmt.Errorf("AMQPClient is closed")
	}
	callInfo.RpcInfo.ReplyTo = this.callbackqueueName
	var correlation_id = callInfo.RpcInfo.Cid

	clinetCallInfo := &ClinetCallInfo{
		correlation_id: correlation_id,
		call:           callback,
		timeout:        callInfo.RpcInfo.Expired,
	}
	this.callinfos.Set(correlation_id, *clinetCallInfo)
	body, err := this.Marshal(&callInfo.RpcInfo)
	if err != nil {
		return err
	}
	return this.nats.Publish(this.rpcId, body)
}

/**
消息请求 不需要回复
*/
func (this *NatsClient) CallNR(callInfo CallInfo) error {
	body, err := this.Marshal(&callInfo.RpcInfo)
	if err != nil {
		return err
	}
	return this.nats.Publish(this.rpcId, body)
}
func (this *NatsClient) on_request_handle() error {
	defer func() {
		if r := recover(); r != nil {
			var rn = ""
			switch r.(type) {

			case string:
				rn = r.(string)
			case error:
				rn = r.(error).Error()
			}
			buf := make([]byte, 1024)
			l := runtime.Stack(buf, false)
			errstr := string(buf[:l])
			log.Panicf("%s\n ----Stack----\n%s", rn, errstr)
		}
	}()

	go func() {
		<-this.done
		this.subs.Unsubscribe()
	}()

	for {
		m, err := this.subs.NextMsg(time.Minute)
		if err != nil && err == nats.ErrTimeout {
			continue
		} else if err != nil {
			return err
		}

		resultInfo, err := this.UnmarshalResult(m.Data)
		if err != nil {
			log.Errorf("Unmarshal faild", err)
		} else {
			correlation_id := resultInfo.Cid
			clinetCallInfo := this.callinfos.Get(correlation_id)
			//删除
			this.callinfos.Delete(correlation_id)
			if clinetCallInfo != nil {
				clinetCallInfo.(ClinetCallInfo).call <- *resultInfo
				close(clinetCallInfo.(ClinetCallInfo).call)
			} else {
				//可能客户端已超时了，但服务端处理完还给回调了
				log.Warnf("rpc callback no found : [%s]", correlation_id)
			}
		}
	}
	return nil
}
func (this *NatsClient) UnmarshalResult(data []byte) (*ResultInfo, error) {
	//fmt.Println(msg)
	//保存解码后的数据，Value可以为任意数据类型
	var resultInfo ResultInfo
	err := proto.Unmarshal(data, &resultInfo)
	if err != nil {
		return nil, err
	} else {
		return &resultInfo, err
	}
}
func (this *NatsClient) Unmarshal(data []byte) (*RPCInfo, error) {
	//fmt.Println(msg)
	//保存解码后的数据，Value可以为任意数据类型
	var rpcInfo RPCInfo
	err := proto.Unmarshal(data, &rpcInfo)
	if err != nil {
		return nil, err
	} else {
		return &rpcInfo, err
	}
	panic("bug")
}
func (this *NatsClient) Marshal(rpcInfo *RPCInfo) ([]byte, error) {
	b, err := proto.Marshal(rpcInfo)
	return b, err
}
