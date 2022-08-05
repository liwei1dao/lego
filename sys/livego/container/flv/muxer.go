package flv

import (
	"fmt"
	"os"
	"path"
	"strings"
	"time"

	"github.com/liwei1dao/lego/sys/livego/codec"
	"github.com/liwei1dao/lego/sys/livego/core"
	"github.com/liwei1dao/lego/sys/livego/utils/pio"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/utils/container/id"
)

var (
	flvHeader = []byte{0x46, 0x4c, 0x56, 0x01, 0x05, 0x00, 0x00, 0x00, 0x09}
)

const (
	headerLen = 11
)

func NewFlvDvr(sys core.ISys, log log.ILogger) *FlvDvr {
	return &FlvDvr{
		sys: sys,
		log: log,
	}
}

type FlvDvr struct {
	sys core.ISys
	log log.ILogger
}

func (this *FlvDvr) GetWriter(info core.Info) core.WriteCloser {
	paths := strings.SplitN(info.Key, "/", 2)
	if len(paths) != 2 {
		this.log.Warnf("invalid info")
		return nil
	}

	flvDir := this.sys.GetFLVDir()

	err := os.MkdirAll(path.Join(flvDir, paths[0]), 0755)
	if err != nil {
		this.log.Errorf("mkdir error: ", err)
		return nil
	}

	fileName := fmt.Sprintf("%s_%d.%s", path.Join(flvDir, info.Key), time.Now().Unix(), "flv")
	this.log.Debugf("flv dvr save stream to: ", fileName)
	w, err := os.OpenFile(fileName, os.O_CREATE|os.O_RDWR, 0755)
	if err != nil {
		this.log.Errorf("open file error: ", err)
		return nil
	}

	writer := NewFLVWriter(paths[0], paths[1], info.URL, w)
	this.log.Debugf("new flv dvr: ", writer.Info())
	return writer
}

func NewFLVWriter(app, title, url string, ctx *os.File) *FLVWriter {
	ret := &FLVWriter{
		Uid:     id.NewXId(),
		app:     app,
		title:   title,
		url:     url,
		ctx:     ctx,
		RWBaser: core.NewRWBaser(time.Second * 10),
		closed:  make(chan struct{}),
		buf:     make([]byte, headerLen),
	}

	ret.ctx.Write(flvHeader)
	pio.PutI32BE(ret.buf[:4], 0)
	ret.ctx.Write(ret.buf[:4])

	return ret
}

type FLVWriter struct {
	Uid string
	core.RWBaser
	app, title, url string
	buf             []byte
	closed          chan struct{}
	ctx             *os.File
	closedWriter    bool
}

func (this *FLVWriter) Info() (ret core.Info) {
	ret.UID = this.Uid
	ret.URL = this.url
	ret.Key = this.app + "/" + this.title
	return
}

func (this *FLVWriter) Write(p *core.Packet) error {
	this.RWBaser.SetPreTime()
	h := this.buf[:headerLen]
	typeID := core.TAG_VIDEO
	if !p.IsVideo {
		if p.IsMetadata {
			var err error
			typeID = core.TAG_SCRIPTDATAAMF0
			p.Data, err = codec.MetaDataReform(p.Data, codec.DEL)
			if err != nil {
				return err
			}
		} else {
			typeID = core.TAG_AUDIO
		}
	}
	dataLen := len(p.Data)
	timestamp := p.TimeStamp
	timestamp += this.BaseTimeStamp()
	this.RWBaser.RecTimeStamp(timestamp, uint32(typeID))

	preDataLen := dataLen + headerLen
	timestampbase := timestamp & 0xffffff
	timestampExt := timestamp >> 24 & 0xff

	pio.PutU8(h[0:1], uint8(typeID))
	pio.PutI24BE(h[1:4], int32(dataLen))
	pio.PutI24BE(h[4:7], int32(timestampbase))
	pio.PutU8(h[7:8], uint8(timestampExt))

	if _, err := this.ctx.Write(h); err != nil {
		return err
	}

	if _, err := this.ctx.Write(p.Data); err != nil {
		return err
	}

	pio.PutI32BE(h[:4], int32(preDataLen))
	if _, err := this.ctx.Write(h[:4]); err != nil {
		return err
	}
	return nil
}

func (this *FLVWriter) Wait() {
	select {
	case <-this.closed:
		return
	}
}
func (this *FLVWriter) Close(error) {
	if this.closedWriter {
		return
	}
	this.closedWriter = true
	this.ctx.Close()
	close(this.closed)
}
