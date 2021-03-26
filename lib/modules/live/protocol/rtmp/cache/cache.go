package cache

import "github.com/liwei1dao/lego/lib/modules/live/configure"

func NewCache() *Cache {
	return &Cache{
		gop:      NewGopCache(configure.Config.GetInt("gop_num")),
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
