package protocol

import (
	"sync"

	"github.com/liwei1dao/lego/core"
)

var msgPool = sync.Pool{
	New: func() interface{} {
		header := Header([12]byte{})
		header[0] = magicNumber

		return &Message{
			Header: &header,
			from:   &core.ServiceNode{},
		}
	},
}

func GetPooledMsg() *Message {
	return msgPool.Get().(*Message)
}

func FreeMsg(msg *Message) {
	msg.Reset()
	msgPool.Put(msg)
}

var poolUint32Data = sync.Pool{
	New: func() interface{} {
		data := make([]byte, 4)
		return &data
	},
}
