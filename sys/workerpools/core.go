package workerpools

import "context"

type (
	IWorkerPool interface {
		Stop()
		StopWait()
		IsStop() bool
		Submit(task func(ctx context.Context, cancel context.CancelFunc, agrs ...interface{}), agrs ...interface{})
		SubmitWait(task func(ctx context.Context, cancel context.CancelFunc, agrs ...interface{}), agrs ...interface{})
		WaitingQueueSize() int
	}
)

var (
	defsys IWorkerPool
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys IWorkerPool, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func NewTaskPools(opt ...Option) (pools IWorkerPool, err error) {
	pools = newWorkerPool(opt...)
	return
}

func Stop()        { defsys.Stop() }
func StopWait()    { defsys.StopWait() }
func IsStop() bool { return defsys.IsStop() }
func Submit(task func(ctx context.Context, cancel context.CancelFunc, agrs ...interface{}), agrs ...interface{}) {
	defsys.Submit(task, agrs...)
}
func SubmitWait(task func(ctx context.Context, cancel context.CancelFunc, agrs ...interface{}), agrs ...interface{}) {
	defsys.SubmitWait(task, agrs...)
}
func WaitingQueueSize() int { return defsys.WaitingQueueSize() }
