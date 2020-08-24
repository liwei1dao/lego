package serialize

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"

	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
	"github.com/liwei1dao/lego/sys/proto"
)

type SerializeData struct {
	dataType        reflect.Type
	SerializeFunc   func(d interface{}) ([]byte, error)
	UnSerializeFunc func(dataType reflect.Type, d []byte) (interface{}, error)
}

type RpcData struct {
	Dtype string
	Data  []byte
}

func GetRpcData(i interface{}) (*RpcData, string) {
	dtype, d, err := Serialize(i)
	if err != nil {
		return nil, err.Error()
	} else {
		return &RpcData{Dtype: dtype, Data: d}, ""
	}
}

const (
	NULL = "null" //nil   null
)

var (
	SerializeObjs map[string]*SerializeData
)

func SerializeInit() {
	SerializeObjs = make(map[string]*SerializeData)
	OnRegister(nil, NullToBytes, BytesToNull)
	OnRegister(byte(0), ByteToBytes, BytesTobyte)
	OnRegister(false, BoolToBytes, BytesToBool)
	OnRegister(int8(0), Int8ToBytes, BytesToInt8)
	OnRegister(int16(0), Int16ToBytes, BytesToInt16)
	OnRegister(uint16(0), UInt16ToBytes, BytesToUInt16)
	OnRegister(int(0), IntToBytes, BytesToInt)
	OnRegister(int32(0), Int32ToBytes, BytesToInt32)
	OnRegister(uint32(0), UInt32ToBytes, BytesToUInt32)
	OnRegister(int64(0), Int64ToBytes, BytesToInt64)
	OnRegister(uint64(0), UInt64ToBytes, BytesToUInt64)
	OnRegister(float32(0), Float32ToBytes, BytesToFloat32)
	OnRegister(float64(0), Float64ToBytes, BytesToFloat64)
	OnRegister(string(""), StringToBytes, BytesToString)
	OnRegister([]byte{}, SliceByteToBytes, BytesToSliceByte)
	OnRegister([]int32{}, SliceInt32ToBytes, BytesToSliceInt32)
	OnRegister([]uint32{}, SliceUInt32ToBytes, BytesToSliceUInt32)
	OnRegister([]string{}, SliceStringToBytes, BytesToSliceString)
	OnRegister([]interface{}{}, SliceInterfaceToBytes, BytesToSliceInterface)
	OnRegister(map[string]string{}, JsonStructMarshal, BytesToMapString)
	OnRegister(map[int32]string{}, JsonStructMarshal, Int32ToMapString)
	OnRegister(map[string]interface{}{}, MapStringInterfaceToBytes, BytesToMapStringRpc)
	OnRegister(map[int32]interface{}{}, MapInt32InterfaceToBytes, BytesToMapInt32Rpc)
	OnRegister(map[uint32]interface{}{}, MapUInt32InterfaceToBytes, BytesToMapUInt32Rpc)
	OnRegister(map[string]*RpcData{}, JsonStructMarshal, BytesToMapStringRpc)
	OnRegister(&proto.Message{}, ProtoStructMarshal, PrototructUnmarshal)
	OnRegister(core.ErrorCode(0), ErrorCodeToBytes, BytesToErrorCode)
	OnRegister(core.CustomRoute(0), CustomRouteToBytes, BytesToCustomRoute)
	OnRegister(map[uint16][]uint16{}, JsonStructMarshal, BytesToMapUnit16SliceUInt16Rpc)
	OnRegister(&core.ServiceMonitor{}, JsonStructMarshal, JsonStructUnmarshal)
}

func OnRegister(d interface{}, sf func(d interface{}) ([]byte, error), unsf func(dataType reflect.Type, d []byte) (interface{}, error)) {
	dtype := "Null"
	if d != nil {
		dtype = reflect.TypeOf(d).String()
	}
	if _, ok := SerializeObjs[dtype]; ok {
		log.Warnf("注册重复序列化数据结构 dtype = %s", dtype)
		return
	}
	SerializeObjs[dtype] = &SerializeData{
		dataType:        reflect.TypeOf(d),
		SerializeFunc:   sf,
		UnSerializeFunc: unsf,
	}
}

func Serialize(d interface{}) (dtype string, b []byte, err error) {
	dtype = "Null"
	if d != nil {
		dtype = reflect.TypeOf(d).String()
	}
	if v, ok := SerializeObjs[dtype]; ok {
		b, err = v.SerializeFunc(d)
		return
	}
	return dtype, nil, fmt.Errorf("没有注册序列化数据结构 dtype = %s", dtype)
}

func UnSerialize(dtype string, d []byte) (interface{}, error) {
	if v, ok := SerializeObjs[dtype]; ok {
		return v.UnSerializeFunc(v.dataType, d)
	}
	return nil, fmt.Errorf("没有注册序列化数据结构 dtype = %s", dtype)
}

//----------------------------------------------内置序列化--------------------------------------------------------------
//序列化----------------------------------------------------------------------------------------------------------------
func NullToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 0)
	return buf, nil
}
func ByteToBytes(v interface{}) ([]byte, error) {
	return []byte{v.(byte)}, nil
}
func BoolToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 1)
	if v.(bool) {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	return buf, nil
}
func Int8ToBytes(v interface{}) ([]byte, error) {
	return []byte{(byte(v.(int8)))}, nil
}
func Int16ToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(v.(int16)))
	return buf, nil
}
func UInt16ToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, v.(uint16))
	return buf, nil
}
func IntToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(v.(int)))
	return buf, nil
}
func Int32ToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(v.(int32)))
	return buf, nil
}
func UInt32ToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v.(uint32))
	return buf, nil
}
func Int64ToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(v.(int64)))
	return buf, nil
}
func UInt64ToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v.(uint64))
	return buf, nil
}
func Float32ToBytes(v interface{}) ([]byte, error) {
	bits := math.Float32bits(v.(float32))
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes, nil
}
func Float64ToBytes(v interface{}) ([]byte, error) {
	bits := math.Float64bits(v.(float64))
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes, nil
}
func StringToBytes(v interface{}) ([]byte, error) {
	return []byte(v.(string)), nil
}
func SliceByteToBytes(v interface{}) ([]byte, error) {
	return v.([]byte), nil
}
func SliceInt32ToBytes(v interface{}) ([]byte, error) {
	d := v.([]int32)
	data := []byte{}
	for _, v := range d {
		var buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v))
		data = append(data, buf[0:]...)
	}
	return data, nil
}
func SliceUInt32ToBytes(v interface{}) ([]byte, error) {
	d := v.([]uint32)
	data := []byte{}
	for _, v := range d {
		var buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, v)
		data = append(data, buf[0:]...)
	}
	return data, nil
}
func SliceStringToBytes(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func SliceInterfaceToBytes(v interface{}) ([]byte, error) {
	data := []*RpcData{}
	for _, v := range v.([]interface{}) {
		t, d, e := Serialize(v)
		if e != nil {
			return nil, e
		}
		data = append(data, &RpcData{Dtype: t, Data: d})
	}
	return json.Marshal(data)
}

func MapStringInterfaceToBytes(v interface{}) ([]byte, error) {
	data := make(map[string]*RpcData)
	for k, v := range v.(map[string]interface{}) {
		t, d, e := Serialize(v)
		if e != nil {
			return nil, e
		}
		data[k] = &RpcData{Dtype: t, Data: d}
	}
	return json.Marshal(data)
}

func MapInt32InterfaceToBytes(v interface{}) ([]byte, error) {
	data := make(map[int32]*RpcData)
	for k, v := range v.(map[int32]interface{}) {
		t, d, e := Serialize(v)
		if e != nil {
			return nil, e
		}
		data[k] = &RpcData{Dtype: t, Data: d}
	}
	return json.Marshal(data)
}
func MapUInt32InterfaceToBytes(v interface{}) ([]byte, error) {
	data := make(map[uint32]*RpcData)
	for k, v := range v.(map[uint32]interface{}) {
		t, d, e := Serialize(v)
		if e != nil {
			return nil, e
		}
		data[k] = &RpcData{Dtype: t, Data: d}
	}
	return json.Marshal(data)
}
func JsonStructMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func ProtoStructMarshal(v interface{}) ([]byte, error) {
	return v.(proto.IMessage).Serializable()
}
func ErrorCodeToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(v.(core.ErrorCode)))
	return buf, nil
}
func CustomRouteToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(v.(core.CustomRoute)))
	return buf, nil
}

//反序列化--------------------------------------------------------------------------------------------------------------
func BytesToNull(dataType reflect.Type, buf []byte) (interface{}, error) {
	return nil, nil
}
func BytesTobyte(dataType reflect.Type, buf []byte) (interface{}, error) {
	var data byte = buf[0]
	return data, nil
}
func BytesToBool(dataType reflect.Type, buf []byte) (interface{}, error) {
	var data bool = buf[0] != 0
	return data, nil
}
func BytesToInt8(dataType reflect.Type, buf []byte) (interface{}, error) {
	var data int8 = int8(buf[0])
	return data, nil
}
func BytesToInt16(dataType reflect.Type, buf []byte) (interface{}, error) {
	return int16(binary.BigEndian.Uint16(buf)), nil
}
func BytesToUInt16(dataType reflect.Type, buf []byte) (interface{}, error) {
	return binary.BigEndian.Uint16(buf), nil
}
func BytesToInt(dataType reflect.Type, buf []byte) (interface{}, error) {
	return int(binary.BigEndian.Uint32(buf)), nil
}
func BytesToInt32(dataType reflect.Type, buf []byte) (interface{}, error) {
	return int32(binary.BigEndian.Uint32(buf)), nil
}
func BytesToUInt32(dataType reflect.Type, buf []byte) (interface{}, error) {
	return binary.BigEndian.Uint32(buf), nil
}
func BytesToInt64(dataType reflect.Type, buf []byte) (interface{}, error) {
	return int64(binary.BigEndian.Uint64(buf)), nil
}
func BytesToUInt64(dataType reflect.Type, buf []byte) (interface{}, error) {
	return binary.BigEndian.Uint64(buf), nil
}
func BytesToFloat32(dataType reflect.Type, buf []byte) (interface{}, error) {
	bits := binary.LittleEndian.Uint32(buf)
	return math.Float32frombits(bits), nil
}
func BytesToFloat64(dataType reflect.Type, buf []byte) (interface{}, error) {
	bits := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(bits), nil
}
func BytesToString(dataType reflect.Type, buf []byte) (interface{}, error) {
	bits := string(buf)
	return bits, nil
}
func BytesToSliceInt32(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make([]int32, len(buf)/4)
	for i, _ := range data {
		data[i] = int32(binary.BigEndian.Uint32(buf[i*4:]))
	}
	return data, nil
}
func BytesToSliceUInt32(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make([]uint32, len(buf)/4)
	for i, _ := range data {
		data[i] = binary.BigEndian.Uint32(buf[i*4:])
	}
	return data, nil
}
func BytesToSliceString(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := []string{}
	err := json.Unmarshal(buf, &data)
	return data, err
}
func BytesToSliceInterface(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := []*RpcData{}
	err := json.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	mapdata := make([]interface{}, len(data))
	for i, v := range data {
		d, e := UnSerialize(v.Dtype, v.Data)
		if e != nil {
			return nil, e
		}
		mapdata[i] = d
	}
	return mapdata, err
}
func BytesToMapString(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make(map[string]string)
	err := json.Unmarshal(buf, &data)
	return data, err
}

func Int32ToMapString(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make(map[int32]string)
	err := json.Unmarshal(buf, &data)
	return data, err
}
func BytesToSliceByte(dataType reflect.Type, buf []byte) (interface{}, error) {
	return buf, nil
}
func BytesToMapStringRpc(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make(map[string]*RpcData)
	err := json.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	mapdata := make(map[string]interface{})
	for k, v := range data {
		d, e := UnSerialize(v.Dtype, v.Data)
		if e != nil {
			return nil, e
		}
		mapdata[k] = d
	}
	return mapdata, err
}

func BytesToMapInt32Rpc(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make(map[int32]*RpcData)
	err := json.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	mapdata := make(map[int32]interface{})
	for k, v := range data {
		d, e := UnSerialize(v.Dtype, v.Data)
		if e != nil {
			return nil, e
		}
		mapdata[k] = d
	}
	return mapdata, err
}

func BytesToMapUInt32Rpc(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make(map[uint32]*RpcData)
	err := json.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	mapdata := make(map[uint32]interface{})
	for k, v := range data {
		d, e := UnSerialize(v.Dtype, v.Data)
		if e != nil {
			return nil, e
		}
		mapdata[k] = d
	}
	return mapdata, err
}

func JsonStructUnmarshal(dataType reflect.Type, buf []byte) (interface{}, error) {
	msg := reflect.New(dataType.Elem()).Interface()
	err := json.Unmarshal(buf, msg)
	return msg, err
}

func PrototructUnmarshal(dataType reflect.Type, buf []byte) (interface{}, error) {
	return proto.MessageFactory.MessageDecodeBybytes(buf)
}

func BytesToErrorCode(dataType reflect.Type, buf []byte) (interface{}, error) {
	return core.ErrorCode(binary.BigEndian.Uint32(buf)), nil
}

func BytesToCustomRoute(dataType reflect.Type, buf []byte) (interface{}, error) {
	return core.CustomRoute(binary.BigEndian.Uint32(buf)), nil
}

func BytesToMapUnit16SliceUInt16Rpc(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make(map[uint16][]uint16)
	err := json.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	return data, err
}
