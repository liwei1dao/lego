package rpc

import (
	"context"
	"reflect"
	"strconv"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc/protocol"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
)

//获取心跳消息包
func getHeartbeat(node *core.ServiceNode) []byte {
	req := protocol.GetPooledMsg()
	req.SetMessageType(rpccore.Request)
	req.SetHeartbeat(true)
	req.SetOneway(true)
	req.SetServiceMethod("")
	codec := codecs[rpccore.ProtoBuffer]
	data, _ := codec.Marshal(node)
	req.SetPayload(data)
	allData := req.EncodeSlicePointer()
	defer func() {
		protocol.PutData(allData)
		protocol.FreeMsg(req)
	}()
	return *allData
}

//获取握手消息包
func getShakehands(node *core.ServiceNode) []byte {
	req := protocol.GetPooledMsg()
	req.SetMessageType(rpccore.Request)
	req.SetHeartbeat(true)
	req.SetOneway(true)
	req.SetServiceMethod("")
	codec := codecs[rpccore.ProtoBuffer]
	data, _ := codec.Marshal(node)
	req.SetPayload(data)
	allData := req.EncodeSlicePointer()
	defer func() {
		protocol.PutData(allData)
		protocol.FreeMsg(req)
	}()
	return *allData
}

//解析rpc请求超时设置
func parseServerTimeout(ctx *rpccore.Context, req rpccore.IMessage) context.CancelFunc {
	if req == nil || req.Metadata() == nil {
		return nil
	}

	st := req.Metadata()[rpccore.ServerTimeout]
	if st == "" {
		return nil
	}

	timeout, err := strconv.ParseInt(st, 10, 64)
	if err != nil {
		return nil
	}

	newCtx, cancel := context.WithTimeout(ctx.Context, time.Duration(timeout)*time.Millisecond)
	ctx.Context = newCtx
	return cancel
}

//设置消息错误信息
func handleError(res rpccore.IMessage, err error) (rpccore.IMessage, error) {
	res.SetMessageStatusType(rpccore.Error)
	if res.Metadata() == nil {
		res.SetMetadata(make(map[string]string))
	}
	res.Metadata()[rpccore.ServiceError] = err.Error()
	return res, err
}

//是否是内置类型
func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	return isExported(t.Name()) || t.PkgPath() == ""
}

func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}
