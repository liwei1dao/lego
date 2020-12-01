package proto

import "bufio"

func newSys(options Options) (sys *Proto, err error) {
	sys = &Proto{
		options: options,
	}
	return
}

type Proto struct {
	options Options
}

func (this *Proto) MessageDecodeBybufio(r *bufio.Reader) (message IMessage, err error) {
	return
}

func (this *Proto) MessageDecodeBybytes(buffer []byte) (message IMessage, err error) {
	return
}

func (this *Proto) MessageMarshal(comId uint16, msgId uint16, msg interface{}) (message IMessage) {
	return
}
