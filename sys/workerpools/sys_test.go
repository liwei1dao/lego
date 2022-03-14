package workerpools_test

import (
	"context"
	"testing"
	"time"

	"github.com/liwei1dao/lego/sys/workerpools"
)

func Test_workerpools(t *testing.T) {
	pools, _ := workerpools.NewSys(workerpools.SetMaxWorkers(1), workerpools.SetTaskTimeOut(time.Second*4))
	go func() {
		pools.Submit(func(ctx context.Context, cancel context.CancelFunc, agrs ...interface{}) {
			agr0 := agrs[0].(string)
			time.Sleep(time.Second * 6)
			t.Logf(agr0)
			cancel()
		}, "liwei1dao")
	}()

	// go func() {
	// 	time.Sleep(time.Second * 2)
	// 	pools.Submit(func(ctx context.Context, cancel context.CancelFunc) {
	// 		time.Sleep(time.Second * 1)
	// 		t.Logf("liwei1dao")
	// 	})
	// }()

	go func() {
		time.Sleep(time.Second * 1)
		pools.Submit(func(ctx context.Context, cancel context.CancelFunc, agrs ...interface{}) {
			agr0 := agrs[0].(string)
			t.Logf(agr0)
			cancel()
		}, "liwei2dao")
	}()

	time.Sleep(time.Second * 10)
}
