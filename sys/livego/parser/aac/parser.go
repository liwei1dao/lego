package aac

import (
	"fmt"
	"io"

	"github.com/liwei1dao/lego/sys/livego/core"
)

const (
	adtsHeaderLen = 7
)

var aacRates = []int{96000, 88200, 64000, 48000, 44100, 32000, 24000, 22050, 16000, 12000, 11025, 8000, 7350}
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

func NewParser() *Parser {
	return &Parser{
		gettedSpecific: false,
		cfgInfo:        &mpegCfgInfo{},
		adtsHeader:     make([]byte, adtsHeaderLen),
	}
}

type Parser struct {
	gettedSpecific bool
	adtsHeader     []byte
	cfgInfo        *mpegCfgInfo
}

func (this *Parser) SampleRate() int {
	rate := 44100
	if this.cfgInfo.sampleRate <= byte(len(aacRates)-1) {
		rate = aacRates[this.cfgInfo.sampleRate]
	}
	return rate
}
func (this *Parser) Parse(b []byte, packetType uint8, w io.Writer) (err error) {
	switch packetType {
	case core.AAC_SEQHDR:
		err = this.specificInfo(b)
	case core.AAC_RAW:
		err = this.adts(b, w)
	}
	return
}
func (this *Parser) specificInfo(src []byte) error {
	if len(src) < 2 {
		return specificBufInvalid
	}
	this.gettedSpecific = true
	this.cfgInfo.objectType = (src[0] >> 3) & 0xff
	this.cfgInfo.sampleRate = ((src[0] & 0x07) << 1) | src[1]>>7
	this.cfgInfo.channel = (src[1] >> 3) & 0x0f
	return nil
}
func (this *Parser) adts(src []byte, w io.Writer) error {
	if len(src) <= 0 || !this.gettedSpecific {
		return audioBufInvalid
	}

	frameLen := uint16(len(src)) + 7

	//first write adts header
	this.adtsHeader[0] = 0xff
	this.adtsHeader[1] = 0xf1

	this.adtsHeader[2] &= 0x00
	this.adtsHeader[2] = this.adtsHeader[2] | (this.cfgInfo.objectType-1)<<6
	this.adtsHeader[2] = this.adtsHeader[2] | (this.cfgInfo.sampleRate)<<2

	this.adtsHeader[3] &= 0x00
	this.adtsHeader[3] = this.adtsHeader[3] | (this.cfgInfo.channel<<2)<<4
	this.adtsHeader[3] = this.adtsHeader[3] | byte((frameLen<<3)>>14)

	this.adtsHeader[4] &= 0x00
	this.adtsHeader[4] = this.adtsHeader[4] | byte((frameLen<<5)>>8)

	this.adtsHeader[5] &= 0x00
	this.adtsHeader[5] = this.adtsHeader[5] | byte(((frameLen<<13)>>13)<<5)
	this.adtsHeader[5] = this.adtsHeader[5] | (0x7C<<1)>>3
	this.adtsHeader[6] = 0xfc

	if _, err := w.Write(this.adtsHeader[0:]); err != nil {
		return err
	}
	if _, err := w.Write(src); err != nil {
		return err
	}
	return nil
}
