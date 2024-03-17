package gin_test

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/liwei1dao/lego/sys/gin"
	"github.com/liwei1dao/lego/sys/gin/engine"
	"github.com/liwei1dao/lego/sys/log"
)

func Test_sys(t *testing.T) {
	if err := log.OnInit(nil,
		log.SetFileName("log.log"),
		log.SetIsDebug(false),
		log.SetEncoder(log.TextEncoder),
	); err != nil {
		fmt.Printf("log init err:%v", err)
		return
	}
	if sys, err := gin.NewSys(); err != nil {
		fmt.Printf("gin init err:%v", err)
	} else {
		sys.GET("/test", func(c *engine.Context) {
			c.JSON(http.StatusOK, "hello")
		})
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

///测试签名
func Test_ParamSign(t *testing.T) {
	origin, sgin := gin.ParamSign("@234%67g12q4*67m12#4l67!", map[string]interface{}{"images": []string{"测试资源.png", "11.jpg"}})
	fmt.Println(origin, sgin)
}
