package core

import (
	"bytes"
	"fmt"
	"io"

	"github.com/liwei1dao/lego/lib/modules/live/amf"
	"github.com/liwei1dao/lego/sys/livego/codec"
	"github.com/liwei1dao/lego/sys/livego/packet"
	"github.com/liwei1dao/lego/sys/log"
)

var (
	cmdConnect       = "connect"
	cmdFcpublish     = "FCPublish"
	cmdReleaseStream = "releaseStream"
	cmdCreateStream  = "createStream"
	cmdPublish       = "publish"
	cmdFCUnpublish   = "FCUnpublish"
	cmdDeleteStream  = "deleteStream"
	cmdPlay          = "play"
)

func NewConnServer(conn *Conn) *ConnServer {
	return &ConnServer{
		conn:     conn,
		streamID: 1,
		bytesw:   bytes.NewBuffer(nil),
		decoder:  &codec.Decoder{},
		encoder:  &codec.Encoder{},
	}
}

var (
	ErrReq = fmt.Errorf("req error")
)

type ConnectInfo struct {
	App            string `amf:"app" json:"app"`
	Flashver       string `amf:"flashVer" json:"flashVer"`
	SwfUrl         string `amf:"swfUrl" json:"swfUrl"`
	TcUrl          string `amf:"tcUrl" json:"tcUrl"`
	Fpad           bool   `amf:"fpad" json:"fpad"`
	AudioCodecs    int    `amf:"audioCodecs" json:"audioCodecs"`
	VideoCodecs    int    `amf:"videoCodecs" json:"videoCodecs"`
	VideoFunction  int    `amf:"videoFunction" json:"videoFunction"`
	PageUrl        string `amf:"pageUrl" json:"pageUrl"`
	ObjectEncoding int    `amf:"objectEncoding" json:"objectEncoding"`
}

type PublishInfo struct {
	Name string
	Type string
}

type ConnServer struct {
	done          bool
	streamID      int
	isPublisher   bool
	conn          *Conn
	transactionID int
	ConnInfo      ConnectInfo
	PublishInfo   PublishInfo
	decoder       *codec.Decoder
	encoder       *codec.Encoder
	bytesw        *bytes.Buffer
}

func (this *ConnServer) handleCmdMsg(c *ChunkStream) error {
	amfType := codec.AMF0
	if c.TypeID == 17 {
		c.Data = c.Data[1:]
	}
	r := bytes.NewReader(c.Data)
	vs, err := this.decoder.DecodeBatch(r, codec.Version(amfType))
	if err != nil && err != io.EOF {
		return err
	}
	// log.Debugf("rtmp req: %#v", vs)
	switch vs[0].(type) {
	case string:
		switch vs[0].(string) {
		case cmdConnect:
			if err = this.connect(vs[1:]); err != nil {
				return err
			}
			if err = this.connectResp(c); err != nil {
				return err
			}
		case cmdCreateStream:
			if err = this.createStream(vs[1:]); err != nil {
				return err
			}
			if err = this.createStreamResp(c); err != nil {
				return err
			}
		case cmdPublish:
			if err = this.publishOrPlay(vs[1:]); err != nil {
				return err
			}
			if err = this.publishResp(c); err != nil {
				return err
			}
			this.done = true
			this.isPublisher = true
			log.Debug("handle publish req done")
		case cmdPlay:
			if err = this.publishOrPlay(vs[1:]); err != nil {
				return err
			}
			if err = this.playResp(c); err != nil {
				return err
			}
			this.done = true
			this.isPublisher = false
			log.Debug("handle play req done")
		case cmdFcpublish:
			this.fcPublish(vs)
		case cmdReleaseStream:
			this.releaseStream(vs)
		case cmdFCUnpublish:
		case cmdDeleteStream:
		default:
			log.Debugf("no support command=", vs[0].(string))
		}
	}

	return nil
}

func (this *ConnServer) ReadMsg() error {
	var c ChunkStream
	for {
		if err := this.conn.Read(&c); err != nil {
			return err
		}
		switch c.TypeID {
		case 20, 17:
			if err := this.handleCmdMsg(&c); err != nil {
				return err
			}
		}
		if this.done {
			break
		}
	}
	return nil
}

func (this *ConnServer) IsPublisher() bool {
	return this.isPublisher
}

func (this *ConnServer) Read(c *ChunkStream) (err error) {
	return this.conn.Read(c)
}

func (this *ConnServer) Write(c ChunkStream) error {
	if c.TypeID == packet.TAG_SCRIPTDATAAMF0 ||
		c.TypeID == packet.TAG_SCRIPTDATAAMF3 {
		var err error
		if c.Data, err = amf.MetaDataReform(c.Data, amf.DEL); err != nil {
			return err
		}
		c.Length = uint32(len(c.Data))
	}
	return this.conn.Write(&c)
}

func (this *ConnServer) Flush() error {
	return this.conn.Flush()
}

func (this *ConnServer) GetInfo() (app string, name string, url string) {
	app = this.ConnInfo.App
	name = this.PublishInfo.Name
	url = this.ConnInfo.TcUrl + "/" + this.PublishInfo.Name
	return
}

func (this *ConnServer) Close(err error) {
	this.conn.Close()
}

func (this *ConnServer) writeMsg(csid, streamID uint32, args ...interface{}) error {
	this.bytesw.Reset()
	for _, v := range args {
		if _, err := this.encoder.Encode(this.bytesw, v, amf.AMF0); err != nil {
			return err
		}
	}
	msg := this.bytesw.Bytes()
	c := ChunkStream{
		Format:    0,
		CSID:      csid,
		Timestamp: 0,
		TypeID:    20,
		StreamID:  streamID,
		Length:    uint32(len(msg)),
		Data:      msg,
	}
	this.conn.Write(&c)
	return this.conn.Flush()
}

func (this *ConnServer) connect(vs []interface{}) error {
	for _, v := range vs {
		switch v.(type) {
		case string:
		case float64:
			id := int(v.(float64))
			if id != 1 {
				return ErrReq
			}
			this.transactionID = id
		case codec.Object:
			obimap := v.(codec.Object)
			if app, ok := obimap["app"]; ok {
				this.ConnInfo.App = app.(string)
			}
			if flashVer, ok := obimap["flashVer"]; ok {
				this.ConnInfo.Flashver = flashVer.(string)
			}
			if tcurl, ok := obimap["tcUrl"]; ok {
				this.ConnInfo.TcUrl = tcurl.(string)
			}
			if encoding, ok := obimap["objectEncoding"]; ok {
				this.ConnInfo.ObjectEncoding = int(encoding.(float64))
			}
		}
	}
	return nil
}

func (this *ConnServer) connectResp(cur *ChunkStream) error {
	c := this.conn.NewWindowAckSize(2500000)
	this.conn.Write(&c)
	c = this.conn.NewSetPeerBandwidth(2500000)
	this.conn.Write(&c)
	c = this.conn.NewSetChunkSize(uint32(1024))
	this.conn.Write(&c)

	resp := make(amf.Object)
	resp["fmsVer"] = "FMS/3,0,1,123"
	resp["capabilities"] = 31

	event := make(amf.Object)
	event["level"] = "status"
	event["code"] = "NetConnection.Connect.Success"
	event["description"] = "Connection succeeded."
	event["objectEncoding"] = this.ConnInfo.ObjectEncoding
	return this.writeMsg(cur.CSID, cur.StreamID, "_result", this.transactionID, resp, event)
}

func (this *ConnServer) createStream(vs []interface{}) error {
	for _, v := range vs {
		switch v.(type) {
		case string:
		case float64:
			this.transactionID = int(v.(float64))
		case amf.Object:
		}
	}
	return nil
}

func (this *ConnServer) createStreamResp(cur *ChunkStream) error {
	return this.writeMsg(cur.CSID, cur.StreamID, "_result", this.transactionID, nil, this.streamID)
}

func (this *ConnServer) publishOrPlay(vs []interface{}) error {
	for k, v := range vs {
		switch v.(type) {
		case string:
			if k == 2 {
				this.PublishInfo.Name = v.(string)
			} else if k == 3 {
				this.PublishInfo.Type = v.(string)
			}
		case float64:
			id := int(v.(float64))
			this.transactionID = id
		case amf.Object:
		}
	}

	return nil
}

func (this *ConnServer) publishResp(cur *ChunkStream) error {
	event := make(amf.Object)
	event["level"] = "status"
	event["code"] = "NetStream.Publish.Start"
	event["description"] = "Start publishing."
	return this.writeMsg(cur.CSID, cur.StreamID, "onStatus", 0, nil, event)
}

func (this *ConnServer) playResp(cur *ChunkStream) error {
	this.conn.SetRecorded()
	this.conn.SetBegin()

	event := make(amf.Object)
	event["level"] = "status"
	event["code"] = "NetStream.Play.Reset"
	event["description"] = "Playing and resetting stream."
	if err := this.writeMsg(cur.CSID, cur.StreamID, "onStatus", 0, nil, event); err != nil {
		return err
	}

	event["level"] = "status"
	event["code"] = "NetStream.Play.Start"
	event["description"] = "Started playing stream."
	if err := this.writeMsg(cur.CSID, cur.StreamID, "onStatus", 0, nil, event); err != nil {
		return err
	}

	event["level"] = "status"
	event["code"] = "NetStream.Data.Start"
	event["description"] = "Started playing stream."
	if err := this.writeMsg(cur.CSID, cur.StreamID, "onStatus", 0, nil, event); err != nil {
		return err
	}

	event["level"] = "status"
	event["code"] = "NetStream.Play.PublishNotify"
	event["description"] = "Started playing notify."
	if err := this.writeMsg(cur.CSID, cur.StreamID, "onStatus", 0, nil, event); err != nil {
		return err
	}
	return this.conn.Flush()
}

func (this *ConnServer) fcPublish(vs []interface{}) error {
	return nil
}

func (this *ConnServer) releaseStream(vs []interface{}) error {
	return nil
}
