package parser

import (
	"fmt"
	"io"

	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/livego/parser/aac"
	"github.com/liwei1dao/lego/sys/livego/parser/h264"
	"github.com/liwei1dao/lego/sys/livego/parser/mp3"
)

var (
	errNoAudio = fmt.Errorf("demuxer no audio")
)

func NewCodecParser() *CodecParser {
	return &CodecParser{}
}

type CodecParser struct {
	aac  *aac.Parser
	mp3  *mp3.Parser
	h264 *h264.Parser
}

func (this *CodecParser) SampleRate() (int, error) {
	if this.aac == nil && this.mp3 == nil {
		return 0, errNoAudio
	}
	if this.aac != nil {
		return this.aac.SampleRate(), nil
	}
	return this.mp3.SampleRate(), nil
}

func (this *CodecParser) Parse(p *core.Packet, w io.Writer) (err error) {

	switch p.IsVideo {
	case true:
		f, ok := p.Header.(core.VideoPacketHeader)
		if ok {
			if f.CodecID() == core.VIDEO_H264 {
				if this.h264 == nil {
					this.h264 = h264.NewParser()
				}
				err = this.h264.Parse(p.Data, f.IsSeq(), w)
			}
		}
	case false:
		f, ok := p.Header.(core.AudioPacketHeader)
		if ok {
			switch f.SoundFormat() {
			case core.SOUND_AAC:
				if this.aac == nil {
					this.aac = aac.NewParser()
				}
				err = this.aac.Parse(p.Data, f.AACPacketType(), w)
			case core.SOUND_MP3:
				if this.mp3 == nil {
					this.mp3 = mp3.NewParser()
				}
				err = this.mp3.Parse(p.Data)
			}
		}

	}
	return
}
