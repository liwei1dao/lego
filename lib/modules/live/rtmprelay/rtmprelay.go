package rtmprelay

import (
	"bytes"
	"fmt"
	"io"

	"github.com/liwei1dao/lego/lib/modules/live/amf"
	"github.com/liwei1dao/lego/lib/modules/live/av"
	"github.com/liwei1dao/lego/lib/modules/live/liveconn"
	"github.com/liwei1dao/lego/sys/log"
)

var (
	STOP_CTRL = "RTMPRELAY_STOP"
)

func NewRtmpRelay(playurl *string, publishurl *string) *RtmpRelay {
	return &RtmpRelay{
		PlayUrl:              *playurl,
		PublishUrl:           *publishurl,
		cs_chan:              make(chan liveconn.ChunkStream, 500),
		sndctrl_chan:         make(chan string),
		connectPlayClient:    nil,
		connectPublishClient: nil,
		startflag:            false,
	}
}

type RtmpRelay struct {
	PlayUrl              string
	PublishUrl           string
	cs_chan              chan liveconn.ChunkStream
	sndctrl_chan         chan string
	connectPlayClient    *liveconn.ConnClient
	connectPublishClient *liveconn.ConnClient
	startflag            bool
}

func (self *RtmpRelay) Start() error {
	if self.startflag {
		return fmt.Errorf("The rtmprelay already started, playurl=%s, publishurl=%s\n", self.PlayUrl, self.PublishUrl)
	}

	self.connectPlayClient = liveconn.NewConnClient()
	self.connectPublishClient = liveconn.NewConnClient()

	log.Debugf("play server addr:%v starting....", self.PlayUrl)
	err := self.connectPlayClient.Start(self.PlayUrl, av.PLAY)
	if err != nil {
		log.Debugf("connectPlayClient.Start url=%v error", self.PlayUrl)
		return err
	}

	log.Debugf("publish server addr:%v starting....", self.PublishUrl)
	err = self.connectPublishClient.Start(self.PublishUrl, av.PUBLISH)
	if err != nil {
		log.Debugf("connectPublishClient.Start url=%v error", self.PublishUrl)
		self.connectPlayClient.Close(nil)
		return err
	}

	self.startflag = true
	go self.rcvPlayChunkStream()
	go self.sendPublishChunkStream()

	return nil
}

func (self *RtmpRelay) Stop() {
	if !self.startflag {
		log.Debugf("The rtmprelay already stoped, playurl=%s, publishurl=%s", self.PlayUrl, self.PublishUrl)
		return
	}

	self.startflag = false
	self.sndctrl_chan <- STOP_CTRL
}

func (self *RtmpRelay) rcvPlayChunkStream() {
	log.Debug("rcvPlayRtmpMediaPacket connectClient.Read...")
	for {
		var rc liveconn.ChunkStream

		if self.startflag == false {
			self.connectPlayClient.Close(nil)
			log.Debugf("rcvPlayChunkStream close: playurl=%s, publishurl=%s", self.PlayUrl, self.PublishUrl)
			break
		}
		err := self.connectPlayClient.Read(&rc)

		if err != nil && err == io.EOF {
			break
		}
		//log.Debugf("connectPlayClient.Read return rc.TypeID=%v length=%d, err=%v", rc.TypeID, len(rc.Data), err)
		switch rc.TypeID {
		case 20, 17:
			r := bytes.NewReader(rc.Data)
			vs, err := self.connectPlayClient.DecodeBatch(r, amf.AMF0)

			log.Debugf("rcvPlayRtmpMediaPacket: vs=%v, err=%v", vs, err)
		case 18:
			log.Debug("rcvPlayRtmpMediaPacket: metadata....")
		case 8, 9:
			self.cs_chan <- rc
		}
	}
}

func (self *RtmpRelay) sendPublishChunkStream() {
	for {
		select {
		case rc := <-self.cs_chan:
			//log.Debugf("sendPublishChunkStream: rc.TypeID=%v length=%d", rc.TypeID, len(rc.Data))
			self.connectPublishClient.Write(rc)
		case ctrlcmd := <-self.sndctrl_chan:
			if ctrlcmd == STOP_CTRL {
				self.connectPublishClient.Close(nil)
				log.Debugf("sendPublishChunkStream close: playurl=%s, publishurl=%s", self.PlayUrl, self.PublishUrl)
				return
			}
		}
	}
}
