package proto

import "bufio"

type (
	IMessage interface {
		GetComId() uint16
		GetMsgId() uint16
		GetMsg() []byte
	}
	IProto interface {
		MessageDecodeBybufio(r *bufio.Reader) (message IMessage, err error)
		MessageDecodeBybytes(buffer []byte) (message IMessage, err error)
		MessageMarshal(comId uint16, msgId uint16, msg interface{}) (message IMessage)
	}
)

var (
	defsys IProto
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys IProto, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func MessageDecodeBybufio(r *bufio.Reader) (IMessage, error) {
	return defsys.MessageDecodeBybufio(r)
}
func MessageDecodeBybytes(buffer []byte) (msg IMessage, err error) {
	return defsys.MessageDecodeBybytes(buffer)
}
