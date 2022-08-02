package livego_test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"

	"github.com/liwei1dao/lego/sys/livego"
	"github.com/liwei1dao/lego/sys/log"
)

func Test_sys(t *testing.T) {
	if err := log.OnInit(nil, log.SetLoglevel(log.DebugLevel)); err != nil {
		fmt.Printf("log init err:%v", err)
		return
	}
	log.Debugf("log init succ")
	if _, err := livego.NewSys(); err != nil {
		log.Debugf("livego init err:%v", err)
		return
	}
	log.Debugf("livego init succ")
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		fmt.Printf("terminating: via signal\n")
	}
}
