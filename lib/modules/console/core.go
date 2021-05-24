package console

import (
	"fmt"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/lib/modules/http"
	"github.com/liwei1dao/lego/sys/mgo"
	"github.com/liwei1dao/lego/sys/redis"
)

const (
	UserKey string = "UserId"
)

const ( //DB
	Sql_ConsoleUserDataIdTable core.SqlTable = "console_usersid" //话数据
	Sql_ConsoleUserDataTable   core.SqlTable = "console_users"
)

const ( //Cache
	Cache_ConsoleToken          core.Redis_Key = "Console_Token:%s"          //用户数据缓存
	Cache_ConsoleUsers          core.Redis_Key = "Console_Users:%d"          //用户数据缓存
	Cache_ConsoleClusterMonitor core.Redis_Key = "Console_Clustermonitor:%s" //集群监听数据
	Cache_ConsoleHostMonitor    core.Redis_Key = "Console_Hostmonitor"       //主机监听数据
	Cache_ConsoleCaptcha        core.Redis_Key = "Console_Captcha:%s-%d"     //验证码缓存
)

const ( //event
	Event_Registered core.Event_Key = "Console_Registered" //用户注册
	Event_Login      core.Event_Key = "Console_Login"      //用户登录
)

type (
	IConsole interface {
		http.IHttp
		Options() IOptions
		Cache() ICache
		DB() IDB
		Captcha() ICaptcha
		Hostmonitorcomp() IHostmonitorcomp
		Clustermonitorcomp() IClustermonitorcomp
		ParamSign(param map[string]interface{}) (sign string)
		CheckToken(c *http.Context)
		HttpStatusOK(c *http.Context, code core.ErrorCode, data interface{})
		CreateToken(uId uint32) (token string)
		checkToken(token string, uId uint32) (check bool)
	}
	ICache interface {
		core.IModuleComp
		GetRedis() redis.IRedis
		QueryToken(token string) (uId uint32, err error)
		WriteToken(token string, uId uint32) (err error)
		CleanToken(token string) (err error)
		QueryUserData(uId uint32) (result *Cache_UserData, err error)
		WriteUserData(data *Cache_UserData) (err error)
		CleanUserData(uid uint32) (err error)
		AddNewClusterMonitor(data map[string]*ClusterMonitor)
		GetClusterMonitor(sIs string, timeleng int) (result []*ClusterMonitor, err error)
		AddNewHostMonitor(data *HostMonitor)
		GetHostMonitor(timeleng int) (result []*HostMonitor, err error)
	}
	IDB interface {
		core.IModuleComp
		GetMgo() mgo.IMongodb
		QueryUserDataById(uid uint32) (result *DB_UserData, err error)
		QueryUserDataByPhonOrEmail(phonoremail string) (result *DB_UserData, err error)
		LoginUserDataByPhonOrEmail(data *DB_UserData) (result *DB_UserData, err error)
		UpDataUserData(data *DB_UserData) (err error)
	}
	ICaptcha interface {
		core.IModuleComp
		SendEmailCaptcha(email, captcha string) (err error)
		QueryCaptcha(cId string, ctype CaptchaType) (captcha string, err error)
		WriteCaptcha(cId, captcha string, ctype CaptchaType)
	}
	IHostmonitorcomp interface {
		core.IModuleComp
		HostInfo() *HostInfo
		CpuInfo() []*CpuInfo
		MemoryInfo() *MemoryInfo
		GetHostMonitorData(queryTime QueryMonitorTime) (result *QueryHostMonitorDataResp)
	}
	IClustermonitorcomp interface {
		core.IModuleComp
		GetClusterMonitorDataResp(queryTime QueryMonitorTime) (result map[string]map[string]interface{})
	}
	ConsoleUserId struct {
		Id uint32 `bson:"_id"`
	}
)

const (
	ErrorCode_Success               core.ErrorCode = 0  //成功
	ErrorCode_NoFindService         core.ErrorCode = 10 //没有找到远程服务器
	ErrorCode_RpcFuncExecutionError core.ErrorCode = 11 //Rpc方法执行错误
	ErrorCode_CacheReadError        core.ErrorCode = 12 //缓存读取失败
	ErrorCode_SqlExecutionError     core.ErrorCode = 13 //数据库执行错误
	ErrorCode_ReqParameterError     core.ErrorCode = 14 //请求参数错误
	ErrorCode_SignError             core.ErrorCode = 15 //签名错误

	//Login登录错误码
	ErrorCode_TokenExpired                core.ErrorCode = 100 //Token过期
	ErrorCode_NoLOgin                     core.ErrorCode = 101 //未登录
	ErrorCode_AuthorityLow                core.ErrorCode = 102 //权限不足
	ErrorCode_AlreadyRegister             core.ErrorCode = 103 //帐号已被注册
	ErrorCode_LoginAccountOrPasswordError core.ErrorCode = 104 //登录帐号密码错误
	ErrorCode_UserCookieLose              core.ErrorCode = 105 //用户Cookie数据丢失
	ErrorCode_UserDataError               core.ErrorCode = 106 //用户数据错误
)

var ErrorCodeMsg = map[core.ErrorCode]string{
	ErrorCode_Success:                     "成功",
	ErrorCode_NoFindService:               "没有找到远程服务器",
	ErrorCode_RpcFuncExecutionError:       "Rpc方法执行错误",
	ErrorCode_CacheReadError:              "缓存读取失败",
	ErrorCode_SqlExecutionError:           "数据库执行错误",
	ErrorCode_ReqParameterError:           "请求参数错误",
	ErrorCode_SignError:                   "签名错误",
	ErrorCode_TokenExpired:                "Token过期",
	ErrorCode_NoLOgin:                     "未登录",
	ErrorCode_AuthorityLow:                "您的权限不足",
	ErrorCode_AlreadyRegister:             "帐号已被注册",
	ErrorCode_LoginAccountOrPasswordError: "登录帐号密码错误",
	ErrorCode_UserCookieLose:              "用户Cookie数据丢失",
}

func GetErrorCodeMsg(code core.ErrorCode) string {
	if v, ok := ErrorCodeMsg[code]; ok {
		return v
	} else {
		return fmt.Sprintf("ErrorCode:%d", code)
	}
}
