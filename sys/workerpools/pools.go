package workerpools

import (
	"context"
	"sync"
	"time"

	"github.com/liwei1dao/lego/utils/container"
)

func newSys(options Options) (sys *WorkerPool, err error) {
	sys = &WorkerPool{
		taskQueue:    make(chan *Task, 1),
		maxWorkers:   options.MaxWorkers,
		readyWorkers: make(chan chan *Task, options.DefWrokers),
		timeout:      options.IdleTimeoutSec,
		tasktimeout:  options.Tasktimeout,
		stoppedChan:  make(chan struct{}),
	}
	go sys.dispatch()
	return
}

type Task struct {
	f    func(ctx context.Context, cancel context.CancelFunc, agrs ...interface{})
	agrs []interface{}
}

type WorkerPool struct {
	maxWorkers   int
	timeout      time.Duration //超时释放空闲工作人员
	tasktimeout  time.Duration //任务执行操超时间
	taskQueue    chan *Task
	readyWorkers chan chan *Task
	stoppedChan  chan struct{}
	waitingQueue container.Deque
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
func (p *WorkerPool) Submit(task func(ctx context.Context, cancel context.CancelFunc, agrs ...interface{}), agrs ...interface{}) {
	if task != nil {
		p.taskQueue <- &Task{f: task, agrs: agrs}
	}
}
func (p *WorkerPool) SubmitWait(task func(ctx context.Context, cancel context.CancelFunc, agrs ...interface{}), agrs ...interface{}) {
	if task == nil {
		return
	}
	doneChan := make(chan struct{})
	p.taskQueue <- &Task{f: func(ctx context.Context, cancel context.CancelFunc, agrs ...interface{}) {
		task(ctx, cancel, agrs...)
		close(doneChan)
	}, agrs: agrs}
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
		task           *Task
		ok, wait       bool
		workerTaskChan chan *Task
	)
	startReady := make(chan chan *Task)
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
				workerTaskChan <- p.waitingQueue.PopFront().(*Task)
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
					go func(t *Task) {
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
			workerTaskChan <- p.waitingQueue.PopFront().(*Task)
		}
	}

	// Stop all remaining workers as they become ready.
	for workerCount > 0 {
		workerTaskChan = <-p.readyWorkers
		close(workerTaskChan)
		workerCount--
	}
}
func startWorker(startReady, readyWorkers chan chan *Task, taskouttime time.Duration) {
	go func() {
		taskChan := make(chan *Task)
		var task *Task
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
			go task.f(ctx, cancel, task.agrs...)
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
