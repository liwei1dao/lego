package live

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
