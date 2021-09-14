package rpc

import (
	fmt "fmt"
	"time"

	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc/core"
	"github.com/liwei1dao/lego/sys/rpc/serialize"
	"github.com/liwei1dao/lego/utils/id"
	"google.golang.org/protobuf/proto"
)

type RPCClient struct {
	ServiceId  string
	rpcExpired time.Duration
	conn       IRPCConnClient
}

func (this *RPCClient) Stop() (err error) {
	err = this.conn.Stop()
	return
}

/**
消息请求 需要回复
*/
func (this *RPCClient) Call(_func string, params ...interface{}) (interface{}, error) {
	var ArgsType []string = make([]string, len(params))
	var args [][]byte = make([][]byte, len(params))

	for k, param := range params {
		var err error = nil
		ArgsType[k], args[k], err = serialize.Serialize(param)
		if err != nil {
			log.Errorf("RPC CallNR ServerId = %s f = %s  ERROR = %v ", this.ServiceId, _func, err)
			return nil, fmt.Errorf("args[%d] error %s", k, err.Error())
		}
	}
	start := time.Now()
	r, errstr := this.CallArgs(_func, ArgsType, args)
	// log.Debugf("RPC Call ServerId = %s f = %s Elapsed = %v Result = %v ERROR = %v", this.options.sId, _func, time.Since(start), r, errstr)
	if errstr != "" {
		log.Errorf("RPC Call ServerId = %s f = %s Elapsed = %v Result = %v ERROR = %v", this.ServiceId, _func, time.Since(start), r, errstr)
		return r, fmt.Errorf(errstr)
	} else {
		return r, nil
	}
}

/**
消息请求 不需要回复
*/
func (this *RPCClient) CallNR(_func string, params ...interface{}) (err error) {
	var ArgsType []string = make([]string, len(params))
	var args [][]byte = make([][]byte, len(params))
	for k, param := range params {
		ArgsType[k], args[k], err = serialize.Serialize(param)
		if err != nil {
			log.Errorf("RPC CallNR ServerId = %s f = %s  ERROR = %v ", this.ServiceId, _func, err)
			return fmt.Errorf("args[%d] error %s", k, err.Error())
		}
	}
	start := time.Now()
	err = this.CallNRArgs(_func, ArgsType, args)
	if err != nil {
		log.Errorf("RPC CallNR ServerId = %s f = %s Elapsed = %v ERROR = %v", this.ServiceId, _func, time.Since(start), err)
	}
	return err
}

func (this *RPCClient) CallArgs(_func string, ArgsType []string, args [][]byte) (interface{}, string) {
	var correlation_id = id.Rand().Hex()
	rpcInfo := &core.RPCInfo{
		Fn:       *proto.String(_func),
		Reply:    *proto.Bool(true),
		Expired:  *proto.Int64((time.Now().UTC().Add(this.rpcExpired).UnixNano()) / 1000000),
		Cid:      *proto.String(correlation_id),
		Args:     args,
		ArgsType: ArgsType,
	}

	callInfo := &core.CallInfo{
		RpcInfo: *rpcInfo,
	}
	callback := make(chan core.ResultInfo, 1)
	var err error
	err = this.conn.Call(*callInfo, callback)
	if err != nil {
		log.Errorf("RPC RPCClient CallArgs err:%v", err)
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
		result, err := serialize.UnSerialize(resultInfo.ResultType, resultInfo.Result)
		if err != nil {
			return nil, err.Error()
		}
		return result, resultInfo.Error
	case <-time.After(this.rpcExpired):
		close(callback)
		this.conn.Delete(rpcInfo.Cid)
		return nil, "timeout"
	}
}

func (this *RPCClient) CallNRArgs(_func string, ArgsType []string, args [][]byte) (err error) {
	var correlation_id = id.Rand().Hex()
	rpcInfo := &core.RPCInfo{
		Fn:       *proto.String(_func),
		Reply:    *proto.Bool(false),
		Expired:  *proto.Int64((time.Now().UTC().Add(this.rpcExpired).UnixNano()) / 1000000),
		Cid:      *proto.String(correlation_id),
		Args:     args,
		ArgsType: ArgsType,
	}
	callInfo := &core.CallInfo{
		RpcInfo: *rpcInfo,
	}
	return this.conn.CallNR(*callInfo)
}
