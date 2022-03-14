package core

import (
	"reflect"
)

type (
	IRpcServer interface {
		Call(callInfo CallInfo) error
	}
	IRPCConnServer interface {
		Callback(callinfo CallInfo) error
	}
	FunctionInfo struct {
		Function  reflect.Value
		Goroutine bool
	}
	CallInfo struct {
		RpcInfo RPCInfo
		Result  ResultInfo
		Props   map[string]interface{}
		Agent   IRPCConnServer
	}
	ClinetCallInfo struct {
		Correlation_id string
		Timeout        int64 //超时
		Call           chan ResultInfo
	}
)
