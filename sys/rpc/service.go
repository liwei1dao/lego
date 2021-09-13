package rpc

import (
	"encoding/json"
	fmt "fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	lgcore "github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc/core"
	"github.com/liwei1dao/lego/sys/rpc/serialize"
)

type RPCService struct {
	rpcId        string
	maxCoroutine int
	wg           sync.WaitGroup //任务阻塞
	runch        chan int       //控制模块可同时开启的最大协程数
	listener     RPCListener    //监控输出
	conn         IRPCConnServer
	functions    map[lgcore.Rpc_Key][]*core.FunctionInfo
}

func (this *RPCService) RpcId() string {
	return this.rpcId
}

func (this *RPCService) Start() (err error) {
	err = this.conn.Start(this)
	return
}

func (this *RPCService) Stop() (err error) {
	this.wg.Wait()
	err = this.conn.Stop()
	return
}
func (this *RPCService) GetRpcInfo() (rfs []lgcore.Rpc_Key) {
	rfs = make([]lgcore.Rpc_Key, 0)
	for k, _ := range this.functions {
		rfs = append(rfs, lgcore.Rpc_Key(k))
	}
	return
}

func (this *RPCService) Register(id lgcore.Rpc_Key, f interface{}) {
	if _, ok := this.functions[id]; !ok {
		this.functions[id] = []*core.FunctionInfo{}
	}
	this.functions[id] = append(this.functions[id], &core.FunctionInfo{
		Function:  reflect.ValueOf(f),
		Goroutine: false,
	})
}
func (this *RPCService) RegisterGO(id lgcore.Rpc_Key, f interface{}) {
	if _, ok := this.functions[id]; !ok {
		this.functions[id] = []*core.FunctionInfo{}
	}
	this.functions[id] = append(this.functions[id], &core.FunctionInfo{
		Function:  reflect.ValueOf(f),
		Goroutine: true,
	})
}

func (this *RPCService) UnRegister(id lgcore.Rpc_Key, f interface{}) {
	if fs, ok := this.functions[id]; ok {
		for i, v := range fs {
			if v.Function == reflect.ValueOf(f) {
				this.functions[id] = append(fs[0:i], fs[i+1:]...)
			}
		}
		if len(this.functions[id]) == 0 {
			delete(this.functions, id)
		}
	}
}

func (this *RPCService) Wait() {
	if len(this.runch) >= this.maxCoroutine {
		log.Warnf("RPCService Run Coroutine Reach the Upper limit")
	}
	this.runch <- 1
	return
}
func (this *RPCService) Finish() {
	<-this.runch
}

func (this *RPCService) Call(callInfo core.CallInfo) error {
	this.runFunc(callInfo)
	return nil
}

func (this *RPCService) runFunc(callInfo core.CallInfo) {
	start := time.Now()
	_errorCallback := func(Cid string, Error string) {
		resultInfo := core.NewResultInfo(Cid, Error, serialize.NULL, nil)
		callInfo.Result = *resultInfo
		this.doCallback(callInfo)
		if this.listener != nil {
			this.listener.OnError(callInfo.RpcInfo.Fn, &callInfo, fmt.Errorf(Error))
		}
	}
	defer func() {
		if r := recover(); r != nil {
			var rn = ""
			switch r.(type) {

			case string:
				rn = r.(string)
			case error:
				rn = r.(error).Error()
			}
			log.Errorf("recover", rn)
			_errorCallback(callInfo.RpcInfo.Cid, rn)
		}
	}()

	funcs, ok := this.functions[lgcore.Rpc_Key(callInfo.RpcInfo.Fn)]
	if !ok {
		_errorCallback(callInfo.RpcInfo.Cid, fmt.Sprintf("Rpc runFunc NoFound Fuc【%s】", callInfo.RpcInfo.Fn))
		return
	}
	if len(funcs) > 1 && callInfo.RpcInfo.Reply {
		_errorCallback(callInfo.RpcInfo.Cid, fmt.Sprintf("Rpc runFunc Fuc【%s】 为订阅型消息 不能返回信息 订阅数量:%d", callInfo.RpcInfo.Fn, len(funcs)))
		return
	}

	for _, v := range funcs {
		f := v.Function
		params := callInfo.RpcInfo.Args
		ArgsType := callInfo.RpcInfo.ArgsType
		if len(params) != f.Type().NumIn() {
			_errorCallback(callInfo.RpcInfo.Cid, fmt.Sprintf("Rpc runFunc Fuc【%s】 请求参数异常  params %v is not adapted.%v", callInfo.RpcInfo.Fn, params, f.String()))
			return
		}
		_runFunc := func(f reflect.Value, ArgsType []string, params [][]byte) {
			this.wg.Add(1)
			defer func() {
				this.wg.Add(-1)
				this.Finish()
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
					allError := fmt.Sprintf("rpc func(%s) error %s ----Stack----%s", callInfo.RpcInfo.Fn, rn, errstr)
					log.Errorf(allError)
					_errorCallback(callInfo.RpcInfo.Cid, allError)
				}
			}()
			var in []reflect.Value
			if len(ArgsType) > 0 {
				in = make([]reflect.Value, len(params))
				for k, v := range ArgsType {
					ty, err := serialize.UnSerialize(v, params[k])
					if err != nil {
						_errorCallback(callInfo.RpcInfo.Cid, err.Error())
						return
					}
					switch v2 := ty.(type) { //多选语句switch
					case []uint8:
						if reflect.TypeOf(ty).AssignableTo(f.Type().In(k)) {
							in[k] = reflect.ValueOf(ty)
						} else {
							elemp := reflect.New(f.Type().In(k))
							err := json.Unmarshal(v2, elemp.Interface())
							if err != nil {
								log.Errorf("%v []uint8--> %v error with='%v'", callInfo.RpcInfo.Fn, f.Type().In(k), err)
								in[k] = reflect.ValueOf(ty)
							} else {
								in[k] = elemp.Elem()
							}
						}
					case nil:
						in[k] = reflect.Zero(f.Type().In(k))
					default:
						in[k] = reflect.ValueOf(ty)
					}

				}
			}

			if this.listener != nil {
				errs := this.listener.BeforeHandle(callInfo.RpcInfo.Fn, &callInfo)
				if errs != nil {
					_errorCallback(callInfo.RpcInfo.Cid, errs.Error())
					return
				}
			}

			out := f.Call(in)
			var rs []interface{}
			if len(out) != 2 {
				_errorCallback(callInfo.RpcInfo.Cid, "The number of prepare is not adapted.")
				return
			}
			if len(out) > 0 { //prepare out paras
				rs = make([]interface{}, len(out), len(out))
				for i, v := range out {
					rs[i] = v.Interface()
				}
			}
			argsType, args, err := serialize.Serialize(rs[0])
			if err != nil {
				_errorCallback(callInfo.RpcInfo.Cid, err.Error())
				return
			}
			resultInfo := core.NewResultInfo(
				callInfo.RpcInfo.Cid,
				rs[1].(string),
				argsType,
				args,
			)
			callInfo.Result = *resultInfo
			this.doCallback(callInfo)
			// log.Debugf("RPC Exec ModuleType = %v f = %v Elapsed = %v", this.opts.sId, callInfo.RpcInfo.Fn, time.Since(start))
			if this.listener != nil {
				this.listener.OnComplete(callInfo.RpcInfo.Fn, &callInfo, resultInfo, time.Since(start).Nanoseconds())
			}
		}

		this.Wait()
		if v.Goroutine {
			go _runFunc(f, ArgsType, params)
		} else {
			_runFunc(f, ArgsType, params)
		}
	}
}

func (this *RPCService) doCallback(callInfo core.CallInfo) {
	if callInfo.RpcInfo.Reply {
		//需要回复的才回复
		err := callInfo.Agent.Callback(callInfo)
		if err != nil {
			log.Warnf("rpc callback erro :%s", err.Error())
		}
	} else {
		//对于不需要回复的消息,可以判断一下是否出现错误，打印一些警告
		if callInfo.Result.Error != "" {
			log.Warnf("rpc callback erro :%s", callInfo.Result.Error)
		}
	}
}
