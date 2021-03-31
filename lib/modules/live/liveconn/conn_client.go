package liveconn

import (
	"bytes"

	"github.com/liwei1dao/lego/lib/modules/live/amf"
)

func NewConnClient() *ConnClient {
	return &ConnClient{
		transID: 1,
		bytesw:  bytes.NewBuffer(nil),
		encoder: &amf.Encoder{},
		decoder: &amf.Decoder{},
	}
}

type ConnClient struct {
	done       bool
	transID    int
	url        string
	tcurl      string
	app        string
	title      string
	query      string
	curcmdName string
	streamid   uint32
	conn       *Conn
	encoder    *amf.Encoder
	decoder    *amf.Decoder
	bytesw     *bytes.Buffer
}
