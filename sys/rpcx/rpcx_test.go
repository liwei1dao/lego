package rpcx

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
)

func Test_Sys(t *testing.T) {
	if err := log.OnInit(nil); err != nil {
		fmt.Printf("err:%v", err)
		return
	}
	node, _ := core.NewServiceNode("stag=%s&stype=%s&sid=%s&version=%s&addr=%s")
	if sys, err := NewSys(
		SetServiceNode(node),
		SetETCDServers([]string{"10.0.0.9:2379"}),
	); err != nil {
		fmt.Printf("err:%v", err)
		return
	} else {
		if err = sys.RegisterFunction(RpcxTestHandle); err != nil {
			fmt.Printf("err:%v", err)
			return
		}
		if err = sys.Start(); err != nil {
			fmt.Printf("err:%v", err)
			return
		}
		go func() {
			time.Sleep(time.Second * 3)
			if err = sys.Call(context.Background(), "worker/worker_1", "Mul", &Args{A: 1, B: 2}, &Reply{}); err != nil {
				fmt.Printf("Call:%v \n", err)
				return
			}
		}()
	}

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		fmt.Printf("terminating: via signal\n")
	}
}

type Args struct {
	A int
	B int
}
type Reply struct {
	Error string
}

func RpcxTestHandle(ctx context.Context, args *Args, reply *Reply) error {
	fmt.Printf("A:%d + B:%d = %d", args.A, args.B, args.A+args.B)
	return nil
}
