package core

type S_Category string //服务类别 例如 网关服务 游戏服务 业务服务   主要用于服务功能分类
type M_Modules string  //模块类型
type S_Comps string    //服务器组件类型
type ErrorCode int32   //错误码
type Event_Key string  //事件Key
type Rpc_Key string    //RPC
type Redis_Key string  //Redis缓存
type SqlTable string   //数据库表定义
type CustomRoute uint8 //自定义网关

const ( //默认事件
	Event_ServiceStartEnd  Event_Key = "ServiceStartEnd"  //服务完全启动完毕
	Event_FindNewService   Event_Key = "FindNewService"   //发现新的服务
	Event_UpDataOldService Event_Key = "UpDataOldService" //发现新的服务
	Event_LoseService      Event_Key = "LoseService"      //丢失服务
	Event_RegistryStart    Event_Key = "RegistryStart"    //注册表系统启动成功
)

const (
	S_Category_SystemService   S_Category = "SystemService"   //系统服务类型
	S_Category_GateService     S_Category = "GateService"     //网关服务类型
	S_Category_BusinessService S_Category = "BusinessService" //业务服务器
)

type ServiceSttings struct {
	Id       string                            //服务id
	Type     string                            //服务类型
	Version  int                               //服务版本
	Settings map[string]interface{}            //服务配置
	Sys      map[string]map[string]interface{} //系统配置
	Modules  map[string]map[string]interface{} //模块配置
}

type IService interface {
	GetId() string                                              //获取服务id
	GetType() string                                            //获取服务类型
	GetVersion() int32                                          //获取服务版本
	GetWorkPath() string                                        //获取服务工作目录
	GetSettings() ServiceSttings                                //获取服务配置表信息
	Init(service IService) (err error)                          //初始化接口
	InitSys()                                                   //初始化系统
	OnInstallComp(cops ...IServiceComp)                         //组装服务组件
	Start() (err error)                                         //启动服务
	Run(mods ...IModule)                                        //运行服务
	Close(closemsg string)                                      //关闭服务
	Destroy() (err error)                                       //销毁服务
	GetComp(CompName S_Comps) (comp IServiceComp, err error)    //获取组件
	GetModule(ModuleName M_Modules) (module IModule, err error) //获取模块
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
type IServiceMonitor interface {
	IModule
	RegisterServiceSettingItem(name string, iswrite bool, value interface{}, f func(newvalue string) (err error))                  //注册服务级别的Setting
	RegisterModuleSettingItem(module M_Modules, name string, iswrite bool, value interface{}, f func(newvalue string) (err error)) //注册模块级别的Setting
}

//Monitor 数据
type (
	SettingItem struct {
		ItemName string
		IsWrite  bool
		Data     interface{}
	}
	ServiceMonitor struct { //服务监听
		ServiceId        string                       //服务Id
		ServiceType      string                       //服务类型
		ServiceCategory  S_Category                   //服务列表
		ServiceVersion   int32                        //服务版本
		ServiceTag       string                       //服务集群
		Pid              int32                        //进程Id
		Pname            string                       //进程名称
		MemoryInfo       []float32                    //内存使用量
		CpuInfo          []float64                    //Cpu使用量
		TotalGoroutine   []int                        //总的协程数
		ServicePreWeight []int32                      //服务权重
		Setting          map[string]*SettingItem      //服务器配置信息
		ModuleMonitor    map[M_Modules]*ModuleMonitor //模块监听信息
	}
	ModuleMonitor struct { //模块监听
		ModuleName M_Modules               //模块名称
		Setting    map[string]*SettingItem //模块配置信息
	}
)
