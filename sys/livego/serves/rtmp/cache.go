package rtmp

import "github.com/liwei1dao/lego/sys/livego/packet"

func NewCache(gopNum int) *Cache {
	return &Cache{
		gop:      NewGopCache(gopNum),
		videoSeq: NewSpecialCache(),
		audioSeq: NewSpecialCache(),
		metadata: NewSpecialCache(),
	}
}

type Cache struct {
	gop      *GopCache
	videoSeq *SpecialCache
	audioSeq *SpecialCache
	metadata *SpecialCache
}

func NewGopCache(num int) *GopCache {
	return &GopCache{
		count: num,
		gops:  make([]*array, num),
	}
}

type GopCache struct {
	start     bool
	num       int
	count     int
	nextindex int
	gops      []*array
}
type array struct {
	index   int
	packets []*packet.Packet
}

func NewSpecialCache() *SpecialCache {
	return &SpecialCache{}
}

type SpecialCache struct {
	full bool
	p    *packet.Packet
}
