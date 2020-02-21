package proto

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"

	"github.com/golang/protobuf/proto"
)

const msgheadsize = 9
const msgmaxleng = 1024 * 1204 * 3
const DK_CHECK_NEW = 0x24 //校验类型

type Message struct {
	ComId   uint16
	MsgId   uint16
	CheckId byte
	MsgLen  uint32
	Buffer  []byte
}

func (this *Message) GetComId() uint16 {
	return this.ComId
}

func (this *Message) GetMsgId() uint16 {
	return this.MsgId
}

func (this *Message) GetMsg() []byte {
	return this.Buffer
}

func (this *Message) Serializable() (bytes []byte, err error) {
	if IsUseBigEndian {
		msg := []byte{}
		_msg := make([]byte, 2)
		binary.BigEndian.PutUint16(_msg, this.ComId)
		msg = append(msg, _msg...)
		_msg = make([]byte, 2)
		binary.BigEndian.PutUint16(_msg, this.MsgId)
		msg = append(msg, _msg...)
		msg = append(msg, this.CheckId)
		_msg = make([]byte, 4)
		binary.BigEndian.PutUint32(_msg, this.MsgLen)
		msg = append(msg, _msg...)
		msg = append(msg, this.Buffer...)
		return msg, nil
	} else {
		msg := []byte{}
		_msg := make([]byte, 2)
		binary.LittleEndian.PutUint16(_msg, this.ComId)
		msg = append(msg, _msg...)
		_msg = make([]byte, 2)
		binary.LittleEndian.PutUint16(_msg, this.MsgId)
		msg = append(msg, _msg...)
		msg = append(msg, this.CheckId)
		_msg = make([]byte, 4)
		binary.LittleEndian.PutUint32(_msg, this.MsgLen)
		msg = append(msg, _msg...)
		msg = append(msg, this.Buffer...)
		return msg, nil
	}
}

func (this *Message) ToString() string {
	if MsgProtoType == Proto_Json {
		return fmt.Sprintf("ComId:%d MsgId:%d Msg:%s", this.ComId, this.MsgId, string(this.Buffer))
	} else {
		return fmt.Sprintf("ComId:%d MsgId:%d Msg:%v", this.ComId, this.MsgId, this.Buffer)
	}
}

type DefMessageFactory struct {
}

func (this *DefMessageFactory) DefMessageDecodeBybufio(r *bufio.Reader) (msg IMessage, err error) {
	message := &Message{}
	message.ComId, err = ReadUInt16(r)
	if err != nil {
		return msg, err
	}
	message.MsgId, err = ReadUInt16(r)
	if err != nil {
		return msg, err
	}
	message.CheckId, err = ReadByte(r)
	if err != nil {
		return msg, err
	}
	if (message.CheckId & DK_CHECK_NEW) != DK_CHECK_NEW {
		return nil, fmt.Errorf("消息校验失败 CheckId: 0x%x", message.CheckId)
	}
	message.MsgLen, err = ReadUInt32(r)
	if err != nil {
		return msg, err
	}
	if message.MsgLen < msgheadsize {
		return nil, fmt.Errorf("消息解析失败 异常 MsgLen: %d", message.MsgLen)
	}
	if message.MsgLen > msgmaxleng {
		return nil, fmt.Errorf("消息解析失败 超长 MsgLen: %d", message.MsgLen)
	}
	message.Buffer = make([]byte, message.MsgLen-msgheadsize)
	_, err = io.ReadFull(r, message.Buffer)
	return message, err
}

func (this *DefMessageFactory) DefMessageDecodeBybytes(buffer []byte) (msg IMessage, err error) {
	if len(buffer) >= msgheadsize {
		message := &Message{}
		if IsUseBigEndian {
			message.ComId = binary.BigEndian.Uint16(buffer[0:])
			message.MsgId = binary.BigEndian.Uint16(buffer[2:])
			message.CheckId = buffer[4]
			message.MsgLen = binary.BigEndian.Uint32(buffer[5:])
		} else {
			message.ComId = binary.LittleEndian.Uint16(buffer[0:])
			message.MsgId = binary.LittleEndian.Uint16(buffer[2:])
			message.CheckId = buffer[4]
			message.MsgLen = binary.LittleEndian.Uint32(buffer[5:])
		}

		if int(message.MsgLen) >= msgheadsize && len(buffer) >= int(message.MsgLen) {
			message.Buffer = buffer[msgheadsize:message.MsgLen]
			return message, nil
		} else {
			return nil, fmt.Errorf("解析数据失败")
		}
	}
	return nil, fmt.Errorf("解析数据失败")
}

func (this *DefMessageFactory) DefMessageMarshal(comId uint16, msgId uint16, msg interface{}) IMessage {
	if MsgProtoType == Proto_Json {
		return jsonDefMessageMarshal(comId, msgId, msg)
	} else {
		return protoDefMessageMarshal(comId, msgId, msg)
	}
}

// Proto 消息 编码  方法
func protoDefMessageMarshal(comId uint16, msgId uint16, msg interface{}) IMessage {
	data, _ := proto.Marshal(msg.(proto.Message))
	message := &Message{
		ComId:   uint16(comId),
		MsgId:   uint16(msgId),
		CheckId: DK_CHECK_NEW,
		MsgLen:  uint32(len(data) + msgheadsize),
		Buffer:  data,
	}
	return message
}

// Json 消息 编码  方法 这个有点坑爹啊 _Msg里面的字段必须是大写开头的哈 不然无法序列化的
func jsonDefMessageMarshal(comId uint16, msgId uint16, msg interface{}) IMessage {
	data, _ := json.Marshal(msg)
	message := &Message{
		ComId:   uint16(comId),
		MsgId:   uint16(msgId),
		CheckId: DK_CHECK_NEW,
		MsgLen:  uint32(len(data) + msgheadsize),
		Buffer:  data,
	}
	return message
}
