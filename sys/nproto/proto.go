package proto

import (
	"bufio"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"reflect"

	"github.com/golang/protobuf/proto"
)

func JsonStructUnmarshal(d []byte, v interface{}) error {
	return json.Unmarshal(d, v)
}
func ProtoStructUnmarshal(d []byte, v interface{}) error {
	if pb, ok := v.(proto.Message); !ok {
		return fmt.Errorf("ProtoStructUnmarshal 对象不是 proto.Message")
	} else {
		return proto.Unmarshal(d, pb)
	}
}
func JsonStructMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func ProtoStructMarshal(v interface{}) ([]byte, error) {
	if pb, ok := v.(proto.Message); !ok {
		return nil, fmt.Errorf("ProtoStructUnmarshal 对象不是 proto.Message")
	} else {
		return proto.Marshal(pb)
	}
}
func msgToString(v interface{}) string {
	if msg, ok := v.(IMsgMarshalString); ok {
		msgstr, err := msg.ToString()
		if err == nil {
			return msgstr
		} else {
			return fmt.Sprintf("解析异常:%s", err.Error())
		}
	} else {
		msgstr, err := json.Marshal(v)
		if err == nil {
			return string(msgstr)
		} else {
			return fmt.Sprintf("解析异常:%s", err.Error())
		}
	}
}

// Json 解码  方法
func jsonStructUnmarshal(msgType reflect.Type, d []byte) (interface{}, error) {
	msg := reflect.New(msgType.Elem()).Interface()
	err := json.Unmarshal(d, msg)
	// if err := json.Unmarshal(d, msg); err != nil {
	// 	log.Errorf("json [%s]序列化 结构错误 err:%s", string(d), err.Error())
	// }
	return msg, err
}

// Proto 解码  方法
func protoStructUnmarshal(msgType reflect.Type, d []byte) (interface{}, error) {
	msg := reflect.New(msgType.Elem()).Interface()
	err := proto.UnmarshalMerge(d, msg.(proto.Message))
	// if err := proto.UnmarshalMerge(d, msg.(proto.Message)); err != nil {
	// 	log.Errorf("proto 序列化 结构错误 err:%s", err.Error())
	// }
	return msg, err
}

func ReadByte(r *bufio.Reader) (byte, error) {
	buf := make([]byte, 1)
	_, err := io.ReadFull(r, buf[:1])
	if err != nil {
		return 0, err
	}
	return buf[0], nil
}

func ReadUInt16(r *bufio.Reader) (uint16, error) {
	buf := make([]byte, 2)
	_, err := io.ReadFull(r, buf[:2])
	if err != nil {
		return 0, err
	}
	if option.IsUseBigEndian {
		return binary.BigEndian.Uint16(buf[:2]), nil
	} else {
		return binary.LittleEndian.Uint16(buf[:2]), nil
	}
}

func ReadUInt32(r *bufio.Reader) (uint32, error) {
	buf := make([]byte, 4)
	_, err := io.ReadFull(r, buf[:4])
	if err != nil {
		return 0, err
	}
	if option.IsUseBigEndian {
		return binary.BigEndian.Uint32(buf[:4]), nil
	} else {
		return binary.LittleEndian.Uint32(buf[:4]), nil
	}
}
