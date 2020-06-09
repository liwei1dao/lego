package workerpools

import (
	"context"
	"testing"
	"time"
)

func Test_workerpools(t *testing.T) {
	pools, _ := NewTaskPools(SetMaxWorkers(1), SetTaskTimeOut(time.Second*4))
	go func() {
		pools.Submit(func(ctx context.Context, cancel context.CancelFunc, agrs ...interface{}) {
			time.Sleep(time.Second * 6)
			t.Logf("liwei2dao")
			cancel()
		})
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
			t.Logf("liwei1dao")
			cancel()
		})
	}()

	time.Sleep(time.Second * 10)
}
