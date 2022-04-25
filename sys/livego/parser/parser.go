package parser

import (
	"io"

	"github.com/liwei1dao/lego/sys/livego/packet"
	"github.com/liwei1dao/lego/sys/livego/parser/aac"
	"github.com/liwei1dao/lego/sys/livego/parser/h264"
	"github.com/liwei1dao/lego/sys/livego/parser/mp3"
)

func NewCodecParser() *CodecParser {
	return &CodecParser{}
}

type CodecParser struct {
	aac  *aac.Parser
	mp3  *mp3.Parser
	h264 *h264.Parser
}

func (codeParser *CodecParser) Parse(p *packet.Packet, w io.Writer) (err error) {

	switch p.IsVideo {
	case true:
		f, ok := p.Header.(packet.VideoPacketHeader)
		if ok {
			if f.CodecID() == packet.VIDEO_H264 {
				if codeParser.h264 == nil {
					codeParser.h264 = h264.NewParser()
				}
				err = codeParser.h264.Parse(p.Data, f.IsSeq(), w)
			}
		}
	case false:
		f, ok := p.Header.(packet.AudioPacketHeader)
		if ok {
			switch f.SoundFormat() {
			case packet.SOUND_AAC:
				if codeParser.aac == nil {
					codeParser.aac = aac.NewParser()
				}
				err = codeParser.aac.Parse(p.Data, f.AACPacketType(), w)
			case packet.SOUND_MP3:
				if codeParser.mp3 == nil {
					codeParser.mp3 = mp3.NewParser()
				}
				err = codeParser.mp3.Parse(p.Data)
			}
		}

	}
	return
}
