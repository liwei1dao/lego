package core

import (
	"fmt"
)

var (
	STATIC_RELAY_STOP_CTRL = "STATIC_RTMPRELAY_STOP"
)

func NewStaticPush(rtmpurl string) *StaticPush {
	return &StaticPush{
		RtmpUrl:       rtmpurl,
		packet_chan:   make(chan *Packet, 500),
		sndctrl_chan:  make(chan string),
		connectClient: nil,
		startflag:     false,
	}
}

type StaticPush struct {
	server        IServer
	RtmpUrl       string
	packet_chan   chan *Packet
	sndctrl_chan  chan string
	connectClient *ConnClient
	startflag     bool
}

func (this *StaticPush) IsStart() bool {
	return this.startflag
}

func (this *StaticPush) Start() error {
	if this.startflag {
		return fmt.Errorf("StaticPush already start %s", this.RtmpUrl)
	}

	this.connectClient = NewConnClient(this.server)
	if this.server.GetDebug() {
		this.server.Debugf("[SYS LiveGo] static publish server addr:%v starting....", this.RtmpUrl)
	}
	err := this.connectClient.Start(this.RtmpUrl, "publish")
	if err != nil {
		this.server.Debugf("connectClient.Start url=%v error", this.RtmpUrl)
		return err
	}
	if this.server.GetDebug() {
		this.server.Debugf("[SYS LiveGo] static publish server addr:%v started, streamid=%d", this.RtmpUrl, this.connectClient.GetStreamId())
	}
	go this.HandleAvPacket()

	this.startflag = true
	return nil
}

func (this *StaticPush) HandleAvPacket() {
	if !this.IsStart() {
		this.server.Debugf("static push %s not started", this.RtmpUrl)
		return
	}

	for {
		select {
		case packet := <-this.packet_chan:
			this.sendPacket(packet)
		case ctrlcmd := <-this.sndctrl_chan:
			if ctrlcmd == STATIC_RELAY_STOP_CTRL {
				this.connectClient.Close(nil)
				this.server.Debugf("Static HandleAvPacket close: publishurl=%s", this.RtmpUrl)
				return
			}
		}
	}
}

func (this *StaticPush) Stop() {
	if !this.startflag {
		return
	}

	this.server.Debugf("StaticPush Stop: %s", this.RtmpUrl)
	this.sndctrl_chan <- STATIC_RELAY_STOP_CTRL
	this.startflag = false
}

func (this *StaticPush) sendPacket(p *Packet) {
	if !this.startflag {
		return
	}
	var cs ChunkStream

	cs.Data = p.Data
	cs.Length = uint32(len(p.Data))
	cs.StreamID = this.connectClient.GetStreamId()
	cs.Timestamp = p.TimeStamp
	if p.IsVideo {
		cs.TypeID = TAG_VIDEO
	} else {
		if p.IsMetadata {
			cs.TypeID = TAG_SCRIPTDATAAMF0
		} else {
			cs.TypeID = TAG_AUDIO
		}
	}
	this.connectClient.Write(cs)
}

func (this *StaticPush) WriteAvPacket(packet *Packet) {
	if !this.startflag {
		return
	}
	this.packet_chan <- packet
}
