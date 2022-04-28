package core

import (
	"sync"
	"time"
)

func NewRWBaser(duration time.Duration) RWBaser {
	return RWBaser{
		timeout: duration,
		PreTime: time.Now(),
	}
}

type RWBaser struct {
	lock               sync.Mutex
	PreTime            time.Time
	timeout            time.Duration
	BaseTimestamp      uint32
	LastVideoTimestamp uint32
	LastAudioTimestamp uint32
}

func (this *RWBaser) Alive() bool {
	this.lock.Lock()
	b := !(time.Now().Sub(this.PreTime) >= this.timeout)
	this.lock.Unlock()
	return b
}

func (this *RWBaser) BaseTimeStamp() uint32 {
	return this.BaseTimestamp
}

func (this *RWBaser) SetPreTime() {
	this.lock.Lock()
	this.PreTime = time.Now()
	this.lock.Unlock()
}

func (this *RWBaser) RecTimeStamp(timestamp, typeID uint32) {
	if typeID == TAG_VIDEO {
		this.LastVideoTimestamp = timestamp
	} else if typeID == TAG_AUDIO {
		this.LastAudioTimestamp = timestamp
	}
}

func (this *RWBaser) CalcBaseTimestamp() {
	if this.LastAudioTimestamp > this.LastVideoTimestamp {
		this.BaseTimestamp = this.LastAudioTimestamp
	} else {
		this.BaseTimestamp = this.LastVideoTimestamp
	}
}
