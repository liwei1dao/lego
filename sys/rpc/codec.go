package rpc

import (
	"bytes"
	"fmt"

	proto "github.com/gogo/protobuf/proto"
	jsoniter "github.com/json-iterator/go"
	"github.com/liwei1dao/lego/sys/rpc/rpccore"
	"github.com/tinylib/msgp/msgp"
	"github.com/vmihailenco/msgpack/v5"
	pb "google.golang.org/protobuf/proto"
)

// Codecs are codecs supported by rpcx. You can add customized codecs in Codecs.
var codecs = map[rpccore.SerializeType]rpccore.ICodec{
	rpccore.JSON:        &JSONCodec{},
	rpccore.ProtoBuffer: &PBCodec{},
	rpccore.MsgPack:     &MsgpackCodec{},
}

type JSONCodec struct{}

func (c JSONCodec) Marshal(i interface{}) ([]byte, error) {
	return jsoniter.Marshal(i)
}

func (c JSONCodec) Unmarshal(data []byte, i interface{}) error {
	d := jsoniter.NewDecoder(bytes.NewBuffer(data))
	d.UseNumber()
	return d.Decode(i)
}

type PBCodec struct{}

func (c PBCodec) Marshal(i interface{}) ([]byte, error) {
	if m, ok := i.(proto.Marshaler); ok {
		return m.Marshal()
	}

	if m, ok := i.(pb.Message); ok {
		return pb.Marshal(m)
	}

	return nil, fmt.Errorf("%T is not a proto.Marshaler or pb.Message", i)
}

func (c PBCodec) Unmarshal(data []byte, i interface{}) error {
	if m, ok := i.(proto.Unmarshaler); ok {
		return m.Unmarshal(data)
	}

	if m, ok := i.(pb.Message); ok {
		return pb.Unmarshal(data, m)
	}

	return fmt.Errorf("%T is not a proto.Unmarshaler or pb.Message", i)
}

type MsgpackCodec struct{}

func (c MsgpackCodec) Marshal(i interface{}) ([]byte, error) {
	if m, ok := i.(msgp.Marshaler); ok {
		return m.MarshalMsg(nil)
	}
	var buf bytes.Buffer
	enc := msgpack.NewEncoder(&buf)
	// enc.UseJSONTag(true)
	err := enc.Encode(i)
	return buf.Bytes(), err
}

func (c MsgpackCodec) Unmarshal(data []byte, i interface{}) error {
	if m, ok := i.(msgp.Unmarshaler); ok {
		_, err := m.UnmarshalMsg(data)
		return err
	}
	dec := msgpack.NewDecoder(bytes.NewReader(data))
	// dec.UseJSONTag(true)
	err := dec.Decode(i)
	return err
}
