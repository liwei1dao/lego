package workerpools

import (
	"context"
	"sync"
	"time"

	cont "github.com/liwei1dao/lego/utils/concurrent"
)

const (
	idleTimeoutSec = 5
)

func newWorkerPool(opt ...Option) IWorkerPool {
	opts := newOptions(opt...)
	pool := &WorkerPool{
		taskQueue:    make(chan func(ctx context.Context, cancel context.CancelFunc), 1),
		maxWorkers:   opts.maxWorkers,
		readyWorkers: make(chan chan func(ctx context.Context, cancel context.CancelFunc), opts.defWrokers),
		timeout:      time.Second * idleTimeoutSec,
		tasktimeout:  opts.tasktimeout,
		stoppedChan:  make(chan struct{}),
	}
	go pool.dispatch()
	return pool
}

type WorkerPool struct {
	maxWorkers   int
	timeout      time.Duration //超时释放空闲工作人员
	tasktimeout  time.Duration //任务执行操超时间
	taskQueue    chan func(ctx context.Context, cancel context.CancelFunc)
	readyWorkers chan chan func(ctx context.Context, cancel context.CancelFunc)
	stoppedChan  chan struct{}
	waitingQueue cont.Deque
	stopMutex    sync.Mutex
	stopped      bool
}

func (p *WorkerPool) Stop() {
	p.stop(false)
}
func (p *WorkerPool) StopWait() {
	p.stop(true)
}
func (p *WorkerPool) IsStop() bool {
	p.stopMutex.Lock()
	defer p.stopMutex.Unlock()
	return p.stopped
}
func (p *WorkerPool) Submit(task func(ctx context.Context, cancel context.CancelFunc)) {
	if task != nil {
		p.taskQueue <- task
	}
}
func (p *WorkerPool) SubmitWait(task func(ctx context.Context, cancel context.CancelFunc)) {
	if task == nil {
		return
	}
	doneChan := make(chan struct{})
	p.taskQueue <- func(ctx context.Context, cancel context.CancelFunc) {
		task(ctx, cancel)
		close(doneChan)
	}
	<-doneChan
}
func (p *WorkerPool) WaitingQueueSize() int {
	return p.waitingQueue.Len()
}
func (p *WorkerPool) dispatch() {
	defer close(p.stoppedChan)
	timeout := time.NewTimer(p.timeout)
	var (
		workerCount    int
		task           func(ctx context.Context, cancel context.CancelFunc)
		ok, wait       bool
		workerTaskChan chan func(ctx context.Context, cancel context.CancelFunc)
	)
	startReady := make(chan chan func(ctx context.Context, cancel context.CancelFunc))
Loop:
	for {
		if p.waitingQueue.Len() != 0 {
			select {
			case task, ok = <-p.taskQueue:
				if !ok {
					break Loop
				}
				if task == nil {
					wait = true
					break Loop
				}
				p.waitingQueue.PushBack(task)
			case workerTaskChan = <-p.readyWorkers:
				// A worker is ready, so give task to worker.
				workerTaskChan <- p.waitingQueue.PopFront().(func(ctx context.Context, cancel context.CancelFunc))
			}
			continue
		}
		timeout.Reset(p.timeout)
		select {
		case task, ok = <-p.taskQueue:
			if !ok || task == nil {
				break Loop
			}
			// Got a task to do.
			select {
			case workerTaskChan = <-p.readyWorkers:
				// A worker is ready, so give task to worker.
				workerTaskChan <- task
			default:
				// No workers ready.
				// Create a new worker, if not at max.
				if workerCount < p.maxWorkers {
					workerCount++
					go func(t func(ctx context.Context, cancel context.CancelFunc)) {
						startWorker(startReady, p.readyWorkers, p.tasktimeout)
						// Submit the task when the new worker.
						taskChan := <-startReady
						taskChan <- t
					}(task)
				} else {
					// Enqueue task to be executed by next available worker.
					p.waitingQueue.PushBack(task)
				}
			}
		case <-timeout.C:
			// Timed out waiting for work to arrive.  Kill a ready worker.
			if workerCount > 0 {
				select {
				case workerTaskChan = <-p.readyWorkers:
					// A worker is ready, so kill.
					close(workerTaskChan)
					workerCount--
				default:
					// No work, but no ready workers.  All workers are busy.
				}
			}
		}
	}

	// If instructed to wait for all queued tasks, then remove from queue and
	// give to workers until queue is empty.
	if wait {
		for p.waitingQueue.Len() != 0 {
			workerTaskChan = <-p.readyWorkers
			// A worker is ready, so give task to worker.
			workerTaskChan <- p.waitingQueue.PopFront().(func(ctx context.Context, cancel context.CancelFunc))
		}
	}

	// Stop all remaining workers as they become ready.
	for workerCount > 0 {
		workerTaskChan = <-p.readyWorkers
		close(workerTaskChan)
		workerCount--
	}
}
func startWorker(startReady, readyWorkers chan chan func(ctx context.Context, cancel context.CancelFunc), taskouttime time.Duration) {
	go func() {
		taskChan := make(chan func(ctx context.Context, cancel context.CancelFunc))
		var task func(ctx context.Context, cancel context.CancelFunc)
		var ok bool
		// Register availability on starReady channel.
		startReady <- taskChan
		for {
			// Read task from dispatcher.
			task, ok = <-taskChan
			if !ok {
				break
			}
			ctx, cancel := context.WithTimeout(context.Background(), taskouttime)
			go task(ctx, cancel)
			select {
			case <-ctx.Done():
				cancel()
			}
			readyWorkers <- taskChan
		}
	}()
}
func (p *WorkerPool) stop(wait bool) {
	p.stopMutex.Lock()
	defer p.stopMutex.Unlock()
	if p.stopped {
		return
	}
	p.stopped = true
	if wait {
		p.taskQueue <- nil
	}
	close(p.taskQueue)
	<-p.stoppedChan
}
