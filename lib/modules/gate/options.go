package gate

import "time"

type (
	IGateOptions interface {
		GetGateMaxGoroutine() int
		GetTcpAddr() string
		GetWSAddr() string
		GetWSOuttime() time.Duration
		GetWSHeartbeat() time.Duration
		GetWSCertFile() string
		GetWSKeyFile() string
	}
)
