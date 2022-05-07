package blockcache_test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/liwei1dao/lego/sys/blockcache"
	"github.com/liwei1dao/lego/sys/log"
)

func Test_sys(t *testing.T) {
	if err := log.OnInit(nil, log.SetLoglevel(log.DebugLevel), log.SetDebugMode(true)); err != nil {
		fmt.Printf("log init err:%v", err)
		return
	}
	log.Debugf("log init succ")
	if sys, err := blockcache.NewSys(blockcache.SetCacheMaxSzie(100)); err != nil {
		log.Debugf("livego init err:%v", err)
		return
	} else {
		go func() {
			for {
				sys.In() <- "liwei1dao"
				log.Debugf("In:liwei1dao")
			}
		}()
		go func() {
			for {
				for v := range sys.Out() {
					log.Debugf("Out:%v", v)
					time.Sleep(time.Second)
				}
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
