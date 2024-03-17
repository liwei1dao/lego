package discovery_test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/discovery"
)

func Test_sys(t *testing.T) {
	node := &core.ServiceNode{
		Tag:  "demo",
		Type: "gate",
		Id:   "gate_1",
		Addr: "127.0.0.1:7851",
	}
	if sys, err := discovery.NewSys(
		discovery.SetStoreType(discovery.StoreConsul),
		discovery.SetEndpoints([]string{"10.0.0.9:8500"}),
		discovery.SetBasePath("demo"),
		discovery.SetUpdateInterval(time.Second*10),
		discovery.SetServiceNode(node),
	); err != nil {
		fmt.Printf("err:%v\n", err)
		return
	} else {
		sys.Start()
		ss := sys.GetServices()
		fmt.Printf("ss:%v\n", ss)
	}

	//监听外部关闭服务信号
	c := make(chan os.Signal, 1)
	//添加进程结束信号
	signal.Notify(c,
		os.Interrupt,    //退出信号 ctrl+c退出
		syscall.SIGHUP,  //终端控制进程结束(终端连接断开)
		syscall.SIGINT,  //用户发送INTR字符(Ctrl+C)触发
		syscall.SIGTERM, //结束程序(可以被捕获、阻塞或忽略)
		syscall.SIGQUIT) //用户发送QUIT字符(Ctrl+/)触发
	select {
	case sig := <-c:
		fmt.Println("关闭 signal\n", sig)
	}
}
