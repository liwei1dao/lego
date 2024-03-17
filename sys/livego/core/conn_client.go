package core

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"math/rand"
	"net"
	neturl "net/url"
	"strings"

	"github.com/liwei1dao/lego/sys/livego/codec"
	"github.com/liwei1dao/lego/sys/log"
)

var (
	respResult     = "_result"
	respError      = "_error"
	onStatus       = "onStatus"
	publishStart   = "NetStream.Publish.Start"
	playStart      = "NetStream.Play.Start"
	connectSuccess = "NetConnection.Connect.Success"
	onBWDone       = "onBWDone"
)

var (
	ErrFail = fmt.Errorf("respone err")
)

func NewConnClient(sys ISys, log log.ILogger) *ConnClient {
	return &ConnClient{
		sys: sys,
		log: log,
	}
}

type ConnClient struct {
	sys        ISys
	log        log.ILogger
	transID    int
	streamid   uint32
	url        string
	tcurl      string
	app        string
	title      string
	curcmdName string
	isRTMPS    bool
	conn       *Conn
	encoder    *codec.Encoder
	decoder    *codec.Decoder
	bytesw     *bytes.Buffer
}

func (this *ConnClient) GetStreamId() uint32 {
	return this.streamid
}

func (this *ConnClient) Start(url string, method string) error {
	u, err := neturl.Parse(url)
	if err != nil {
		return err
	}
	this.url = url
	path := strings.TrimLeft(u.Path, "/")
	ps := strings.SplitN(path, "/", 2)
	if len(ps) != 2 {
		return fmt.Errorf("u path err: %s", path)
	}
	this.app = ps[0]
	this.title = ps[1]
	if u.RawQuery != "" {
		this.title += "?" + u.RawQuery
	}
	this.isRTMPS = strings.HasPrefix(url, "rtmps://")

	var port string
	if this.isRTMPS {
		this.tcurl = "rtmps://" + u.Host + "/" + this.app
		port = ":443"
	} else {
		this.tcurl = "rtmp://" + u.Host + "/" + this.app
		port = ":1935"
	}

	host := u.Host
	localIP := ":0"
	var remoteIP string
	if strings.Index(host, ":") != -1 {
		host, port, err = net.SplitHostPort(host)
		if err != nil {
			return err
		}
		port = ":" + port
	}
	ips, err := net.LookupIP(host)
	this.log.Debugf("ips: %v, host: %v", ips, host)
	if err != nil {
		this.log.Errorf("err:%v", err)
		return err
	}
	remoteIP = ips[rand.Intn(len(ips))].String()
	if strings.Index(remoteIP, ":") == -1 {
		remoteIP += port
	}

	local, err := net.ResolveTCPAddr("tcp", localIP)
	if err != nil {
		this.log.Errorf("err:%v", err)
		return err
	}
	if this.sys.GetDebug() {
		this.log.Debugf("remoteIP:%v", remoteIP)
	}
	remote, err := net.ResolveTCPAddr("tcp", remoteIP)
	if err != nil {
		this.log.Errorf("err:%v", err)
		return err
	}

	var conn net.Conn
	if this.isRTMPS {
		var config tls.Config
		if this.sys.GetEnableTLSVerify() {
			roots, err := x509.SystemCertPool()
			if err != nil {
				this.log.Errorf("err:%v", err)
				return err
			}
			config.RootCAs = roots
		} else {
			config.InsecureSkipVerify = true
		}

		conn, err = tls.Dial("tcp", remoteIP, &config)
		if err != nil {
			this.log.Errorf("err:%v", err)
			return err
		}
	} else {
		conn, err = net.DialTCP("tcp", local, remote)
		if err != nil {
			this.log.Errorf(" err:%v", err)
			return err
		}
	}

	this.log.Debugf(" connection:", "local:", conn.LocalAddr(), "remote:", conn.RemoteAddr())
	this.conn = NewConn(conn, this.sys, this.log)
	this.log.Debugf(" HandshakeClient....")
	if err := this.conn.HandshakeClient(); err != nil {
		return err
	}

	this.log.Debugf(" writeConnectMsg....")
	if err := this.writeConnectMsg(); err != nil {
		return err
	}
	this.log.Debugf(" writeCreateStreamMsg....")
	if err := this.writeCreateStreamMsg(); err != nil {
		this.log.Errorf(" writeCreateStreamMsg error", err)
		return err
	}
	this.log.Debugf(" method control:", method, PUBLISH, PLAY)
	if method == PUBLISH {
		if err := this.writePublishMsg(); err != nil {
			return err
		}
	} else if method == PLAY {
		if err := this.writePlayMsg(); err != nil {
			return err
		}
	}

	return nil
}
func (this *ConnClient) DecodeBatch(r io.Reader, ver codec.Version) (ret []interface{}, err error) {
	vs, err := this.decoder.DecodeBatch(r, ver)
	return vs, err
}
func (this *ConnClient) Read(c *ChunkStream) (err error) {
	return this.conn.Read(c)
}

func (this *ConnClient) Write(c ChunkStream) error {
	if c.TypeID == TAG_SCRIPTDATAAMF0 ||
		c.TypeID == TAG_SCRIPTDATAAMF3 {
		var err error
		if c.Data, err = codec.MetaDataReform(c.Data, codec.ADD); err != nil {
			return err
		}
		c.Length = uint32(len(c.Data))
	}
	return this.conn.Write(&c)
}

func (this *ConnClient) Close(err error) {
	this.conn.Close()
}

func (this *ConnClient) readRespMsg() error {
	var err error
	var rc ChunkStream
	for {
		if err = this.conn.Read(&rc); err != nil {
			return err
		}
		if err != nil && err != io.EOF {
			return err
		}
		switch rc.TypeID {
		case 20, 17:
			r := bytes.NewReader(rc.Data)
			vs, _ := this.decoder.DecodeBatch(r, codec.AMF0)
			this.log.Debugf("readRespMsg: vs=%v", vs)
			for k, v := range vs {
				switch v.(type) {
				case string:
					switch this.curcmdName {
					case cmdConnect, cmdCreateStream:
						if v.(string) != respResult {
							return fmt.Errorf(v.(string))
						}

					case cmdPublish:
						if v.(string) != onStatus {
							return ErrFail
						}
					}
				case float64:
					switch this.curcmdName {
					case cmdConnect, cmdCreateStream:
						id := int(v.(float64))

						if k == 1 {
							if id != this.transID {
								return ErrFail
							}
						} else if k == 3 {
							this.streamid = uint32(id)
						}
					case cmdPublish:
						if int(v.(float64)) != 0 {
							return ErrFail
						}
					}
				case codec.Object:
					objmap := v.(codec.Object)
					switch this.curcmdName {
					case cmdConnect:
						code, ok := objmap["code"]
						if ok && code.(string) != connectSuccess {
							return ErrFail
						}
					case cmdPublish:
						code, ok := objmap["code"]
						if ok && code.(string) != publishStart {
							return ErrFail
						}
					}
				}
			}
			return nil
		}
	}
}

func (this *ConnClient) writeMsg(args ...interface{}) error {
	this.bytesw.Reset()
	for _, v := range args {
		if _, err := this.encoder.Encode(this.bytesw, v, codec.AMF0); err != nil {
			return err
		}
	}
	msg := this.bytesw.Bytes()
	c := ChunkStream{
		Format:    0,
		CSID:      3,
		Timestamp: 0,
		TypeID:    20,
		StreamID:  this.streamid,
		Length:    uint32(len(msg)),
		Data:      msg,
	}
	this.conn.Write(&c)
	return this.conn.Flush()
}

func (this *ConnClient) writeConnectMsg() error {
	event := make(codec.Object)
	event["app"] = this.app
	event["type"] = "nonprivate"
	event["flashVer"] = "FMS.3.1"
	event["tcUrl"] = this.tcurl
	this.curcmdName = cmdConnect
	this.log.Debugf(" writeConnectMsg: connClient.transID=%d, event=%v", this.transID, event)
	if err := this.writeMsg(cmdConnect, this.transID, event); err != nil {
		return err
	}
	return this.readRespMsg()
}

func (this *ConnClient) writeCreateStreamMsg() error {
	this.transID++
	this.curcmdName = cmdCreateStream
	this.log.Debugf(" writeCreateStreamMsg: this.transID=%d", this.transID)
	if err := this.writeMsg(cmdCreateStream, this.transID, nil); err != nil {
		return err
	}

	for {
		err := this.readRespMsg()
		if err == nil {
			return err
		}
		if err == ErrFail {
			this.log.Errorf(" writeCreateStreamMsg readRespMsg err=%v", err)
			return err
		}
	}
}

func (this *ConnClient) writePublishMsg() error {
	this.transID++
	this.curcmdName = cmdPublish
	if err := this.writeMsg(cmdPublish, this.transID, nil, this.title, publishLive); err != nil {
		return err
	}
	return this.readRespMsg()
}

func (this *ConnClient) writePlayMsg() error {
	this.transID++
	this.curcmdName = cmdPlay
	this.log.Debugf("writePlayMsg: connClient.transID=%d, cmdPlay=%v, connClient.title=%v",
		this.transID, cmdPlay, this.title)
	if err := this.writeMsg(cmdPlay, 0, nil, this.title); err != nil {
		return err
	}
	return this.readRespMsg()
}
