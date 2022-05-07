package hls

import (
	"bytes"
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

func NewTSCacheItem(id string) *TSCacheItem {
	return &TSCacheItem{
		id:  id,
		ll:  list.New(),
		num: maxTSCacheNum,
		lm:  make(map[string]TSItem),
	}
}

type TSCacheItem struct {
	id   string
	num  int
	lock sync.RWMutex
	ll   *list.List
	lm   map[string]TSItem
}

// TODO: found data race, fix it
func (this *TSCacheItem) GenM3U8PlayList() ([]byte, error) {
	var seq int
	var getSeq bool
	var maxDuration int
	m3u8body := bytes.NewBuffer(nil)
	for e := this.ll.Front(); e != nil; e = e.Next() {
		key := e.Value.(string)
		v, ok := this.lm[key]
		if ok {
			if v.Duration > maxDuration {
				maxDuration = v.Duration
			}
			if !getSeq {
				getSeq = true
				seq = v.SeqNum
			}
			fmt.Fprintf(m3u8body, "#EXTINF:%.3f,\n%s\n", float64(v.Duration)/float64(1000), v.Name)
		}
	}
	w := bytes.NewBuffer(nil)
	fmt.Fprintf(w,
		"#EXTM3U\n#EXT-X-VERSION:3\n#EXT-X-ALLOW-CACHE:NO\n#EXT-X-TARGETDURATION:%d\n#EXT-X-MEDIA-SEQUENCE:%d\n\n",
		maxDuration/1000+1, seq)
	w.Write(m3u8body.Bytes())
	return w.Bytes(), nil
}

func (this *TSCacheItem) GetItem(key string) (TSItem, error) {
	item, ok := this.lm[key]
	if !ok {
		return item, ErrNoKey
	}
	return item, nil
}

func (this *TSCacheItem) SetItem(key string, item TSItem) {
	if this.ll.Len() == this.num {
		e := this.ll.Front()
		this.ll.Remove(e)
		k := e.Value.(string)
		delete(this.lm, k)
	}
	this.lm[key] = item
	this.ll.PushBack(key)
}
