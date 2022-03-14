package timewheel_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/liwei1dao/lego/sys/timewheel"
)

func checkTimeCost(t *testing.T, start, end time.Time, before int, after int) bool {
	due := end.Sub(start)
	if due > time.Duration(after)*time.Millisecond {
		t.Error("delay run")
		return false
	}

	if due < time.Duration(before)*time.Millisecond {
		t.Error("run ahead")
		return false
	}

	return true
}

func TestAddFunc(t *testing.T) {
	tw, _ := timewheel.NewSys(timewheel.SetTick(100), timewheel.SetBucketsNum(10))
	tw.Start()
	defer tw.Stop()

	for index := 1; index < 6; index++ {
		queue := make(chan bool, 0)
		start := time.Now()
		tw.Add(time.Duration(index)*time.Second, func(*timewheel.Task, ...interface{}) {
			queue <- true
		})

		<-queue

		before := index*1000 - 200
		after := index*1000 + 200
		checkTimeCost(t, start, time.Now(), before, after)
		fmt.Println("time since: ", time.Since(start).String())
	}
}
