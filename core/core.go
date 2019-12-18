package core

type S_Category string //服务类别 例如 网关服务 游戏服务 业务服务   主要用于服务功能分类
type M_Modules string  //模块类型
type S_Comps string    //服务器组件类型
type ErrorCode int32   //错误码
type Event_Key string  //事件Key
type Rpc_Key string    //RPC
type Redis_Key string  //Redis缓存
type SqlTable string   //数据库表定义

const ( //默认事件
	Event_ServiceStartEnd  Event_Key = "ServiceStartEnd"  //服务完全启动完毕
	Event_FindNewService   Event_Key = "FindNewService"   //发现新的服务
	Event_UpDataOldService Event_Key = "UpDataOldService" //发现新的服务
	Event_LoseService      Event_Key = "LoseService"      //丢失服务
)

type ServiceSttings struct {
	Settings map[string]interface{}
	Modules  map[string]map[string]interface{}
}

type IService interface {
	GetId() string
	GetType() string
	GetVersion() int32
	GetWorkPath() string
	GetSettings() ServiceSttings
	Init(service IService) (err error)
	InitSys()
	OnInstallComp(cops ...IServiceComp)
	Start() (err error)
	Run(mods ...IModule)
	Close(closemsg string)
	Destroy() (err error)
	GetComp(CompName S_Comps) (comp IServiceComp, err error)
	GetModule(ModuleName M_Modules) (module IModule, err error)
}

type IServiceComp interface {
	GetName() S_Comps
	Init(service IService, comp IServiceComp) (err error)
	Start() (err error)
	Destroy() (err error)
}
type IModule interface {
	GetType() M_Modules
	Init(service IService, module IModule, setting map[string]interface{}) (err error)
	OnInstallComp()
	Start() (err error)
	Run(closeSig chan bool) (err error)
	Destroy() (err error)
}
type IModuleComp interface {
	Init(service IService, module IModule, comp IModuleComp, setting map[string]interface{}) (err error)
	Start() (err error)
	Destroy() (err error)
}
type IServiceSession interface {
	GetId() string
	GetRpcId() string
	GetType() string
	GetVersion() int32
	SetVersion(v int32)
	GetPreWeight() int32
	SetPreWeight(p int32)
	Done()
	CallNR(_func Rpc_Key, params ...interface{}) (err error)
	Call(_func Rpc_Key, params ...interface{}) (interface{}, error)
}

type IUserSession interface {
	GetSessionId() string
	GetIP() string
	GetGateId() string
	SendMsg(comdId uint16, msgId uint16, msg interface{}) (err error)
	Close() (err error)
}
