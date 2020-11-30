package cbase

import (
	"fmt"
	"os"
	"os/signal"
	"runtime"
	"sync"
	"syscall"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/event"
	"github.com/liwei1dao/lego/sys/log"
)

type defaultModule struct {
	seetring map[string]interface{}
	mi       core.IModule
	closeSig chan bool
	wg       sync.WaitGroup
}

func (this *defaultModule) run() {
	this.mi.Run(this.closeSig)
	this.wg.Done()
}
func (this *defaultModule) destroy() (err error) {
	defer Recover()
	err = this.mi.Destroy()
	if err != nil {
		err = fmt.Errorf("关闭模块【%s】失败 err:%s", this.mi.GetType(), err.Error())
	}
	return
}

func Recover() {
	if r := recover(); r != nil {
		buf := make([]byte, 1024)
		l := runtime.Stack(buf, false)
		log.Panicf("%v: %s", r, buf[:l])
	}
}

type ServiceBase struct {
	closesig chan string
	Service  core.IService
	comps    map[core.S_Comps]core.IServiceComp
	modules  map[core.M_Modules]*defaultModule
}

func (this *ServiceBase) Init(service core.IService) (err error) {
	this.closesig = make(chan string, 1)
	this.Service = service
	this.modules = make(map[core.M_Modules]*defaultModule)
	this.Service.InitSys()
	for _, v := range this.comps {
		err = v.Init(this.Service, v)
		if err != nil {
			return
		}
	}
	log.Infof("服务Init完成 %s", this.Service.GetId())
	return nil
}


//配置服务组件
func (this *ServiceBase) OnInstallComp(cops ...core.IServiceComp) {
	this.comps = make(map[core.S_Comps]core.IServiceComp)
	for _, v := range cops {
		if _, ok := this.comps[v.GetName()]; ok {
			log.Errorf("覆盖注册组件【%s】", v.GetName())
		}
		this.comps[v.GetName()] = v
	}
}

func (this *ServiceBase) Start() (err error) {
	for _, v := range this.comps {
		err = v.Start()
		if err != nil {
			return
		}
	}
	log.Infof("服务Start完成 %s", this.Service.GetId())
	return
}

func (this *ServiceBase) Run(mod ...core.IModule) {
	for _, v := range mod {
		if sf, ok := this.Service.GetSettings().Modules[string(v.GetType())]; ok {
			this.modules[v.GetType()] = &defaultModule{
				seetring: sf,
				mi:       v,
				closeSig: make(chan bool, 1),
			}
		} else {
			this.modules[v.GetType()] = &defaultModule{
				seetring: make(map[string]interface{}),
				mi:       v,
				closeSig: make(chan bool, 1),
			}
			log.Warnf("注册模块【%s】 没有对应的配置信息", v.GetType())
		}
	}
	for _, v := range this.modules {
		err := v.mi.Init(this.Service, v.mi, v.seetring)
		if err != nil {
			panic(fmt.Sprintf("初始化模块【%s】错误 err:%s", v.mi.GetType(), err.Error()))
		}
	}
	for _, v := range this.modules {
		err := v.mi.Start()
		if err != nil {
			panic(fmt.Sprintf("启动模块【%s】错误 err:%v", v.mi.GetType(), err))
		}
	}
	for _, v := range this.modules {
		v.wg.Add(1)
		go v.run()
	}
	event.TriggerEvent(core.Event_ServiceStartEnd) //广播事件
	//监听外部关闭服务信号
	c := make(chan os.Signal, 1)
	//添加进程结束信号
	signal.Notify(c,
		os.Interrupt,    //退出信号 ctrl+c退出
		os.Kill,         //kill 信号
		syscall.SIGHUP,  //终端控制进程结束(终端连接断开)
		syscall.SIGINT,  //用户发送INTR字符(Ctrl+C)触发
		syscall.SIGTERM, //结束程序(可以被捕获、阻塞或忽略)
		syscall.SIGQUIT) //用户发送QUIT字符(Ctrl+/)触发
	select {
	case sig := <-c:
		log.Errorf("服务器关闭 signal = %v\n", sig)
	case <-this.closesig:
		log.Errorf("服务器关闭\n")
	}
}

func (this *ServiceBase) Close(closemsg string) {
	this.closesig <- closemsg
}

func (this *ServiceBase) Destroy() (err error) {
	for _, v := range this.modules {
		v.closeSig <- true
		v.wg.Wait()
		err = v.destroy()
		if err != nil {
			return
		}
	}
	for _, v := range this.comps {
		err = v.Destroy()
		if err != nil {
			return
		}
	}
	return
}

func (this *ServiceBase) GetModule(ModuleName core.M_Modules) (module core.IModule, err error) {
	if v, ok := this.modules[ModuleName]; ok {
		return v.mi, nil
	} else {
		return nil, fmt.Errorf("未装配模块【%s】", ModuleName)
	}
}

func (this *ServiceBase) GetComp(CompName core.S_Comps) (comp core.IServiceComp, err error) {
	if v, ok := this.comps[CompName]; ok {
		return v, nil
	} else {
		return nil, fmt.Errorf("未装配组件【%s】", CompName)
	}
}
