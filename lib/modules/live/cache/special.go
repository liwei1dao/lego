package cache

import "github.com/liwei1dao/lego/lib/modules/live/av"

func NewSpecialCache() *SpecialCache {
	return &SpecialCache{}
}

type SpecialCache struct {
	full bool
	p    *av.Packet
}

func (specialCache *SpecialCache) Write(p *av.Packet) {
	specialCache.p = p
	specialCache.full = true
}

func (specialCache *SpecialCache) Send(w av.WriteCloser) error {
	if !specialCache.full {
		return nil
	}
	return w.Write(specialCache.p)
}
