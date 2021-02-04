package core

import "fmt"

const (
	ErrorCode_Success               ErrorCode = 0  //成功
	ErrorCode_NoFindService         ErrorCode = 10 //没有找到远程服务器
	ErrorCode_RpcFuncExecutionError ErrorCode = 11 //Rpc方法执行错误
	ErrorCode_CacheReadError        ErrorCode = 12 //缓存读取失败
	ErrorCode_SqlExecutionError     ErrorCode = 13 //数据库执行错误
	ErrorCode_ReqParameterError     ErrorCode = 14 //请求参数错误
	ErrorCode_SignError             ErrorCode = 15 //签名错误
	//gate
	ErrorCode_NoRoute                  ErrorCode = 201 //没有路由
	ErrorCode_LocalRouteExecutionError ErrorCode = 202 //本地路由执行错误
)

var ErrorCodeMsg = map[ErrorCode]string{
	ErrorCode_Success:                  "成功",
	ErrorCode_NoFindService:            "没有找到远程服务器",
	ErrorCode_RpcFuncExecutionError:    "Rpc方法执行错误",
	ErrorCode_CacheReadError:           "缓存读取失败",
	ErrorCode_SqlExecutionError:        "数据库执行错误",
	ErrorCode_ReqParameterError:        "请求参数错误",
	ErrorCode_SignError:                "签名错误",
	ErrorCode_NoRoute:                  "没有路由",
	ErrorCode_LocalRouteExecutionError: "本地路由执行错误",
}

func GetErrorCodeMsg(code ErrorCode) string {
	if v, ok := ErrorCodeMsg[code]; ok {
		return v
	} else {
		return fmt.Sprintf("ErrorCode:%d", code)
	}
}
