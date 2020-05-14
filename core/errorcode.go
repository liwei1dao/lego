package core

const (
	ErrorCode_Success               ErrorCode = 0  //成功
	ErrorCode_NoFindService         ErrorCode = 10 //没有找到远程服务器
	ErrorCode_RpcFuncExecutionError ErrorCode = 11 //Rpc方法执行错误
	ErrorCode_CacheReadError        ErrorCode = 12 //缓存读取失败
	ErrorCode_SqlExecutionError     ErrorCode = 13 //数据库执行错误
)
