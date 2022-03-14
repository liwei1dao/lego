package conn

import (
	"fmt"
	"time"

	"github.com/golang/protobuf/proto"
	"github.com/liwei1dao/lego"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc/core"
	"github.com/nats-io/nats.go"
)

func NewNatsService(natsaddr string, rpcId string) (natsService *NatsService, err error) {
	var (
		conn *nats.Conn
		subs *nats.Subscription
	)
	conn, err = nats.Connect(natsaddr)
	if err != nil {
		err = fmt.Errorf("RPC NewNatsService nats.Connect err:%v", err)
		return
	}
	if subs, err = conn.SubscribeSync(rpcId); err != nil {
		err = fmt.Errorf("RPC NatsServer conn.SubscribeSync err:%v", err)
		return
	}
	natsService = &NatsService{
		conn: conn,
		subs: subs,
	}
	return
}

type NatsService struct {
	conn    *nats.Conn
	subs    *nats.Subscription
	service core.IRpcServer
}

func (this *NatsService) Callback(callinfo core.CallInfo) error {
	body, _ := this.MarshalResult(callinfo.Result)
	reply_to := callinfo.Props["reply_to"].(string)
	return this.conn.Publish(reply_to, body)
}

func (this *NatsService) Start(service core.IRpcServer) (err error) {
	this.service = service
	go this.on_request_handle()
	return
}

func (this *NatsService) Stop() (err error) {
	err = this.subs.Unsubscribe()
	return
}
func (this *NatsService) on_request_handle() {
	defer lego.Recover("RPC NatsService")
locp:
	for {
		m, err := this.subs.NextMsg(time.Minute)
		if err != nil && err == nats.ErrTimeout {
			continue
		} else if err != nil {
			break locp
		}
		rpcInfo, err := this.Unmarshal(m.Data)
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
	log.Debugf("RPC NatsService on_request_handle exit")
}
func (s *NatsService) Unmarshal(data []byte) (*core.RPCInfo, error) {
	var rpcInfo core.RPCInfo
	err := proto.Unmarshal(data, &rpcInfo)
	if err != nil {
		return nil, err
	} else {
		return &rpcInfo, err
	}
}

func (s *NatsService) MarshalResult(resultInfo core.ResultInfo) ([]byte, error) {
	b, err := proto.Marshal(&resultInfo)
	return b, err
}
