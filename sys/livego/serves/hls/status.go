package hls

import "time"

func newStatus() *status {
	return &status{
		seqId:         0,
		hasSetFirstTs: false,
		segBeginAt:    time.Now(),
	}
}

type status struct {
	hasVideo       bool
	seqId          int64
	createdAt      time.Time
	segBeginAt     time.Time
	hasSetFirstTs  bool
	firstTimestamp int64
	lastTimestamp  int64
}

func (this *status) update(isVideo bool, timestamp uint32) {
	if isVideo {
		this.hasVideo = true
	}
	if !this.hasSetFirstTs {
		this.hasSetFirstTs = true
		this.firstTimestamp = int64(timestamp)
	}
	this.lastTimestamp = int64(timestamp)
}
func (this *status) resetAndNew() {
	this.seqId++
	this.hasVideo = false
	this.createdAt = time.Now()
	this.hasSetFirstTs = false
}

func (this *status) durationMs() int64 {
	return this.lastTimestamp - this.firstTimestamp
}
