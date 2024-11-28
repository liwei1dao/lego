package rpc

import (
	"context"
	"reflect"
	"runtime"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/discovery"
	"github.com/liwei1dao/lego/sys/rpc/connpool"
	"github.com/liwei1dao/lego/sys/rpc/protocol"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

var TypeOfError = reflect.TypeOf((*error)(nil)).Elem()
var TypeOfContext = reflect.TypeOf((*context.Context)(nil)).Elem()

func newSys(options *Options) (sys *rpc, err error) {
	sys = &rpc{
		options: options,
	}

	if sys.cpool, err = connpool.NewConnPool(sys, options.Log, &rpccore.Config{
		ConnectType: options.ConnectType,
		Endpoints:   options.CommAddrs,
	}); err != nil {
		return
	}

	return
}

type rpc struct {
	options *Options
	cpool   rpccore.IConnPool
	service *Service
}

// 启动系统
func (this *rpc) Start() (err error) {
	if err = this.cpool.Start(); err != nil {
		return
	}
	if err = discovery.Start(); err != nil {
		return
	}
	return
}

// 关闭系统
func (this *rpc) Close() (err error) {
	if err = discovery.Stop(); err != nil {
		return
	}
	if err = this.cpool.Close(); err != nil {
		return
	}
	return
}

func (this *rpc) ServiceNode() core.IServiceNode {
	return this.options.ServiceNode
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
			this.service.ShakehandsResponse(ctx, client, message)
			return
		}
		cancelFunc := parseServerTimeout(ctx, message)
		if cancelFunc != nil {
			defer cancelFunc()
		}
		if res, _ := this.service.handleRequest(ctx, message); res != nil {
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
		this.handleresponse(rpccore.NewContext(context.Background()), message)
	}
}
