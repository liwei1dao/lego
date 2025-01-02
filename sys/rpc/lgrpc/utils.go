package lgrpc

import (
	"context"
	"reflect"
	"strconv"
	"time"
	"unicode"
	"unicode/utf8"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/rpc"
	"github.com/liwei1dao/lego/sys/rpc/protocol"
)

// 获取心跳消息包
func getHeartbeat(node *core.ServiceNode) []byte {
	req := protocol.GetPooledMsg()
	req.SetMessageType(rpc.Request)
	req.SetHeartbeat(true)
	req.SetOneway(true)
	req.SetService("")
	codec := codecs[rpc.ProtoBuffer]
	data, _ := codec.Marshal(node)
	req.SetPayload(data)
	allData := req.EncodeSlicePointer()
	defer func() {
		protocol.PutData(allData)
		protocol.FreeMsg(req)
	}()
	return *allData
}

// 获取握手消息包
func getShakehands(node *core.ServiceNode) []byte {
	req := protocol.GetPooledMsg()
	req.SetMessageType(rpc.Request)
	req.SetHeartbeat(true)
	req.SetOneway(true)
	req.SetService("")
	codec := codecs[rpc.ProtoBuffer]
	data, _ := codec.Marshal(node)
	req.SetPayload(data)
	allData := req.EncodeSlicePointer()
	defer func() {
		protocol.PutData(allData)
		protocol.FreeMsg(req)
	}()
	return *allData
}

// 解析rpc请求超时设置
func parseServerTimeout(ctx *Context, req rpc.IMessage) context.CancelFunc {
	if req == nil || req.Metadata() == nil {
		return nil
	}

	st := req.Metadata()[ServerTimeout]
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

// 设置消息错误信息
func handleError(res rpc.IMessage, err error) (rpc.IMessage, error) {
	res.SetMessageStatusType(rpc.Error)
	if res.Metadata() == nil {
		res.SetMetadata(make(map[string]string))
	}
	res.Metadata()[ServiceError] = err.Error()
	return res, err
}

// 是否是内置类型
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
