package timewheel

import (
	"time"
)

type (
	ISys interface {
		Start()
		Stop()
		Add(delay time.Duration, handler func(*Task, ...interface{}), args ...interface{}) *Task
		AddCron(delay time.Duration, handler func(*Task, ...interface{}), args ...interface{}) *Task
		Remove(task *Task) error
		NewTimer(delay time.Duration) *Timer
		NewTicker(delay time.Duration) *Ticker
		AfterFunc(delay time.Duration, callback func()) *Timer
		After(delay time.Duration) <-chan time.Time
		Sleep(delay time.Duration)
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	if defsys, err = newsys(newOptions(config, option...)); err == nil {
		defsys.Start()
	}
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	if sys, err = newsys(newOptionsByOption(option...)); err == nil {
		sys.Start()
	}
	return
}

func Add(delay time.Duration, handler func(*Task, ...interface{}), args ...interface{}) *Task {
	return defsys.Add(delay, handler, args...)
}

func AddCron(delay time.Duration, handler func(*Task, ...interface{}), args ...interface{}) *Task {
	return defsys.AddCron(delay, handler, args...)
}

func Remove(task *Task) error {
	return defsys.Remove(task)
}

func NewTimer(delay time.Duration) *Timer {
	return defsys.NewTimer(delay)
}

func NewTicker(delay time.Duration) *Ticker {
	return defsys.NewTicker(delay)
}

func AfterFunc(delay time.Duration, callback func()) *Timer {
	return defsys.AfterFunc(delay, callback)
}

func After(delay time.Duration) <-chan time.Time {
	return defsys.After(delay)
}

func Sleep(delay time.Duration) {
	defsys.Sleep(delay)
}
