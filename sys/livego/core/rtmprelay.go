package core

import (
	"bytes"
	"fmt"
	"io"

	"github.com/liwei1dao/lego/sys/livego/codec"
)

var (
	STOP_CTRL = "RTMPRELAY_STOP"
)

func NewRtmpRelay(sys ISys, playurl string, publishurl string) *RtmpRelay {
	return &RtmpRelay{
		sys:                  sys,
		PlayUrl:              playurl,
		PublishUrl:           publishurl,
		cs_chan:              make(chan ChunkStream, 500),
		sndctrl_chan:         make(chan string),
		connectPlayClient:    nil,
		connectPublishClient: nil,
		startflag:            false,
	}
}

type RtmpRelay struct {
	sys                  ISys
	PlayUrl              string
	PublishUrl           string
	cs_chan              chan ChunkStream
	sndctrl_chan         chan string
	connectPlayClient    *ConnClient
	connectPublishClient *ConnClient
	startflag            bool
}

func (this *RtmpRelay) Start() error {
	if this.startflag {
		return fmt.Errorf("The rtmprelay already started, playurl=%s, publishurl=%s\n", this.PlayUrl, this.PublishUrl)
	}

	this.connectPlayClient = NewConnClient(this.sys)
	this.connectPublishClient = NewConnClient(this.sys)

	this.sys.Debugf("play server addr:%v starting....", this.PlayUrl)
	err := this.connectPlayClient.Start(this.PlayUrl, PLAY)
	if err != nil {
		this.sys.Debugf("connectPlayClient.Start url=%v error", this.PlayUrl)
		return err
	}

	this.sys.Debugf("publish server addr:%v starting....", this.PublishUrl)
	err = this.connectPublishClient.Start(this.PublishUrl, PUBLISH)
	if err != nil {
		this.sys.Debugf("connectPublishClient.Start url=%v error", this.PublishUrl)
		this.connectPlayClient.Close(nil)
		return err
	}

	this.startflag = true
	go this.rcvPlayChunkStream()
	go this.sendPublishChunkStream()

	return nil
}

func (this *RtmpRelay) Stop() {
	if !this.startflag {
		this.sys.Debugf("The rtmprelay already stoped, playurl=%s, publishurl=%s", this.PlayUrl, this.PublishUrl)
		return
	}
	this.startflag = false
	this.sndctrl_chan <- STOP_CTRL
}

func (this *RtmpRelay) rcvPlayChunkStream() {
	this.sys.Debugf("rcvPlayRtmpMediaPacket connectClient.Read...")
	for {
		var rc ChunkStream

		if this.startflag == false {
			this.connectPlayClient.Close(nil)
			this.sys.Debugf("rcvPlayChunkStream close: playurl=%s, publishurl=%s", this.PlayUrl, this.PublishUrl)
			break
		}
		err := this.connectPlayClient.Read(&rc)

		if err != nil && err == io.EOF {
			break
		}
		//log.Debugf("connectPlayClient.Read return rc.TypeID=%v length=%d, err=%v", rc.TypeID, len(rc.Data), err)
		switch rc.TypeID {
		case 20, 17:
			r := bytes.NewReader(rc.Data)
			vs, err := this.connectPlayClient.DecodeBatch(r, codec.AMF0)

			this.sys.Debugf("rcvPlayRtmpMediaPacket: vs=%v, err=%v", vs, err)
		case 18:
			this.sys.Debugf("rcvPlayRtmpMediaPacket: metadata....")
			this.cs_chan <- rc
		case 8, 9:
			this.cs_chan <- rc
		}
	}
}

func (this *RtmpRelay) sendPublishChunkStream() {
	for {
		select {
		case rc := <-this.cs_chan:
			//log.Debugf("sendPublishChunkStream: rc.TypeID=%v length=%d", rc.TypeID, len(rc.Data))
			this.connectPublishClient.Write(rc)
		case ctrlcmd := <-this.sndctrl_chan:
			if ctrlcmd == STOP_CTRL {
				this.connectPublishClient.Close(nil)
				this.sys.Debugf("sendPublishChunkStream close: playurl=%s, publishurl=%s", this.PlayUrl, this.PublishUrl)
				return
			}
		}
	}
}
