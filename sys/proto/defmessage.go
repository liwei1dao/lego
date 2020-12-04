package proto

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/liwei1dao/lego/sys/rpc"
)

//默认消息体
type DefMessage struct {
	ComId  uint16 //主Id
	MsgId  uint16 //次Id
	MsgLen uint32 //消息体长度
	Buffer []byte //消息体
}

func (this *DefMessage) GetComId() uint16 {
	return this.ComId
}

func (this *DefMessage) GetMsgId() uint16 {
	return this.MsgId
}

func (this *DefMessage) GetMsgLen() uint32 {
	return this.MsgLen
}

func (this *DefMessage) GetBuffer() []byte {
	return this.Buffer
}

func (this *DefMessage) ToStriing() string {
	return fmt.Sprintf("ComId：%d MsgId:%d MsgLen:%d", this.ComId, this.MsgId, this.MsgLen)
}

//默认消息工厂
type DefMessageFactory struct {
	msgProtoType   ProtoType //消息体类型
	isUseBigEndian bool      //消息传输码
	msgheadsize    uint32    //消息头大小
	msgmaxleng     uint32    //消息体最大长度
}

func (this *DefMessageFactory) SetMessageConfig(MsgProtoType ProtoType, IsUseBigEndian bool) {
	this.msgProtoType = MsgProtoType
	this.isUseBigEndian = IsUseBigEndian
	this.msgheadsize = 8
	this.msgmaxleng = 1024 * 1204 * 3
	//自己rpc消息解析
	rpc.OnRegisterRpcData(&DefMessage{}, this.RpcEncodeMessage, this.RpcDecodeMessage)
}

func (this *DefMessageFactory) DecodeMessageBybufio(r *bufio.Reader) (message IMessage, err error) {
	msg := &DefMessage{}
	msg.ComId, err = this.readUInt16(r)
	if err != nil {
		return msg, err
	}
	msg.MsgId, err = this.readUInt16(r)
	if err != nil {
		return msg, err
	}
	msg.MsgLen, err = this.readUInt32(r)
	if err != nil {
		return msg, err
	}
	if msg.MsgLen > this.msgmaxleng {
		return nil, fmt.Errorf("DecodeMessageBybufio err msg.MsgLen:%d Super long", msg.MsgLen)
	}
	msg.Buffer = make([]byte, msg.MsgLen)
	_, err = io.ReadFull(r, msg.Buffer)
	return msg, err
}

func (this *DefMessageFactory) DecodeMessageBybytes(buffer []byte) (message IMessage, err error) {
	if uint32(len(buffer)) >= this.msgheadsize {
		msg := &DefMessage{}
		if this.isUseBigEndian {
			msg.ComId = binary.BigEndian.Uint16(buffer[0:])
			msg.MsgId = binary.BigEndian.Uint16(buffer[2:])
			msg.MsgLen = binary.BigEndian.Uint32(buffer[4:])
		} else {
			msg.ComId = binary.LittleEndian.Uint16(buffer[0:])
			msg.MsgId = binary.LittleEndian.Uint16(buffer[2:])
			msg.MsgLen = binary.LittleEndian.Uint32(buffer[4:])
		}

		if uint32(len(buffer)) >= msg.MsgLen+this.msgheadsize {
			msg.Buffer = buffer[this.msgheadsize : msg.MsgLen+this.msgheadsize]
			return msg, nil
		} else {
			return nil, fmt.Errorf("DecodeMessageBybytes err package:%v msg.MsgLen:%d", buffer, msg.MsgLen)
		}
	}
	return nil, fmt.Errorf("DecodeMessageBybytes err package:%v", buffer)
}

func (this *DefMessageFactory) EncodeToMesage(comId uint16, msgId uint16, msg interface{}) (message IMessage) {
	defmessage := &DefMessage{
		ComId: uint16(comId),
		MsgId: uint16(msgId),
	}
	if this.msgProtoType == Proto_Buff {
		defmessage.Buffer, _ = proto.Marshal(msg.(proto.Message))
	} else {
		defmessage.Buffer, _ = json.Marshal(msg)
	}
	defmessage.MsgLen = uint32(len(defmessage.Buffer))
	return defmessage
}

func (this *DefMessageFactory) EncodeToByte(message IMessage) (buffer []byte) {
	if this.isUseBigEndian {
		buffer := []byte{}
		_msg := make([]byte, 2)
		binary.BigEndian.PutUint16(_msg, message.GetComId())
		buffer = append(buffer, _msg...)
		_msg = make([]byte, 2)
		binary.BigEndian.PutUint16(_msg, message.GetMsgId())
		buffer = append(buffer, _msg...)
		_msg = make([]byte, 4)
		binary.BigEndian.PutUint32(_msg, message.GetMsgLen())
		buffer = append(buffer, _msg...)
		buffer = append(buffer, message.GetBuffer()...)
		return buffer
	} else {
		buffer := []byte{}
		_msg := make([]byte, 2)
		binary.LittleEndian.PutUint16(_msg, message.GetComId())
		buffer = append(buffer, _msg...)
		_msg = make([]byte, 2)
		binary.LittleEndian.PutUint16(_msg, message.GetMsgId())
		buffer = append(buffer, _msg...)
		_msg = make([]byte, 4)
		binary.LittleEndian.PutUint32(_msg, message.GetMsgLen())
		buffer = append(buffer, _msg...)
		buffer = append(buffer, message.GetBuffer()...)
		return buffer
	}
}

func (this *DefMessageFactory) RpcEncodeMessage(d interface{}) ([]byte, error) {
	message := d.(IMessage)
	if this.isUseBigEndian {
		buffer := []byte{}
		_msg := make([]byte, 2)
		binary.BigEndian.PutUint16(_msg, message.GetComId())
		buffer = append(buffer, _msg...)
		_msg = make([]byte, 2)
		binary.BigEndian.PutUint16(_msg, message.GetMsgId())
		buffer = append(buffer, _msg...)
		_msg = make([]byte, 4)
		binary.BigEndian.PutUint32(_msg, message.GetMsgLen())
		buffer = append(buffer, _msg...)
		buffer = append(buffer, message.GetBuffer()...)
		return buffer, nil
	} else {
		buffer := []byte{}
		_msg := make([]byte, 2)
		binary.LittleEndian.PutUint16(_msg, message.GetComId())
		buffer = append(buffer, _msg...)
		_msg = make([]byte, 2)
		binary.LittleEndian.PutUint16(_msg, message.GetMsgId())
		buffer = append(buffer, _msg...)
		_msg = make([]byte, 4)
		binary.LittleEndian.PutUint32(_msg, message.GetMsgLen())
		buffer = append(buffer, _msg...)
		buffer = append(buffer, message.GetBuffer()...)
		return buffer, nil
	}
}

func (this *DefMessageFactory) RpcDecodeMessage(dataType reflect.Type, d []byte) (interface{}, error) {
	return this.DecodeMessageBybytes(d)
}

func (this *DefMessageFactory) readByte(r *bufio.Reader) (byte, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(r, buf[:1])
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func (this *DefMessageFactory) readUInt16(r *bufio.Reader) (uint16, error) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(r, buf[:2])
	if err != nil {
		return 0, err
	}
	if this.isUseBigEndian {
		return binary.BigEndian.Uint16(buf[:2]), nil
	} else {
		return binary.LittleEndian.Uint16(buf[:2]), nil
	}
}

func (this *DefMessageFactory) readUInt32(r *bufio.Reader) (uint32, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(r, buf[:4])
	if err != nil {
		return 0, err
	}
	if this.isUseBigEndian {
		return binary.BigEndian.Uint32(buf[:4]), nil
	} else {
		return binary.LittleEndian.Uint32(buf[:4]), nil
	}
}
