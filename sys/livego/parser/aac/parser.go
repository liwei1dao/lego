package aac

import (
	"fmt"
	"io"

	"github.com/liwei1dao/lego/sys/livego/packet"
)

/*
AAC 解码
*/

func NewParser() *Parser {
	return &Parser{
		gettedSpecific: false,
		cfgInfo:        &mpegCfgInfo{},
		adtsHeader:     make([]byte, adtsHeaderLen),
	}
}

const (
	adtsHeaderLen = 7
)

var (
	specificBufInvalid = fmt.Errorf("audio mpegspecific error")
	audioBufInvalid    = fmt.Errorf("audiodata  invalid")
)

type mpegExtension struct {
	objectType byte
	sampleRate byte
}

type mpegCfgInfo struct {
	objectType     byte
	sampleRate     byte
	channel        byte
	sbr            byte
	ps             byte
	frameLen       byte
	exceptionLogTs int64
	extension      *mpegExtension
}

type Parser struct {
	gettedSpecific bool
	adtsHeader     []byte
	cfgInfo        *mpegCfgInfo
}

func (parser *Parser) Parse(b []byte, packetType uint8, w io.Writer) (err error) {
	switch packetType {
	case packet.AAC_SEQHDR:
		err = parser.specificInfo(b)
	case packet.AAC_RAW:
		err = parser.adts(b, w)
	}
	return
}

func (parser *Parser) specificInfo(src []byte) error {
	if len(src) < 2 {
		return specificBufInvalid
	}
	parser.gettedSpecific = true
	parser.cfgInfo.objectType = (src[0] >> 3) & 0xff
	parser.cfgInfo.sampleRate = ((src[0] & 0x07) << 1) | src[1]>>7
	parser.cfgInfo.channel = (src[1] >> 3) & 0x0f
	return nil
}

func (parser *Parser) adts(src []byte, w io.Writer) error {
	if len(src) <= 0 || !parser.gettedSpecific {
		return audioBufInvalid
	}

	frameLen := uint16(len(src)) + 7

	//first write adts header
	parser.adtsHeader[0] = 0xff
	parser.adtsHeader[1] = 0xf1

	parser.adtsHeader[2] &= 0x00
	parser.adtsHeader[2] = parser.adtsHeader[2] | (parser.cfgInfo.objectType-1)<<6
	parser.adtsHeader[2] = parser.adtsHeader[2] | (parser.cfgInfo.sampleRate)<<2

	parser.adtsHeader[3] &= 0x00
	parser.adtsHeader[3] = parser.adtsHeader[3] | (parser.cfgInfo.channel<<2)<<4
	parser.adtsHeader[3] = parser.adtsHeader[3] | byte((frameLen<<3)>>14)

	parser.adtsHeader[4] &= 0x00
	parser.adtsHeader[4] = parser.adtsHeader[4] | byte((frameLen<<5)>>8)

	parser.adtsHeader[5] &= 0x00
	parser.adtsHeader[5] = parser.adtsHeader[5] | byte(((frameLen<<13)>>13)<<5)
	parser.adtsHeader[5] = parser.adtsHeader[5] | (0x7C<<1)>>3
	parser.adtsHeader[6] = 0xfc

	if _, err := w.Write(parser.adtsHeader[0:]); err != nil {
		return err
	}
	if _, err := w.Write(src); err != nil {
		return err
	}
	return nil
}
