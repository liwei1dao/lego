@[TOC](lego  开发框架使用)

## lego GoLang开发框架
lego是一款微服务快速开发框架,框架采用容器化组件设计理念，实现业务功能的随意组合和部署，框架由四大核心快组成:Server,Module,Sys,utils 内部更是集成了诸多基础模块和使用系统，例如 网关模块(gate) Web模块(http),监控后台(console),通信系统(rpc),服务发现注册系统(registry),数据系统(redis,mgo,sql...

## lego 功能介绍
lego 的设计之初大量参考了mqant设计理念，所以还有一部分功能设计与mqant一致，后期发现mqant的设计模式对service 对象定义的不是很清晰,所以我后面引入了service + module 的设计方式，后期受到gorm的插件设计形成了sys功能库
1. servic 服务对象，进程的主要载体,可以挂在service-component(服务组件),module(功能模块),服务可以根据自己的业务需求挂在相应的服务组件以及业务模块同时可以启动相应的系统
```
  package main

  import (
    "flag"
    "github.com/liwei1dao/lego"
    "github.com/liwei1dao/lego/base/cluster"
    "github.com/liwei1dao/lego/core"
  )

  var (
    sID = flag.String("sId", "console_1", "获取需要启动的服务id,id不同,读取的配置文件也不同") //启动服务的Id
  )

  func main() {
    flag.Parse()
    s := NewService(
      cluster.SetId(*sID),
    )
    s.OnInstallComp( //装备组件
    )
    lego.Run(s, //运行模块
    )
  }

  func NewService(ops ...cluster.Option) core.IService {
    s := new(Demo1Service)
    s.Configure(ops...)
    return s
  }

  type Demo1Service struct {
	  cluster.ClusterService
  }

  func (this *Demo1Service) InitSys() {
    this.ClusterService.InitSys()
  }
```
2. service-component 服务组件 为service服务对象提供扩展功能,例如:lego/lib/s_comps/comp_gateroute.go,此服务组件可以实现非网关服务的业务服务器接收到来自网关服务转发的用户请求信息
```
  func main() {
    flag.Parse()
    s := NewService(
      cluster.SetId(*sID),
    )
    s.OnInstallComp( //装备组件
      s_comps.NewGateRouteComp(), //挂在网关服务组件 接收用户消息
    )
    lego.Run(s, //运行模块
    )
  }
```
3. module 业务功能模块，项目模块化设计方案，采用Module+Component的设计模式，一个Module下可以挂载多个Component，Module负责管理这些组件以及提供外部接口，Component实现具体的模块业务，这个设计模式参考lgu3d框架设计思路(:smile::smile::smile:),框架内置有console(控制台),gate(网关),http(web业务模块集成于gin框架),live(流媒体 还不完善,属于测试阶段),monitor(监控模块 配合console可以在后台实时监控服务运行情况以及动态设置模块配置数据)
  
console 模块实现 配合lgvue可以搭建项目后台
```
  package console

  import (
    "github.com/liwei1dao/lego/core"
    "github.com/liwei1dao/lego/lib"
    "github.com/liwei1dao/lego/lib/modules/console"
  )

  func NewModule() core.IModule {
    m := new(Console)
    return m
  }

  type Console struct {
    console.Console
  }

  func (this *Console) GetType() core.M_Modules {
    return lib.SM_ConsoleModule
  }

  func (this *Console) NewOptions() (options core.IModuleOptions) {
    return new(Options)
  }

  func (this *Console) Init(service core.IService, module core.IModule, options core.IModuleOptions) (err error) {
    err = this.Console.Init(service, module, options)
    return
  }

  func (this *Console) Start() (err error) {
    err = this.Console.Start()
    return
  }
```
options.go 模块配置文件,自动序列化服务配置文件.toml文件下模块配置
```
  package console

  import (
    "github.com/liwei1dao/lego/lib/modules/console"
    "github.com/liwei1dao/lego/utils/mapstructure"
  )

  type Options struct {
    console.Options
  }

  func (this *Options) LoadConfig(settings map[string]interface{}) (err error) {
    if err = this.Options.LoadConfig(settings); err == nil {
      if settings != nil {
        err = mapstructure.Decode(settings, this)
      }
    }
    return
  }
```
4. module-component 模块业务组件,模块内具体实现业务功能的单元
```
package console

import (
	"fmt"
	reflect "reflect"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/core/cbase"
	"github.com/liwei1dao/lego/sys/redis"
)

type CacheComp struct {
	cbase.ModuleCompBase
	module IConsole
	redis  redis.IRedis
}

func (this *CacheComp) Init(service core.IService, module core.IModule, comp core.IModuleComp, options core.IModuleOptions) (err error) {
	err = this.ModuleCompBase.Init(service, module, comp, options)
	this.module = module.(IConsole)
	if this.redis, err = redis.NewSys(
		redis.SetRedisUrl(this.module.Options().GetRedisUrl()),
		redis.SetRedisDB(this.module.Options().GetRedisDB()),
		redis.SetRedisPassword(this.module.Options().GetRedisPassword()),
	); err != nil {
		err = fmt.Errorf("redis[%s]err:%v", this.module.Options().GetRedisUrl(), err)
	}
	return
}
```
## lego 目录结构介绍
1. base 基础集群服务和基础单服务的实现
2. core lego基础框架设计
3. lib lego代码库,内置服务端开发需要的基本模块业务和服务组件
4. sys lego系统库,集成大量第三方服务插件以及集群服务相关通信支持系统
5. utils lego工具集
## lego 安装
go github.com/liwei1dao/lego
## lego 使用
 可以参考 https://github.com/liwei1dao/demo 下 lego_demo 的实现
 此演示项目可以直接通过docker运行