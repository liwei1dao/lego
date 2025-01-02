package lgrpc

import (
	"context"
	"io"
	"reflect"
	"runtime"
	"sync"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/rpc"
	"github.com/liwei1dao/lego/sys/rpc/protocol"
)

var TypeOfError = reflect.TypeOf((*error)(nil)).Elem()
var TypeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()

func newSys(options *Options) (sys *LGRpc, err error) {
	sys = &LGRpc{
		options: options,
	}

	return
}

type LGRpc struct {
	options  *Options
	metadata string
	cpools   rpc.IConnPool
	server   *Server
	mutex    sync.Mutex
	clients  map[string]*Client
}

// 启动系统
func (this *LGRpc) Start() (err error) {

	return
}

// 关闭系统
func (this *LGRpc) Close() (err error) {

	return
}

func (this *LGRpc) ServiceNode() core.IServiceNode {
	return this.options.ServiceNode
}

// 注册服务
func (this *LGRpc) Register(name string, fn interface{}) (err error) {
	err = this.server.Register(name, fn, this.metadata)
	return
}

// 注销服务
func (this *LGRpc) UnRegister() (err error) {
	err = this.server.UnregisterAll()
	return
}

// 同步执行
func (this *LGRpc) Call(ctx context.Context, service string, req interface{}, reply interface{}) (err error) { //同步调用 等待结果
	if client, ok := this.clients[service]; ok {
		err = client.Call(ctx, req, reply)
	} else {

	}
	return
}

// 异步执行 异步返回
func (this *LGRpc) Go(ctx context.Context, service string, req interface{}, reply interface{}) (call *MessageCall, err error) { //异步调用 异步返回
	if client, ok := this.clients[service]; ok {
		call, err = client.Go(ctx, req, reply)
	} else {

	}
	return
}

func (this *LGRpc) Broadcast(ctx context.Context, service string, req interface{}) (err error) {
	if client, ok := this.clients[service]; ok {
		err = client.Broadcast(ctx, req)
	} else {

	}
	return
}

// 读取远程消息
func (this *LGRpc) ReadProtocol(client rpc.IConnClient, r io.Reader) (err error) {
	req := protocol.GetPooledMsg()
	err = req.Decode(r)
	if err != nil {
		return
	}
	go this.Handle(client, req)
	return
}

// 接收到远程消息
func (this *LGRpc) Handle(client rpc.IConnClient, message rpc.IMessage) {
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 1024)
			buf = buf[:runtime.Stack(buf, true)]
			this.options.Log.Errorf("failed to handle the request: %v,stacks: %s", r, buf)
		}
	}()
	ctx := WithValue(context.Background(), RemoteConnContextKey, client)
	this.options.Log.Debug("[lgrpc Handle] message", log.Field{Key: "Header", Value: message.PrintHeader()}, log.Field{Key: "Service", Value: message.GetService()})
	if message.MessageType() == rpc.Request { //请求消息
		if message.IsHeartbeat() { //心跳
			client.ResetHbeat()
			return
		}
		if message.IsShakeHands() { //握手
			this.ShakehandsResponse(ctx, client, message)
			return
		}
		cancelFunc := parseServerTimeout(ctx, message)
		if cancelFunc != nil {
			defer cancelFunc()
		}
		if res, _ := this.server.handleRequest(ctx, message); res != nil {
			if !message.IsOneway() { //需要回应
				if len(res.Payload()) > 1024 && res.CompressType() != rpc.CompressNone {
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
