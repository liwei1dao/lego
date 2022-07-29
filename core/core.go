package core

import "fmt"

type S_Category string //服务类别 例如 网关服务 游戏服务 业务服务   主要用于服务功能分类
type M_Modules string  //模块类型
type S_Comps string    //服务器组件类型
type ErrorCode int32   //错误码
type Event_Key string  //事件Key
type Rpc_Key string    //RPC
type Redis_Key string  //Redis缓存
type SqlTable string   //数据库表定义
type CustomRoute uint8 //自定义网关

const (
	AutoIp = "0.0.0.0"
	AllIp  = "255.255.255.255"
)

const ( //默认事件
	Event_ServiceStartEnd  Event_Key = "ServiceStartEnd"  //服务完全启动完毕
	Event_FindNewService   Event_Key = "FindNewService"   //发现新的服务
	Event_UpDataOldService Event_Key = "UpDataOldService" //发现新的服务
	Event_LoseService      Event_Key = "LoseService"      //丢失服务
	Event_RegistryStart    Event_Key = "RegistryStart"    //注册表系统启动成功
)

type ServiceSttings struct {
	Id       string                            //服务Id
	Type     string                            //服务类型 (相同的服务可以启动多个)
	Tag      string                            //服务集群标签 (相同标签的集群服务可以互相发现和发现)
	Category S_Category                        //服务列表 (用于区分集群服务下相似业务功能的服务器 例如:游戏服务器)
	Ip       string                            //服务所在Ip				()
	Port     int                               //服务rpcx监听端口
	Comps    map[string]map[string]interface{} //服务组件配置
	Sys      map[string]map[string]interface{} //服务系统配置
	Modules  map[string]map[string]interface{} //服务模块配置
}

type IService interface {
	GetId() string                                              //获取服务id
	GetType() string                                            //获取服务类型
	GetVersion() string                                         //获取服务版本
	GetIp() string                                              //获取服务器运ip
	GetPort() int                                               //服务默认端口
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
	NewOptions() (options ICompOptions)
	Init(service IService, comp IServiceComp, options ICompOptions) (err error)
	Start() (err error)
	Destroy() (err error)
}
type IModule interface {
	GetType() M_Modules
	NewOptions() (options IModuleOptions)
	Init(service IService, module IModule, options IModuleOptions) (err error)
	OnInstallComp()
	Start() (err error)
	Run(closeSig chan bool) (err error)
	Destroy() (err error)
}

type ICompOptions interface {
	LoadConfig(settings map[string]interface{}) (err error)
}

type IModuleOptions interface {
	LoadConfig(settings map[string]interface{}) (err error)
}

type IModuleComp interface {
	Init(service IService, module IModule, comp IModuleComp, options IModuleOptions) (err error)
	Start() (err error)
	Destroy() (err error)
}

//服务节点路径
func (this *ServiceNode) GetNodePath() string {
	return fmt.Sprintf("%s/%s/%s", this.Tag, this.Type, this.Id)
}
