package timewheel

import (
	"time"

	"github.com/liwei1dao/lego/core"
)

type (
	ITimeWheel interface {
		Add(delay time.Duration, handler func(TaskID, ...interface{}), args ...interface{}) *Task
		AddCron(delay time.Duration, handler func(TaskID, ...interface{}), args ...interface{}) *Task
		Remove(task *Task) error
		NewTimer(delay time.Duration) *Timer
		NewTicker(delay time.Duration) *Ticker
		AfterFunc(delay time.Duration, callback func()) *Timer
		After(delay time.Duration) <-chan time.Time
		Sleep(delay time.Duration)
	}
)

var (
	defaultTimeWheel ITimeWheel
)

func OnInit(s core.IService, opt ...Option) (err error) {
	defaultTimeWheel, err = NewTimeWheel(opt...)
	return
}

func Add(delay time.Duration, handler func(TaskID, ...interface{}), args ...interface{}) *Task {
	return defaultTimeWheel.Add(delay, handler, args...)
}

func AddCron(delay time.Duration, handler func(TaskID, ...interface{}), args ...interface{}) *Task {
	return defaultTimeWheel.AddCron(delay, handler, args...)
}

func Remove(task *Task) error {
	return defaultTimeWheel.Remove(task)
}

func NewTimer(delay time.Duration) *Timer {
	return defaultTimeWheel.NewTimer(delay)
}

func NewTicker(delay time.Duration) *Ticker {
	return defaultTimeWheel.NewTicker(delay)
}

func AfterFunc(delay time.Duration, callback func()) *Timer {
	return defaultTimeWheel.AfterFunc(delay, callback)
}

func After(delay time.Duration) <-chan time.Time {
	return defaultTimeWheel.After(delay)
}

func Sleep(delay time.Duration) {
	defaultTimeWheel.Sleep(delay)
}
