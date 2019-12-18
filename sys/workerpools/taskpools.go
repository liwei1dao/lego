package workerpools

import (
	"context"
	"lego/core"
)

var (
	pools IWorkerPool
)

func OnInit(s core.IService, opt ...Option) (err error) {
	pools = newWorkerPool(opt...)
	return
}

func NewTaskPools(opt ...Option) (pools IWorkerPool, err error) {
	pools = newWorkerPool(opt...)
	return
}

func Stop()                                                                { pools.Stop() }
func StopWait()                                                            { pools.StopWait() }
func IsStop() bool                                                         { return pools.IsStop() }
func Submit(task func(ctx context.Context, cancel context.CancelFunc))     { pools.Submit(task) }
func SubmitWait(task func(ctx context.Context, cancel context.CancelFunc)) { pools.SubmitWait(task) }
func WaitingQueueSize() int                                                { return pools.WaitingQueueSize() }
