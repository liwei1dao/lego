package workerpools

import "context"

type (
	IWorkerPool interface {
		Stop()
		StopWait()
		IsStop() bool
		Submit(task func(ctx context.Context, cancel context.CancelFunc))
		SubmitWait(task func(ctx context.Context, cancel context.CancelFunc))
		WaitingQueueSize() int
	}
)
