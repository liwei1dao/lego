package lgrpc

import (
	"context"
	"errors"
)

const (
	ServiceError         = "__rpcx_error__"   //服务错误信息字段
	ServerTimeout        = "__ServerTimeout"  //服务超时字段
	ReqMetaDataKey       = "__req_metadata"   //请求元数据字段
	ResMetaDataKey       = "__res_metadata"   //返回元数据字段
	ServiceAddrKey       = "__service_addr__" //服务端地址
	CallSeqKey           = "__call_seq__"     //客户端请求id存储key
	RemoteConnContextKey = "__remote_conn__"  //远程连接上下文
)

var (
	ErrServerClosed          = errors.New("http: Server closed")                                   //服务关闭
	ErrMetaKVMissing         = errors.New("wrong metadata lines. some keys or values are missing") //解析Meta对象错误
	ErrUnsupportedCompressor = errors.New("unsupported compressor")                                //解压缩错误
	ErrXClientNoServer       = errors.New("can not found any server")
	ErrUnsupportedCodec      = errors.New("unsupported codec")
)

type (
	ISys interface {
		Start() (err error)
		Close() (err error)
		Register(name string, fn interface{}) (err error)
		UnRegister() (err error)
		Call(ctx context.Context, service string, args interface{}, reply interface{}) (err error)                  //同步调用 等待结果
		Go(ctx context.Context, service string, args interface{}, reply interface{}) (call *MessageCall, err error) //异步调用 异步返回
		Broadcast(ctx context.Context, service string, args interface{}) (err error)
	}

	IClient interface {
	}

	ICodec interface {
		Marshal(v interface{}) ([]byte, error)
		Unmarshal(data []byte, v interface{}) error
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, opt ...Option) (err error) {
	var option *Options
	if option, err = newOptions(config, opt...); err != nil {
		return
	}
	defsys, err = newSys(option)
	return
}

func NewSys(opt ...Option) (sys ISys, err error) {
	var option *Options
	if option, err = newOptionsByOption(opt...); err != nil {
		return
	}
	sys, err = newSys(option)
	return
}

func Start() (err error) {
	return defsys.Start()
}
func Close() (err error) {
	return defsys.Close()
}

func Register(name string, fn interface{}) error {
	return defsys.Register(name, fn)
}
func UnRegister() {
	defsys.UnRegister()
}
func Call(ctx context.Context, service string, args interface{}, reply interface{}) (err error) {
	return defsys.Call(ctx, service, args, reply)
}
func Go(ctx context.Context, service string, args interface{}, reply interface{}) (call *MessageCall, err error) {
	return defsys.Go(ctx, service, args, reply)
}
func Broadcast(ctx context.Context, service string, args interface{}) (err error) {
	return defsys.Broadcast(ctx, service, args)
}
