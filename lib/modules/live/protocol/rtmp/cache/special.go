package cache

import "github.com/liwei1dao/lego/lib/modules/live/av"

func NewSpecialCache() *SpecialCache {
	return &SpecialCache{}
}

type SpecialCache struct {
	full bool
	p    *av.Packet
}
