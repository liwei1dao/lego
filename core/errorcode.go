package core

import "fmt"

///内置错误码 0-1000  请外部应用服务不要占用
const (
	ErrorCode_Success                  ErrorCode = 0  //成功
	ErrorCode_NoFindService            ErrorCode = 10 //没有找到远程服务器
	ErrorCode_RpcFuncExecutionError    ErrorCode = 11 //Rpc方法执行错误
	ErrorCode_CacheReadError           ErrorCode = 12 //缓存读取失败
	ErrorCode_SqlExecutionError        ErrorCode = 13 //数据库执行错误
	ErrorCode_ReqParameterError        ErrorCode = 14 //请求参数错误
	ErrorCode_SignError                ErrorCode = 15 //签名错误
	ErrorCode_InsufficientPermissions  ErrorCode = 16 //权限不足
	ErrorCode_NoRoute                  ErrorCode = 17 //没有路由
	ErrorCode_LocalRouteExecutionError ErrorCode = 18 //本地路由执行错误
	//gate

	ErrorCode_NoLogin ErrorCode = 101 //未登录
)

var ErrorCodeMsg = map[ErrorCode]string{
	ErrorCode_Success:                  "成功",
	ErrorCode_NoFindService:            "没有找到远程服务器",
	ErrorCode_RpcFuncExecutionError:    "Rpc方法执行错误",
	ErrorCode_CacheReadError:           "缓存读取失败",
	ErrorCode_SqlExecutionError:        "数据库执行错误",
	ErrorCode_ReqParameterError:        "请求参数错误",
	ErrorCode_SignError:                "签名错误",
	ErrorCode_InsufficientPermissions:  "权限不足",
	ErrorCode_NoLogin:                  "未登录",
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
