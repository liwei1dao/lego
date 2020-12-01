package proto

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"io"
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

func (this *DefMessage) GetMsg() []byte {
	return this.Buffer
}

//默认消息工厂
type DefMessageFactory struct {
	msgProtoType   ProtoType //消息体类型
	isUseBigEndian bool      //消息传输码
	msgheadsize    uint32    //消息头大小
	msgmaxleng     uint32    //消息体最大长度
}

func (this *DefMessageFactory) MessageDecodeBybufio(r *bufio.Reader) (message IMessage, err error) {
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
	if msg.MsgLen < this.msgheadsize {
		return nil, fmt.Errorf("消息解析失败 异常 MsgLen: %d", message.MsgLen)
	}
	if msg.MsgLen > this.msgmaxleng {
		return nil, fmt.Errorf("消息解析失败 超长 MsgLen: %d", message.MsgLen)
	}
	msg.Buffer = make([]byte, msg.MsgLen-this.msgheadsize)
	_, err = io.ReadFull(r, msg.Buffer)
	return msg, err
}

func (this *DefMessageFactory) MessageDecodeBybytes(buffer []byte) (message IMessage, err error) {
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

		if uint32(msg.MsgLen) >= this.msgheadsize && len(buffer) >= int(msg.MsgLen) {
			msg.Buffer = buffer[this.msgheadsize:msg.MsgLen]
			return msg, nil
		} else {
			return nil, fmt.Errorf("解析数据失败")
		}
	}
	return nil, fmt.Errorf("DefMessageFactory MessageDecodeBybytes fail")
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
