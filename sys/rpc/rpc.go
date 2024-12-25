package rpc

import (
	"context"
	"reflect"
	"runtime"
	"sync"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc/protocol"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

var TypeOfError = reflect.TypeOf((*error)(nil)).Elem()
var TypeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()

func newSys(options *Options) (sys *rpc, err error) {
	sys = &rpc{
		options: options,
	}

	return
}

type rpc struct {
	options  *Options
	metadata string
	server   *Server
	mutex    sync.Mutex
	clients  map[string]*Client
}

// 启动系统
func (this *rpc) Start() (err error) {

	return
}

// 关闭系统
func (this *rpc) Close() (err error) {

	return
}

func (this *rpc) ServiceNode() core.IServiceNode {
	return this.options.ServiceNode
}

// 注册服务
func (this *rpc) Register(name string, fn interface{}) (err error) {
	err = this.server.Register(name, fn, this.metadata)
	return
}

// 注销服务
func (this *rpc) UnRegister(name string) (err error) {
	err = this.server.UnregisterAll()
	return
}

// 同步执行
func (this *rpc) Call(ctx context.Context, service string, req interface{}, reply interface{}) (err error) { //同步调用 等待结果
	if client, ok := this.clients[service]; ok {
		err = client.Call(ctx, req, reply)
	} else {

	}
	return
}

// 异步执行 异步返回
func (this *rpc) Go(ctx context.Context, service string, req interface{}, reply interface{}) (call *MessageCall, err error) { //异步调用 异步返回
	if client, ok := this.clients[service]; ok {
		call, err = client.Go(ctx, req, reply)
	} else {

	}
	return
}

// 接收到远程消息
func (this *rpc) Handle(client rpccore.IConnClient, message rpccore.IMessage) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			buf = buf[:runtime.Stack(buf, true)]
			this.options.Log.Errorf("failed to handle the request: %v， stacks: %s", r, buf)
		}
	}()
	ctx := rpccore.WithValue(context.Background(), rpccore.RemoteConnContextKey, client)
	// this.options.Log.Debug("[handle] message", log.Field{Key: "Header", Value: message.PrintHeader()}, log.Field{Key: "ServiceMethod", Value: message.ServiceMethod()})
	if message.MessageType() == rpccore.Request { //请求消息
		if message.IsHeartbeat() { //心跳
			client.ResetHbeat()
			return
		}
		if message.IsShakeHands() {
			this.ShakehandsResponse(ctx, client, message)
			return
		}
		cancelFunc := parseServerTimeout(ctx, message)
		if cancelFunc != nil {
			defer cancelFunc()
		}
		if res, _ := this.server.handleRequest(ctx, message); res != nil {
			if !message.IsOneway() { //需要回应
				if len(res.Payload()) > 1024 && res.CompressType() != rpccore.CompressNone {
					res.SetCompressType(res.CompressType())
				}
				data := res.EncodeSlicePointer()
				client.Write(*data)
				protocol.PutData(data)
			}
		}
	} else { //回应
		if client, ok := this.clients[message.GetService()]; ok {
			client.handleresponse(ctx, message)
		} else {
			this.options.Log.Errorf("no found client : %s", message.GetService())
		}
	}
}

// 执行远程服务---------------------------------------------------------------------------------------
