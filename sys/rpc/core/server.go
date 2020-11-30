package core

import (
	"encoding/json"
	"fmt"
	"reflect"
	"runtime"
	"sync"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
)

func NewRpcServer(opt ...sOption) (srpc *RpcServer, err error) {
	srpc = new(RpcServer)
	srpc.opts = newSOptions(opt...)
	srpc.functions = make(map[string][]*FunctionInfo)
	srpc.ch = make(chan int, srpc.opts.MaxCoroutine)
	srpc.msgserver, err = NewNatsServer(srpc)
	return
}

type RpcServer struct {
	opts      sOptions
	msgserver IRpcMsgServer
	functions map[string][]*FunctionInfo
	listener  RPCListener
	wg        sync.WaitGroup //任务阻塞
	ch        chan int       //控制模块可同时开启的最大协程数
	executing int
}

func (this *RpcServer) RpcId() string {
	return this.msgserver.RpcId()
}

func (this *RpcServer) Done() (err error) {
	this.wg.Wait()
	if this.msgserver != nil {
		err = this.msgserver.Done()
	}
	return
}
func (this *RpcServer) Wait() error {
	this.ch <- 1
	return nil
}
func (this *RpcServer) Finish() {
	<-this.ch
}
func (this *RpcServer) Call(callInfo CallInfo) error {
	this.runFunc(callInfo)
	return nil
}

func (this *RpcServer) GetRpcInfo() (rfs []core.Rpc_Key) {
	rfs = make([]core.Rpc_Key, 0)
	for k, _ := range this.functions {
		rfs = append(rfs, core.Rpc_Key(k))
	}
	return
}

func (this *RpcServer) Register(id string, f interface{}) {
	if _, ok := this.functions[id]; !ok {
		this.functions[id] = []*FunctionInfo{}
	}
	this.functions[id] = append(this.functions[id], &FunctionInfo{
		Function:  reflect.ValueOf(f),
		Goroutine: false,
	})
}
func (this *RpcServer) RegisterGO(id string, f interface{}) {
	if _, ok := this.functions[id]; !ok {
		this.functions[id] = []*FunctionInfo{}
	}
	this.functions[id] = append(this.functions[id], &FunctionInfo{
		Function:  reflect.ValueOf(f),
		Goroutine: true,
	})
}

func (this *RpcServer) UnRegister(id string, f interface{}) {
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

func (this *RpcServer) runFunc(callInfo CallInfo) {
	start := time.Now()
	_errorCallback := func(Cid string, Error string) {
		resultInfo := NewResultInfo(Cid, Error, NULL, nil)
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

	funcs, ok := this.functions[callInfo.RpcInfo.Fn]
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
			this.executing++
			defer func() {
				this.wg.Add(-1)
				this.executing--
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
					allError := fmt.Sprintf("%s rpc func(%s) error %s ----Stack----%s", this.opts.sId, callInfo.RpcInfo.Fn, rn, errstr)
					log.Errorf(allError)
					_errorCallback(callInfo.RpcInfo.Cid, allError)
				}
			}()
			var in []reflect.Value
			if len(ArgsType) > 0 {
				in = make([]reflect.Value, len(params))
				for k, v := range ArgsType {
					ty, err := UnSerialize(v, params[k])
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
			argsType, args, err := Serialize(rs[0])
			if err != nil {
				_errorCallback(callInfo.RpcInfo.Cid, err.Error())
				return
			}
			resultInfo := NewResultInfo(
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
func (this *RpcServer) doCallback(callInfo CallInfo) {
	if callInfo.RpcInfo.Reply {
		//需要回复的才回复
		err := callInfo.Agent.(MQServer).Callback(callInfo)
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
