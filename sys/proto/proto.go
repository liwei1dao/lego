package proto

import (
	"bufio"
	"encoding/json"
	"reflect"

	"github.com/golang/protobuf/proto"
)

func newSys(options Options) (sys *Proto, err error) {
	sys = &Proto{
		options: options,
	}
	options.MessageFactory.SetMessageConfig(options.MsgProtoType, options.IsUseBigEndian)
	return
}

type Proto struct {
	options Options
}

func (this *Proto) DecodeMessageBybufio(r *bufio.Reader) (message IMessage, err error) {
	return this.options.MessageFactory.DecodeMessageBybufio(r)
}

func (this *Proto) DecodeMessageBybytes(buffer []byte) (message IMessage, err error) {
	return this.options.MessageFactory.DecodeMessageBybytes(buffer)
}

func (this *Proto) EncodeToMesage(comId uint16, msgId uint16, msg interface{}) (message IMessage) {
	return this.options.MessageFactory.EncodeToMesage(comId, msgId, msg)
}

func (this *Proto) EncodeToByte(message IMessage) (buffer []byte) {
	return this.options.MessageFactory.EncodeToByte(message)
}

func (this *Proto) ByteDecodeToStruct(t reflect.Type, d []byte) (data interface{}, err error) {
	if this.options.MsgProtoType == Proto_Json {
		data = reflect.New(t.Elem()).Interface()
		err = json.Unmarshal(d, data)
	} else {
		data = reflect.New(t.Elem()).Interface()
		err = proto.UnmarshalMerge(d, data.(proto.Message))
	}
	return
}
