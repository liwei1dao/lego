package lego

import (
	"runtime"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
)

//启动服务
func Run(service core.IService, mod ...core.IModule) {
	cpuNum := runtime.NumCPU() //获得当前设备的cpu核心数
	runtime.GOMAXPROCS(cpuNum) //设置需要用到的cpu数量
	err := service.Init(service)
	if err != nil {
		log.Panicf("服务初始化失败 err=%s", err.Error())
	}
	err = service.Start()
	if err != nil {
		log.Panicf("服务启动失败 err=%s", err.Error())
	}
	service.Run(mod...)
	err = service.Destroy()
	if err != nil {
		log.Panicf("服务销毁失败 err=%s", err.Error())
	}
	log.Infof("服务【%s】关闭成功", service.GetId())
}

//错误采集
func Recover(tag string) {
	if r := recover(); r != nil {
		buf := make([]byte, 1024)
		l := runtime.Stack(buf, false)
		log.Errorf("%s - %v: %s", tag, r, buf[:l])
	}
}
