package live

import (
	"container/list"
	"fmt"
	"sync"
)

const (
	maxTSCacheNum = 3
)

var (
	ErrNoKey = fmt.Errorf("No key for cache")
)

type TSCacheItem struct {
	id   string
	num  int
	lock sync.RWMutex
	ll   *list.List
	lm   map[string]TSItem
}
