package blockcache_test

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"testing"
	"time"

	"github.com/liwei1dao/lego/sys/blockcache"
)

func Test_sys(t *testing.T) {
	if sys, err := blockcache.NewSys(blockcache.SetCacheMaxSzie(100)); err != nil {
		fmt.Printf("livego init err:%v \n", err)
		return
	} else {
		closeSignal := make(chan struct{})
		go func() {
		locp:
			for {
				select {
				case <-closeSignal:
					break locp
				default:
					sys.In() <- "liwei1dao"
					fmt.Printf("In:liwei1dao\n")
				}
			}
			fmt.Printf("In:End\n")
		}()
		go func() {
			for v := range sys.Out() {
				fmt.Printf("Out:%v\n", v)
				time.Sleep(time.Second)
			}
			fmt.Printf("Out:End\n")
		}()
		go func() {
			time.Sleep(time.Second * 3)
			closeSignal <- struct{}{}
			sys.Close()
		}()
	}
	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		fmt.Printf("terminating: via signal\n")
	}
}
