package core

import (
	"sync"
	"time"
)

type RWBaser struct {
	lock    sync.Mutex
	PreTime time.Time

	BaseTimestamp      uint32
	LastVideoTimestamp uint32
	LastAudioTimestamp uint32
}

func (this *RWBaser) BaseTimeStamp() uint32 {
	return this.BaseTimestamp
}

func (this *RWBaser) SetPreTime() {
	this.lock.Lock()
	this.PreTime = time.Now()
	this.lock.Unlock()
}

func (rw *RWBaser) RecTimeStamp(timestamp, typeID uint32) {
	if typeID == TAG_VIDEO {
		rw.LastVideoTimestamp = timestamp
	} else if typeID == TAG_AUDIO {
		rw.LastAudioTimestamp = timestamp
	}
}
