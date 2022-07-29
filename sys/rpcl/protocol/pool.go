package protocol

import "sync"

var msgPool = sync.Pool{
	New: func() interface{} {
		header := Header([12]byte{})
		header[0] = magicNumber

		return &Message{
			Header: &header,
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
