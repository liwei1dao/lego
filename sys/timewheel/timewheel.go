package timewheel

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// 创建一个时间轮
func newsys(options Options) (sys *TimeWheel, err error) {
	sys = &TimeWheel{
		// tick
		tick:      time.Millisecond * time.Duration(options.Tick),
		tickQueue: make(chan time.Time, 10),

		// store
		bucketsNum:    options.BucketsNum,
		bucketIndexes: make(map[taskID]int, 1024*100),
		buckets:       make([]map[taskID]*Task, options.BucketsNum),
		currentIndex:  0,

		// signal
		addC:     make(chan *Task, 1024*5),
		removeC:  make(chan *Task, 1024*2),
		stopC:    make(chan struct{}),
		syncPool: options.IsSyncPool,
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
		async    bool
		stop     bool
		circle   bool
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
		syncPool      bool
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

//启动时间轮
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

func (tw *TimeWheel) Add(delay time.Duration, handler func(*Task, ...interface{}), args ...interface{}) *Task {
	return tw.addAny(delay, modeNotCircle, modeIsAsync, handler, args...)
}

// AddCron add interval task
func (tw *TimeWheel) AddCron(delay time.Duration, handler func(*Task, ...interface{}), args ...interface{}) *Task {
	return tw.addAny(delay, modeIsCircle, modeIsAsync, handler, args...)
}

func (tw *TimeWheel) Remove(task *Task) error {
	tw.removeC <- task
	return nil
}

//停止时间轮
func (tw *TimeWheel) Stop() {
	tw.stopC <- struct{}{}
}

func (tw *TimeWheel) tickGenerator() {
	if tw.tickQueue != nil {
		return
	}

	for !tw.exited {
		select {
		case <-tw.ticker.C:
			select {
			case tw.tickQueue <- time.Now():
			default:
				panic("raise long time blocking")
			}
		}
	}
}

//d调度器
func (this *TimeWheel) schduler() {
	queue := this.ticker.C
	if this.tickQueue == nil {
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

//清理
func (tw *TimeWheel) collectTask(task *Task) {
	index := tw.bucketIndexes[task.id]
	delete(tw.bucketIndexes, task.id)
	delete(tw.buckets[index], task.id)

	if tw.syncPool && !task.circle {
		defaultTaskPool.put(task)
	}
}

func (tw *TimeWheel) handleTick() {
	bucket := tw.buckets[tw.currentIndex]
	for k, task := range bucket {
		if task.stop {
			tw.collectTask(task)
			continue
		}

		if bucket[k].round > 0 {
			bucket[k].round--
			continue
		}

		if task.async {
			go task.callback(task, task.args...)
		} else {
			// optimize gopool
			task.callback(task, task.args...)
		}

		// circle
		if task.circle == true {
			tw.collectTask(task)
			tw.putCircle(task, modeIsCircle)
			continue
		}

		// gc
		tw.collectTask(task)
	}

	if tw.currentIndex == tw.bucketsNum-1 {
		tw.currentIndex = 0
		return
	}

	tw.currentIndex++
}

func (tw *TimeWheel) addAny(delay time.Duration, circle, async bool, callback func(*Task, ...interface{}), agr ...interface{}) *Task {
	if delay <= 0 {
		delay = tw.tick
	}

	id := tw.genUniqueID()

	var task *Task
	if tw.syncPool {
		task = defaultTaskPool.get()
	} else {
		task = new(Task)
	}

	task.delay = delay
	task.id = id
	task.args = agr
	task.callback = callback
	task.circle = circle
	task.async = async // refer to src/runtime/time.go

	tw.addC <- task
	return task
}

func (tw *TimeWheel) put(task *Task) {
	tw.store(task, false)
}

func (tw *TimeWheel) putCircle(task *Task, circleMode bool) {
	tw.store(task, circleMode)
}

func (tw *TimeWheel) store(task *Task, circleMode bool) {
	round := tw.calculateRound(task.delay)
	index := tw.calculateIndex(task.delay)

	if round > 0 && circleMode {
		task.round = round - 1
	} else {
		task.round = round
	}

	tw.bucketIndexes[task.id] = index
	tw.buckets[index][task.id] = task
}

func (tw *TimeWheel) calculateRound(delay time.Duration) (round int) {
	delaySeconds := delay.Seconds()
	tickSeconds := tw.tick.Seconds()
	round = int(delaySeconds / tickSeconds / float64(tw.bucketsNum))
	return
}

func (tw *TimeWheel) calculateIndex(delay time.Duration) (index int) {
	delaySeconds := delay.Seconds()
	tickSeconds := tw.tick.Seconds()
	index = (int(float64(tw.currentIndex) + delaySeconds/tickSeconds)) % tw.bucketsNum
	return
}

func (tw *TimeWheel) remove(task *Task) {
	tw.collectTask(task)
}

func (tw *TimeWheel) NewTimer(delay time.Duration) *Timer {
	queue := make(chan bool, 1) // buf = 1, refer to src/time/sleep.go
	task := tw.addAny(delay,
		modeNotCircle,
		modeNotAsync,
		func(*Task, ...interface{}) {
			notfiyChannel(queue)
		},
	)

	// init timer
	ctx, cancel := context.WithCancel(context.Background())
	timer := &Timer{
		tw:     tw,
		C:      queue, // faster
		task:   task,
		Ctx:    ctx,
		cancel: cancel,
	}

	return timer
}

func (tw *TimeWheel) AfterFunc(delay time.Duration, callback func()) *Timer {
	queue := make(chan bool, 1)
	task := tw.addAny(delay,
		modeNotCircle, modeIsAsync,
		func(*Task, ...interface{}) {
			callback()
			notfiyChannel(queue)
		},
	)

	// init timer
	ctx, cancel := context.WithCancel(context.Background())
	timer := &Timer{
		tw:     tw,
		C:      queue, // faster
		task:   task,
		Ctx:    ctx,
		cancel: cancel,
		fn:     callback,
	}

	return timer
}

func (tw *TimeWheel) NewTicker(delay time.Duration) *Ticker {
	queue := make(chan bool, 1)
	task := tw.addAny(delay,
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
		tw:     tw,
		C:      queue,
		Ctx:    ctx,
		cancel: cancel,
	}

	return ticker
}

func (tw *TimeWheel) After(delay time.Duration) <-chan time.Time {
	queue := make(chan time.Time, 1)
	tw.addAny(delay,
		modeNotCircle, modeNotAsync,
		func(*Task, ...interface{}) {
			queue <- time.Now()
		},
	)
	return queue
}

func (tw *TimeWheel) Sleep(delay time.Duration) {
	queue := make(chan bool, 1)
	tw.addAny(delay,
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
	tw   *TimeWheel
	fn   func() // external custom func
	C    chan bool

	cancel context.CancelFunc
	Ctx    context.Context
}

func (t *Timer) Reset(delay time.Duration) {
	var task *Task
	if t.fn != nil { // use AfterFunc
		task = t.tw.addAny(delay,
			modeNotCircle, modeIsAsync, // must async mode
			func(*Task, ...interface{}) {
				t.fn()
				notfiyChannel(t.C)
			},
		)
	} else {
		task = t.tw.addAny(delay,
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
	t.tw.Remove(t.task)
}

func (t *Timer) StopFunc(callback func()) {
	t.fn = callback
}

type Ticker struct {
	tw     *TimeWheel
	task   *Task
	cancel context.CancelFunc

	C   chan bool
	Ctx context.Context
}

func (t *Ticker) Stop() {
	t.task.stop = true
	t.cancel()
	t.tw.Remove(t.task)
}

func notfiyChannel(q chan bool) {
	select {
	case q <- true:
	default:
	}
}

func (tw *TimeWheel) genUniqueID() taskID {
	id := atomic.AddInt64(&tw.randomID, 1)
	return taskID(id)
}
