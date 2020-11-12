package rcore

import (
	"fmt"
	"time"

	"github.com/liwei1dao/lego/sys/log"
	rpcserialize "github.com/liwei1dao/lego/sys/rpc/serialize"
	"github.com/liwei1dao/lego/utils/id"

	"github.com/golang/protobuf/proto"
)

func NewRpcClient(opt ...cOption) (srpc *RpcClient, err error) {
	srpc = new(RpcClient)
	srpc.opts = newCOptions(opt...)
	srpc.msgclient, err = NewNatsClient(srpc.opts.rpcId, srpc.opts.Nats)
	return
}

type RpcClient struct {
	opts cOptions

	msgclient IRpcMsgClient
}

func (this *RpcClient) Done() (err error) {
	if this.msgclient != nil {
		err = this.msgclient.Done()
	}
	return
}

/**
消息请求 需要回复
*/
func (this *RpcClient) Call(_func string, params ...interface{}) (interface{}, error) {
	var ArgsType []string = make([]string, len(params))
	var args [][]byte = make([][]byte, len(params))

	for k, param := range params {
		var err error = nil
		ArgsType[k], args[k], err = rpcserialize.Serialize(param)
		if err != nil {
			log.Errorf("RPC CallNR ServerId = %s f = %s  ERROR = %v ", this.opts.sId, _func, err)
			return nil, fmt.Errorf("args[%d] error %s", k, err.Error())
		}
	}
	start := time.Now()
	r, errstr := this.CallArgs(_func, ArgsType, args)
	// log.Debugf("RPC Call ServerId = %s f = %s Elapsed = %v Result = %v ERROR = %v", this.opts.sId, _func, time.Since(start), r, errstr)
	if errstr != "" {
		log.Errorf("RPC Call ServerId = %s f = %s Elapsed = %v Result = %v ERROR = %v", this.opts.sId, _func, time.Since(start), r, errstr)
		return r, fmt.Errorf(errstr)
	} else {
		return r, nil
	}
}

/**
消息请求 不需要回复
*/
func (this *RpcClient) CallNR(_func string, params ...interface{}) (err error) {
	var ArgsType []string = make([]string, len(params))
	var args [][]byte = make([][]byte, len(params))
	for k, param := range params {
		ArgsType[k], args[k], err = rpcserialize.Serialize(param)
		if err != nil {
			log.Errorf("RPC CallNR ServerId = %s f = %s  ERROR = %v ", this.opts.sId, _func, err)
			return fmt.Errorf("args[%d] error %s", k, err.Error())
		}
	}
	start := time.Now()
	err = this.CallNRArgs(_func, ArgsType, args)
	if err != nil {
		log.Errorf("RPC CallNR ServerId = %s f = %s Elapsed = %v ERROR = %v", this.opts.sId, _func, time.Since(start), err)
	}
	return err
}

func (this *RpcClient) CallArgs(_func string, ArgsType []string, args [][]byte) (interface{}, string) {
	var correlation_id = id.Rand().Hex()
	rpcInfo := &RPCInfo{
		Fn:       *proto.String(_func),
		Reply:    *proto.Bool(true),
		Expired:  *proto.Int64((time.Now().UTC().Add(time.Second * time.Duration(this.opts.RpcExpired)).UnixNano()) / 1000000),
		Cid:      *proto.String(correlation_id),
		Args:     args,
		ArgsType: ArgsType,
	}

	callInfo := &CallInfo{
		RpcInfo: *rpcInfo,
	}
	callback := make(chan ResultInfo, 1)
	var err error
	err = this.msgclient.Call(*callInfo, callback)
	if err != nil {
		return nil, err.Error()
	}
	select {
	case resultInfo, ok := <-callback:
		if !ok {
			return nil, "client closed"
		}
		if resultInfo.Error != "" {
			return nil, resultInfo.Error
		}
		result, err := rpcserialize.UnSerialize(resultInfo.ResultType, resultInfo.Result)
		if err != nil {
			return nil, err.Error()
		}
		return result, resultInfo.Error
	case <-time.After(time.Second * time.Duration(this.opts.RpcExpired)):
		close(callback)
		this.msgclient.Delete(rpcInfo.Cid)
		return nil, "deadline exceeded"
	}
}

func (this *RpcClient) CallNRArgs(_func string, ArgsType []string, args [][]byte) (err error) {
	var correlation_id = id.Rand().Hex()
	rpcInfo := &RPCInfo{
		Fn:       *proto.String(_func),
		Reply:    *proto.Bool(false),
		Expired:  *proto.Int64((time.Now().UTC().Add(time.Second * time.Duration(this.opts.RpcExpired)).UnixNano()) / 1000000),
		Cid:      *proto.String(correlation_id),
		Args:     args,
		ArgsType: ArgsType,
	}
	callInfo := &CallInfo{
		RpcInfo: *rpcInfo,
	}
	return this.msgclient.CallNR(*callInfo)
}
