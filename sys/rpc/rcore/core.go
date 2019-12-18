package rcore

import "reflect"

type FunctionInfo struct {
	Function  reflect.Value
	Goroutine bool
}

type CallInfo struct {
	RpcInfo RPCInfo
	Result  ResultInfo
	Props   map[string]interface{}
	Agent   MQServer //代理者  AMQPServer / LocalServer 都继承 Callback(callinfo CallInfo)(error) 方法
}

type MQServer interface {
	Callback(callinfo CallInfo) error
}

type ClinetCallInfo struct {
	correlation_id string
	timeout        int64 //超时
	call           chan ResultInfo
}

type RPCListener interface {
	NoFoundFunction(fn string) (*FunctionInfo, error)
	BeforeHandle(fn string, callInfo *CallInfo) error
	OnTimeOut(fn string, Expired int64)
	OnError(fn string, callInfo *CallInfo, err error)
	OnComplete(fn string, callInfo *CallInfo, result *ResultInfo, exec_time int64)
}

type IRpcMsgServer interface {
	RpcId() string
	Done() (err error)
}

type IRpcMsgClient interface {
	Done() (err error)
	Delete(key string) (err error)
	Call(callInfo CallInfo, callback chan ResultInfo) error
	CallNR(callInfo CallInfo) error
}
