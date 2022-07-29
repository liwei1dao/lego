package rpcl

import (
	"reflect"
	"unicode"
	"unicode/utf8"

	"github.com/liwei1dao/lego/core"
	lcore "github.com/liwei1dao/lego/sys/rpcl/core"
	"github.com/liwei1dao/lego/sys/rpcl/protocol"
)

//获取心跳消息包
func getHeartbeat(node core.ServiceNode) []byte {
	req := protocol.GetPooledMsg()
	req.SetMessageType(lcore.Request)
	req.SetHeartbeat(true)
	req.SetOneway(true)
	req.SetServiceMethod("")
	codec := codecs[lcore.ProtoBuffer]
	data, _ := codec.Marshal(node)
	req.SetPayload(data)
	allData := req.EncodeSlicePointer()
	defer func() {
		protocol.PutData(allData)
		protocol.FreeMsg(req)
	}()
	return *allData
}

//设置消息错误信息
func handleError(res lcore.IMessage, err error) (lcore.IMessage, error) {
	res.SetMessageStatusType(lcore.Error)
	if res.Metadata() == nil {
		res.SetMetadata(make(map[string]string))
	}
	res.Metadata()[lcore.ServiceError] = err.Error()
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
