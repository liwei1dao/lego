package cache

import "github.com/liwei1dao/lego/lib/modules/live/av"

type array struct {
	index   int
	packets []*av.Packet
}

type GopCache struct {
	start     bool
	num       int
	count     int
	nextindex int
	gops      []*array
}

func NewGopCache(num int) *GopCache {
	return &GopCache{
		count: num,
		gops:  make([]*array, num),
	}
}
