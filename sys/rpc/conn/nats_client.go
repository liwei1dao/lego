package conn

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/liwei1dao/lego/sys/rpc/core"
	"github.com/liwei1dao/lego/utils/container"
	"github.com/nats-io/nats.go"
)

func NewNatsClient(sys core.ISys, natsaddr string, pushrpcId, receiveId string) (natsClient *NatsClient, err error) {
	var (
		conn *nats.Conn
		subs *nats.Subscription
	)
	conn, err = nats.Connect(natsaddr)
	if err != nil {
		err = fmt.Errorf("RPC NewNatsClient nats.Connect err:%v", err)
		return
	}
	if subs, err = conn.SubscribeSync(receiveId); err != nil {
		return nil, fmt.Errorf("rpc NatsClient 连接错误 err:%s", err.Error())
	}
	natsClient = &NatsClient{
		sys:               sys,
		rpcId:             pushrpcId,
		callinfos:         container.NewBeeMap(),
		callbackqueueName: receiveId,
		conn:              conn,
		subs:              subs,
	}
	go natsClient.on_request_handle()
	return
}

type NatsClient struct {
	sys               core.ISys
	callinfos         *container.BeeMap
	callbackqueueName string
	rpcId             string
	conn              *nats.Conn
	subs              *nats.Subscription
}

func (this *NatsClient) Stop() (err error) {
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
	err = this.subs.Unsubscribe()
	return
}

func (this *NatsClient) Delete(key string) (err error) {
	this.callinfos.Delete(key)
	return
}

/**
消息请求
*/
func (this *NatsClient) Call(callInfo core.CallInfo, callback chan core.ResultInfo) error {
	//var err error
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
	body, err := this.Marshal(&callInfo.RpcInfo)
	if err != nil {
		return err
	}
	return this.conn.Publish(this.rpcId, body)
}

/**
消息请求 不需要回复
*/
func (this *NatsClient) CallNR(callInfo core.CallInfo) error {
	body, err := this.Marshal(&callInfo.RpcInfo)
	if err != nil {
		return err
	}
	return this.conn.Publish(this.rpcId, body)
}

func (this *NatsClient) on_request_handle() {
	defer func() {
		if r := recover(); r != nil {
			this.sys.Errorf("NatsClient panic:", r)
		}
	}()
locp:
	for {
		m, err := this.subs.NextMsg(time.Minute)
		if err != nil && err == nats.ErrTimeout {
			continue
		} else if err != nil {
			break locp
		}

		resultInfo, err := this.UnmarshalResult(m.Data)
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
	this.sys.Debugf("RPC NatsClient on_request_handle exit")
}

func (this *NatsClient) UnmarshalResult(data []byte) (*core.ResultInfo, error) {
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
func (this *NatsClient) Unmarshal(data []byte) (*core.RPCInfo, error) {
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
func (this *NatsClient) Marshal(rpcInfo *core.RPCInfo) ([]byte, error) {
	b, err := proto.Marshal(rpcInfo)
	return b, err
}
