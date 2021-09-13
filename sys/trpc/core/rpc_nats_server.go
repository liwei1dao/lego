package core

import (
	"fmt"
	"runtime"
	"time"

	"github.com/liwei1dao/lego/sys/log"

	"github.com/golang/protobuf/proto"
	"github.com/nats-io/nats.go"
)

func NewNatsServer(server *RpcServer) (natsServer *NatsServer, err error) {
	natsServer = new(NatsServer)
	natsServer.server = server
	natsServer.nats = server.opts.Nats
	natsServer.rpcId = nats.NewInbox()
	natsServer.done = make(chan error)
	if subs, err := natsServer.nats.SubscribeSync(natsServer.rpcId); err != nil {
		return nil, fmt.Errorf("rpc NatsServer 连接错误 err:%s", err.Error())
	} else {
		natsServer.subs = subs
		go natsServer.on_request_handle()
		return natsServer, nil
	}
}

type NatsServer struct {
	nats   *nats.Conn
	subs   *nats.Subscription
	rpcId  string
	done   chan error
	server *RpcServer
}

func (this *NatsServer) RpcId() string {
	return this.rpcId
}

/**
注销消息队列
*/
func (this *NatsServer) Done() (err error) {
	this.done <- nil
	return
}

func (this *NatsServer) Callback(callinfo CallInfo) error {
	body, _ := this.MarshalResult(callinfo.Result)
	reply_to := callinfo.Props["reply_to"].(string)
	return this.nats.Publish(reply_to, body)
}

/**
接收请求信息
*/
func (this *NatsServer) on_request_handle() error {
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
			log.Errorf("%s\n ----Stack----\n%s", rn, errstr)
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
		rpcInfo, err := this.Unmarshal(m.Data)
		if err == nil {
			callInfo := &CallInfo{
				RpcInfo: *rpcInfo,
			}
			callInfo.Props = map[string]interface{}{
				"reply_to": rpcInfo.ReplyTo,
			}
			callInfo.Agent = this //设置代理为NatsServer
			this.server.Call(*callInfo)
		} else {
			fmt.Println("error ", err)
		}
	}
}

func (s *NatsServer) Unmarshal(data []byte) (*RPCInfo, error) {
	var rpcInfo RPCInfo
	err := proto.Unmarshal(data, &rpcInfo)
	if err != nil {
		return nil, err
	} else {
		return &rpcInfo, err
	}
}

// goroutine safe
func (s *NatsServer) MarshalResult(resultInfo ResultInfo) ([]byte, error) {
	b, err := proto.Marshal(&resultInfo)
	return b, err
}
