package timewheel

import (
	"context"
	"fmt"
	"runtime"
	"sync"
	"sync/atomic"
	"time"

	"github.com/liwei1dao/lego"
	"github.com/liwei1dao/lego/sys/log"
)

// 创建一个时间轮
func newsys(options Options) (sys *TimeWheel, err error) {
	sys = &TimeWheel{
		// tick
		tick:      options.Tick,
		tickQueue: make(chan time.Time, 10),

		// store
		bucketsNum:    options.BucketsNum,
		bucketIndexes: make(map[taskID]int, 1024*100),
		buckets:       make([]map[taskID]*Task, options.BucketsNum),
		currentIndex:  0,

		// signal
		addC:    make(chan *Task, 1024*5),
		removeC: make(chan *Task, 1024*2),
		stopC:   make(chan struct{}),
	}

	for i := 0; i < options.BucketsNum; i++ {
		sys.buckets[i] = make(map[taskID]*Task, 16)
	}

	return
}

const (
	typeTimer taskType = iota
	typeTicker

	modeIsCircle  = true
	modeNotCircle = false

	modeIsAsync  = true
	modeNotAsync = false
)

type (
	taskType int64
	taskID   int64
	Task     struct {
		delay    time.Duration
		id       taskID
		round    int
		args     []interface{}
		callback func(*Task, ...interface{})
		async    bool //异步执行
		stop     bool //是否停止
		circle   bool //是否循环
	}
	TimeWheel struct {
		randomID      int64
		tick          time.Duration
		ticker        *time.Ticker
		tickQueue     chan time.Time
		bucketsNum    int
		buckets       []map[taskID]*Task // key: added item, value: *Task
		bucketIndexes map[taskID]int     // key: added item, value: bucket position
		currentIndex  int
		onceStart     sync.Once
		addC          chan *Task
		removeC       chan *Task
		stopC         chan struct{}
		exited        bool
	}
)

// for sync.Pool
func (t *Task) Reset() {
	t.round = 0
	t.callback = nil

	t.async = false
	t.stop = false
	t.circle = false
}

// 启动时间轮
func (this *TimeWheel) Start() {
	// onlye once start
	this.onceStart.Do(
		func() {
			this.ticker = time.NewTicker(this.tick)
			go this.schduler()
			go this.tickGenerator()
		},
	)
}

func (this *TimeWheel) Add(delay time.Duration, handler func(*Task, ...interface{}), args ...interface{}) *Task {
	return this.addAny(delay, modeNotCircle, modeIsAsync, handler, args...)
}

// AddCron add interval task
func (this *TimeWheel) AddCron(delay time.Duration, handler func(*Task, ...interface{}), args ...interface{}) *Task {
	return this.addAny(delay, true, modeIsAsync, handler, args...)
}

func (this *TimeWheel) Remove(task *Task) error {
	this.removeC <- task
	return nil
}

// 停止时间轮
func (this *TimeWheel) Stop() {
	this.stopC <- struct{}{}
}

// 此处写法 为监控时间轮是否正常执行
func (this *TimeWheel) tickGenerator() {
	if this.tickQueue == nil {
		return
	}

	for !this.exited {
		select {
		case <-this.ticker.C:
			select {
			case this.tickQueue <- time.Now():
			default:
				panic("raise long time blocking")
			}
		}
	}
}

// 调度器
func (this *TimeWheel) schduler() {
	queue := this.ticker.C
	if this.tickQueue != nil {
		queue = this.tickQueue
	}

	for {
		select {
		case <-queue:
			this.handleTick()
		case task := <-this.addC:
			this.put(task)
		case key := <-this.removeC:
			this.remove(key)
		case <-this.stopC:
			this.exited = true
			this.ticker.Stop()
			return
		}
	}
}

// 清理
func (this *TimeWheel) collectTask(task *Task) {
	if index, ok := this.bucketIndexes[task.id]; ok {
		delete(this.bucketIndexes, task.id)
		delete(this.buckets[index], task.id)
	}
}

func (this *TimeWheel) recoverTask(task *Task) {
	defaultTaskPool.put(task)
}

func (this *TimeWheel) handleTick() {
	bucket := this.buckets[this.currentIndex]
	for k, task := range bucket {
		if task.stop {
			this.collectTask(task)
			this.recoverTask(task)
			continue
		}

		if bucket[k].round > 0 {
			bucket[k].round--
			continue
		}
		this.collectTask(task)
		if task.async {
			go func(_task *Task) {
				defer func() { //程序异常 收集异常信息传递给前端显示
					if r := recover(); r != nil {
						buf := make([]byte, 4096)
						l := runtime.Stack(buf, false)
						err := fmt.Errorf("%v: %s", r, buf[:l])
						log.Errorf("[timewheel] calltask err:%s", err.Error())
					}
				}()
				this.calltask(_task, _task.args...)
				if _task.circle {
					this.putCircle(_task, true)
				} else {
					this.recoverTask(_task)
				}
			}(task)
		} else {
			this.calltask(task, task.args...)
			//循环执行
			if task.circle {
				this.putCircle(task, true)
			} else {
				this.recoverTask(task)
			}
		}
	}

	if this.currentIndex == this.bucketsNum-1 {
		this.currentIndex = 0
		return
	}

	this.currentIndex++
}

// 执行时间轮事件 捕捉异常错误 防止程序崩溃
func (this *TimeWheel) calltask(task *Task, args ...interface{}) {
	defer lego.Recover("TimeWheel")
	if task.callback == nil {
		log.Error("sys.timeWheel task callback err!", log.Field{Key: "task", Value: task})
		return
	}
	task.callback(task, task.args...)
}

func (this *TimeWheel) addAny(delay time.Duration, circle, async bool, callback func(*Task, ...interface{}), agr ...interface{}) *Task {
	if delay <= 0 {
		delay = this.tick
	}

	id := this.genUniqueID()

	var task *Task
	task = defaultTaskPool.get()
	task.delay = delay
	task.id = id
	task.args = agr
	task.callback = callback
	task.circle = circle
	task.async = async // refer to src/runtime/time.go

	this.addC <- task
	return task
}

func (this *TimeWheel) put(task *Task) {
	this.store(task, false)
}

func (this *TimeWheel) putCircle(task *Task, circleMode bool) {
	this.store(task, circleMode)
}

func (this *TimeWheel) store(task *Task, circleMode bool) {
	round := this.calculateRound(task.delay)
	index := this.calculateIndex(task.delay)

	if round > 0 && circleMode {
		task.round = round - 1
	} else {
		task.round = round
	}

	this.bucketIndexes[task.id] = index
	this.buckets[index][task.id] = task
}

func (this *TimeWheel) calculateRound(delay time.Duration) (round int) {
	delaySeconds := delay.Seconds()
	tickSeconds := this.tick.Seconds()
	round = int(delaySeconds / tickSeconds / float64(this.bucketsNum))
	return
}

func (this *TimeWheel) calculateIndex(delay time.Duration) (index int) {
	delaySeconds := delay.Seconds()
	tickSeconds := this.tick.Seconds()
	index = (int(float64(this.currentIndex) + delaySeconds/tickSeconds)) % this.bucketsNum
	return
}

func (this *TimeWheel) remove(task *Task) {
	this.collectTask(task)
	this.recoverTask(task)
}

func (this *TimeWheel) NewTimer(delay time.Duration) *Timer {
	queue := make(chan bool, 1) // buf = 1, refer to src/time/sleep.go
	task := this.addAny(delay,
		modeNotCircle,
		modeNotAsync,
		func(*Task, ...interface{}) {
			notfiyChannel(queue)
		},
	)

	// init timer
	ctx, cancel := context.WithCancel(context.Background())
	timer := &Timer{
		this:   this,
		C:      queue, // faster
		task:   task,
		Ctx:    ctx,
		cancel: cancel,
	}

	return timer
}

func (this *TimeWheel) AfterFunc(delay time.Duration, callback func()) *Timer {
	queue := make(chan bool, 1)
	task := this.addAny(delay,
		modeNotCircle, modeIsAsync,
		func(*Task, ...interface{}) {
			callback()
			notfiyChannel(queue)
		},
	)

	// init timer
	ctx, cancel := context.WithCancel(context.Background())
	timer := &Timer{
		this:   this,
		C:      queue, // faster
		task:   task,
		Ctx:    ctx,
		cancel: cancel,
		fn:     callback,
	}

	return timer
}

func (this *TimeWheel) NewTicker(delay time.Duration) *Ticker {
	queue := make(chan bool, 1)
	task := this.addAny(delay,
		modeIsCircle,
		modeNotAsync,
		func(*Task, ...interface{}) {
			notfiyChannel(queue)
		},
	)

	// init ticker
	ctx, cancel := context.WithCancel(context.Background())
	ticker := &Ticker{
		task:   task,
		this:   this,
		C:      queue,
		Ctx:    ctx,
		cancel: cancel,
	}

	return ticker
}

func (this *TimeWheel) After(delay time.Duration) <-chan time.Time {
	queue := make(chan time.Time, 1)
	this.addAny(delay,
		modeNotCircle, modeNotAsync,
		func(*Task, ...interface{}) {
			queue <- time.Now()
		},
	)
	return queue
}

func (this *TimeWheel) Sleep(delay time.Duration) {
	queue := make(chan bool, 1)
	this.addAny(delay,
		modeNotCircle, modeNotAsync,
		func(*Task, ...interface{}) {
			queue <- true
		},
	)
	<-queue
}

// similar to golang std timer
type Timer struct {
	task *Task
	this *TimeWheel
	fn   func() // external custom func
	C    chan bool

	cancel context.CancelFunc
	Ctx    context.Context
}

func (t *Timer) Reset(delay time.Duration) {
	var task *Task
	if t.fn != nil { // use AfterFunc
		task = t.this.addAny(delay,
			modeNotCircle, modeIsAsync, // must async mode
			func(*Task, ...interface{}) {
				t.fn()
				notfiyChannel(t.C)
			},
		)
	} else {
		task = t.this.addAny(delay,
			modeNotCircle, modeNotAsync,
			func(*Task, ...interface{}) {
				notfiyChannel(t.C)
			},
		)
	}

	t.task = task
}

func (t *Timer) Stop() {
	t.task.stop = true
	t.cancel()
	t.this.Remove(t.task)
}

func (t *Timer) StopFunc(callback func()) {
	t.fn = callback
}

type Ticker struct {
	this   *TimeWheel
	task   *Task
	cancel context.CancelFunc

	C   chan bool
	Ctx context.Context
}

func (t *Ticker) Stop() {
	t.task.stop = true
	t.cancel()
	t.this.Remove(t.task)
}

func notfiyChannel(q chan bool) {
	select {
	case q <- true:
	default:
	}
}

func (this *TimeWheel) genUniqueID() taskID {
	id := atomic.AddInt64(&this.randomID, 1)
	return taskID(id)
}
