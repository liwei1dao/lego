@[TOC](lego  开发框架使用)

# lego 开发框架简介
在我前端框架完成之后就一直想进军后端框架设计，刚好golang 这是兴起了,于是顺势开始研究其golang，详细看过leaf和mqant这两套开源框架的源码以及思路，原先的几个版本也基本上行都是拿着别人的框架改的，后面逐渐修改加入自己的想法以及理解，可以看到其实lego 的核心设计开始和前端u3d_fw的设计一致 只不过多了一个服务的概念，此框架也继承了我框架设计的基本想法 横向和纵向的可延申性和业务无关性 两大特点


# lego  的导入
go get github.com/liwei1dao/lego
# lego 框架目录结构说明
- logo 总目录
	* base    //框架基础服务封装内置single(独立)和cluster(集群)两类服务基类
	* core 	//框架核心接口定义以及service和module以及各类组件和数据类型定义
	* lib  //框架代码库 用户代码积累，先继承进来 模块(gate,http),服务组件(comp_gateroute 路由组件) ,以及模块组件(comp_gate 网关业务组件)，后期还将继续集成进来更多个功能性模块以及组件
	* sys //框架系统库，类似于工具库，但是由于功能型比较强而且涉及到自己的数据结构管理切不舍和做外部扩展所以开辟了这个系统库包含例如(monodb,redis,log,rpc...)这类独立的功能集，后期也将继续扩展
	* utils  //工具集 封装各类实用型工具以及数据容器之类的对象
# base 服务封装
- IService 所有服务都必须继承于它 (在core目录下)
```
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
```

 
- single 独立服务器 
	* 此服务为独立服务，一般用于提供单独功能且不需要其他服务协助的功能性服务开发，例如简单的web服务器和一些微小型业务服务器
```
type ISingleService interface {
	core.IService
}
```

- cluster 集群服务器
	* 此服务用于大型服务集群开发使用，涉及到服务间通信以及服务发现，更新，丢失，以及服务点对点，一对多，以及消息订阅于发布等服务间消息发送机制，
	
```
type IClusterService interface {
	core.IService
	GetTag() string                                                                                                   //获取集群标签
	GetCategory() core.S_Category                                                                                     //服务类别 例如游戏服
	GetRpcId() string                                                                                                 //获取rpc通信id
	GetPreWeight() int32                                                                                              //集群服务负载值 暂时可以不用理会
	GetSessionsByCategory(category core.S_Category) (ss []core.IServiceSession)                                       //按服务类别获取服务列表
	DefauleRpcRouteRules(stype string) (ss core.IServiceSession, err error)                                           //默认rpc路由规则
	RpcInvokeById(sId string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (result interface{}, err error)     //执行远程服务Rpc方法
	RpcInvokeByType(sType string, rkey core.Rpc_Key, iscall bool, arg ...interface{}) (result interface{}, err error) //根据路由规则执行远程方法
	ReleaseRpc(rkey core.Rpc_Key, arg ...interface{})                                                                 //发布Rpc
	Register(id core.Rpc_Key, f interface{}) (err error)                                                              //注册RPC远程方法
	RegisterGO(id core.Rpc_Key, f interface{}) (err error)                                                            //注册RPC远程方法
	Subscribe(id core.Rpc_Key, f interface{}) (err error)                                                             //订阅Rpc
	UnSubscribe(id core.Rpc_Key, f interface{}) (err error)                                                           //订阅Rpc
}
```
# core 接口结构定义
- 自定义数据结构·
	 type S_Category string //服务类别 例如 网关服务 游戏服务 业务服务   主要用于服务功能分类
	type M_Modules string  //模块类型
	type S_Comps string    //服务器组件类型
	type ErrorCode int32   //错误码
	type Event_Key string  //事件Key名
	type Rpc_Key string    //RPC接口名
	type Redis_Key string  //Redis缓存键
	type SqlTable string   //数据库表名
- 服务接口

	```
	//服务配置数据结构
	type ServiceSttings struct {
		Settings map[string]interface{}
		Modules  map[string]map[string]interface{}
	}
	//服务接口
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
	//服务器组件接口
	type IServiceComp interface {
		GetName() S_Comps
		Init(service IService, comp IServiceComp) (err error)
		Start() (err error)
		Destroy() (err error)
	}
	```
- 模块接口定义

	```
	//模块接口
	type IModule interface {
		GetType() M_Modules
		Init(service IService, module IModule, setting map[string]interface{}) (err error)
		OnInstallComp()
		Start() (err error)
		Run(closeSig chan bool) (err error)
		Destroy() (err error)
	}
	//模块组件接口
	type IModuleComp interface {
		Init(service IService, module IModule, comp IModuleComp, setting map[string]interface{}) (err error)
		Start() (err error)
		Destroy() (err error)
	}
	```
- 会话对象
 	* 服务会话 
		```
		//服务会话 服务见通信以及基本信息传递的对象
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
		```
	* 用户会话 客户端对象消息发送以及基本信息传递的对象
		```
		type IUserSession interface {
			GetSessionId() string
			GetIP() string
			GetGateId() string
			SendMsg(comdId uint16, msgId uint16, msg interface{}) (err error)
			Close() (err error)
		}
		
		```

# lib 集成功能模块及组件
- s_comps 
	* comp_gateroute 为服务提供接受网关服务的消息，只有组装了这个组件的服务才可以接收到来自自己或者其他网关服务的消息推送
- modules
	* gate 网关服务模块 提供tcp和websocket服务给客户端 模块需要自行扩展实现必要接口 可参考demo的gate模块实现实例
	* http web服务模块 提供http或者https服务 一样可以参考demo的实现实例
- module-comps
	* comp_gate 配合网关服务模块 集成这个组件的业务组件可以处理来自网关分配的消息

# sys 集成系统 
- registry 注册表系统 配合consul服务实现服务集群的发现更新和丢失事件以及集群下rpc发布订阅数据的保存
- rpc rpc消息系统，为集群服务下提供服务见消息通信
- workerpools 工作池系统，为高并发服务以及业务通过安全切高性能的处理流
- log 日志系统 集成了zap系统进去，此系统可以扩展日志接口采用自定义以及其他第三方日志系统都可
- proto 协议系统 为gate和rpc消息序列化支持，内部封装有默认消息结构，可以通过启动参数设置自定义消息结构
- event 事件系统，为服务提供事件的注册监听处理功能
- mgo mogodb数据库系统，封装有官方驱动库，简化系统调用接口
- sqls sqlserver数据库系统，封装有官方驱动库，简化系统调用接口
- redis redis缓存数据系统，提供了各类数据的存储及读取接口
- timer 计时器系统，封装了高效的时间轮系统
- sdks 内置集成各类第三方服务代码例如阿里云的短信以及邮件消息推送实现
# demo下载地址
demo 项目地址:
https://github.com/liwei1dao/lego_demo.git or
https://gitee.com/liwei1dao/lego_demo.git
