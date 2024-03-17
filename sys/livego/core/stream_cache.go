package core

import (
	"fmt"
)

var (
	maxGOPCap    int = 1024
	ErrGopTooBig     = fmt.Errorf("gop to big")
)

func NewCache() *Cache {
	return &Cache{}
}

type Cache struct {
	gop      *GopCache
	videoSeq *SpecialCache
	audioSeq *SpecialCache
	metadata *SpecialCache
}

func (this *Cache) Write(p Packet) {
	if p.IsMetadata {
		this.metadata.Write(&p)
		return
	} else {
		if !p.IsVideo {
			ah, ok := p.Header.(AudioPacketHeader)
			if ok {
				if ah.SoundFormat() == SOUND_AAC &&
					ah.AACPacketType() == AAC_SEQHDR {
					this.audioSeq.Write(&p)
					return
				} else {
					return
				}
			}

		} else {
			vh, ok := p.Header.(VideoPacketHeader)
			if ok {
				if vh.IsSeq() {
					this.videoSeq.Write(&p)
					return
				}
			} else {
				return
			}

		}
	}
	this.gop.Write(&p)
}

func (this *Cache) Send(w WriteCloser) error {
	if err := this.metadata.Send(w); err != nil {
		return err
	}

	if err := this.videoSeq.Send(w); err != nil {
		return err
	}

	if err := this.audioSeq.Send(w); err != nil {
		return err
	}

	if err := this.gop.Send(w); err != nil {
		return err
	}
	return nil
}

func newArray() *array {
	ret := &array{
		index:   0,
		packets: make([]*Packet, 0, maxGOPCap),
	}
	return ret
}

type array struct {
	index   int
	packets []*Packet
}

func (this *array) reset() {
	this.index = 0
	this.packets = this.packets[:0]
}

func (this *array) write(packet *Packet) error {
	if this.index >= maxGOPCap {
		return ErrGopTooBig
	}
	this.packets = append(this.packets, packet)
	this.index++
	return nil
}
func (this *array) send(w WriteCloser) error {
	var err error
	for i := 0; i < this.index; i++ {
		packet := this.packets[i]
		if err = w.Write(packet); err != nil {
			return err
		}
	}
	return err
}

type GopCache struct {
	start     bool
	num       int
	count     int
	nextindex int
	gops      []*array
}

func (this *GopCache) Write(p *Packet) {
	var ok bool
	if p.IsVideo {
		vh := p.Header.(VideoPacketHeader)
		if vh.IsKeyFrame() && !vh.IsSeq() {
			ok = true
		}
	}
	if ok || this.start {
		this.start = true
		this.writeToArray(p, ok)
	}
}

func (this *GopCache) Send(w WriteCloser) error {
	return this.sendTo(w)
}

func (this *GopCache) writeToArray(chunk *Packet, startNew bool) error {
	var ginc *array
	if startNew {
		ginc = this.gops[this.nextindex]
		if ginc == nil {
			ginc = newArray()
			this.num++
			this.gops[this.nextindex] = ginc
		} else {
			ginc.reset()
		}
		this.nextindex = (this.nextindex + 1) % this.count
	} else {
		ginc = this.gops[(this.nextindex+1)%this.count]
	}
	ginc.write(chunk)
	return nil
}

func (this *GopCache) sendTo(w WriteCloser) error {
	var err error
	pos := (this.nextindex + 1) % this.count
	for i := 0; i < this.num; i++ {
		index := (pos - this.num + 1) + i
		if index < 0 {
			index += this.count
		}
		g := this.gops[index]
		err = g.send(w)
		if err != nil {
			return err
		}
	}
	return nil
}

type SpecialCache struct {
	full bool
	p    *Packet
}

func (this *SpecialCache) Write(p *Packet) {
	this.p = p
	this.full = true
}

func (this *SpecialCache) Send(w WriteCloser) error {
	if !this.full {
		return nil
	}

	// demux in hls will change p.Data, only send a copy here
	newPacket := *this.p
	return w.Write(&newPacket)
}
