package core

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"math"
	"reflect"

	"github.com/golang/protobuf/proto"
	"github.com/liwei1dao/lego/core"
	"github.com/liwei1dao/lego/sys/log"
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
	serializeObjs map[string]*SerializeData
)

func init() {
	serializeObjs = make(map[string]*SerializeData)
	OnRegisterRpcData(nil, nullToBytes, bytesToNull)
	OnRegisterRpcData(byte(0), byteToBytes, bytesTobyte)
	OnRegisterRpcData(false, boolToBytes, bytesToBool)
	OnRegisterRpcData(int8(0), int8ToBytes, bytesToInt8)
	OnRegisterRpcData(int16(0), int16ToBytes, bytesToInt16)
	OnRegisterRpcData(uint16(0), uInt16ToBytes, bytesToUInt16)
	OnRegisterRpcData(int(0), intToBytes, bytesToInt)
	OnRegisterRpcData(int32(0), int32ToBytes, bytesToInt32)
	OnRegisterRpcData(uint32(0), uInt32ToBytes, bytesToUInt32)
	OnRegisterRpcData(int64(0), int64ToBytes, bytesToInt64)
	OnRegisterRpcData(uint64(0), uInt64ToBytes, bytesToUInt64)
	OnRegisterRpcData(float32(0), float32ToBytes, bytesToFloat32)
	OnRegisterRpcData(float64(0), float64ToBytes, bytesToFloat64)
	OnRegisterRpcData(string(""), stringToBytes, bytesToString)
	OnRegisterRpcData([]byte{}, sliceByteToBytes, bytesToSliceByte)
	OnRegisterRpcData([]int32{}, sliceInt32ToBytes, bytesToSliceInt32)
	OnRegisterRpcData([]uint32{}, sliceUInt32ToBytes, bytesToSliceUInt32)
	OnRegisterRpcData([]uint64{}, sliceUInt64ToBytes, bytesToSliceUInt64)
	OnRegisterRpcData([]string{}, sliceStringToBytes, bytesToSliceString)
	OnRegisterRpcData([]interface{}{}, sliceInterfaceToBytes, bytesToSliceInterface)
	OnRegisterRpcData(map[string]string{}, jsonStructMarshal, bytesToMapString)
	OnRegisterRpcData(map[int32]string{}, jsonStructMarshal, int32ToMapString)
	OnRegisterRpcData(map[string]interface{}{}, mapStringInterfaceToBytes, bytesToMapStringRpc)
	OnRegisterRpcData(map[int32]interface{}{}, mapInt32InterfaceToBytes, bytesToMapInt32Rpc)
	OnRegisterRpcData(map[uint32]interface{}{}, mapUInt32InterfaceToBytes, bytesToMapUInt32Rpc)
	OnRegisterRpcData(map[string]*RpcData{}, jsonStructMarshal, bytesToMapStringRpc)
	OnRegisterRpcData(core.ErrorCode(0), errorCodeToBytes, bytesToErrorCode)
	OnRegisterRpcData(core.CustomRoute(0), customRouteToBytes, bytesToCustomRoute)
	OnRegisterRpcData(map[uint16][]uint16{}, jsonStructMarshal, bytesToMapUnit16SliceUInt16Rpc)
	OnRegisterRpcData(&core.ServiceMonitor{}, jsonStructMarshal, jsonStructUnmarshal)
}

func Serialize(d interface{}) (dtype string, b []byte, err error) {
	dtype = "Null"
	if d != nil {
		dtype = reflect.TypeOf(d).String()
	}
	if v, ok := serializeObjs[dtype]; ok {
		b, err = v.SerializeFunc(d)
		return
	}
	return dtype, nil, fmt.Errorf("没有注册序列化数据结构 dtype = %s", dtype)
}

func UnSerialize(dtype string, d []byte) (interface{}, error) {
	if v, ok := serializeObjs[dtype]; ok {
		return v.UnSerializeFunc(v.dataType, d)
	}
	return nil, fmt.Errorf("没有注册序列化数据结构 dtype = %s", dtype)
}

func OnRegisterRpcData(d interface{}, sf func(d interface{}) ([]byte, error), unsf func(dataType reflect.Type, d []byte) (interface{}, error)) {
	dtype := "Null"
	if d != nil {
		dtype = reflect.TypeOf(d).String()
	}
	if _, ok := serializeObjs[dtype]; ok {
		log.Warnf("注册重复序列化数据结构 dtype = %s", dtype)
		return
	}
	serializeObjs[dtype] = &SerializeData{
		dataType:        reflect.TypeOf(d),
		SerializeFunc:   sf,
		UnSerializeFunc: unsf,
	}
}

//注册Proto Or Json 数据结构到RPC
func OnRegisterProtoData(d interface{}) (err error) {
	if _, ok := d.(proto.Message); ok {
		OnRegisterRpcData(d, protoStructMarshal, protoStructUnmarshal)
	} else {
		log.Errorf("rpc OnRegisterProtoData d:%v no is proto.Message", d)
	}
	return
}
func OnRegisterJsonRpc(d interface{}) (err error) {
	OnRegisterRpcData(d, jsonStructMarshal, jsonStructUnmarshal)
	return
}

//----------------------------------------------内置序列化--------------------------------------------------------------
//序列化----------------------------------------------------------------------------------------------------------------
func nullToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 0)
	return buf, nil
}
func byteToBytes(v interface{}) ([]byte, error) {
	return []byte{v.(byte)}, nil
}
func boolToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 1)
	if v.(bool) {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	return buf, nil
}
func int8ToBytes(v interface{}) ([]byte, error) {
	return []byte{(byte(v.(int8)))}, nil
}
func int16ToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(v.(int16)))
	return buf, nil
}
func uInt16ToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, v.(uint16))
	return buf, nil
}
func intToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(v.(int)))
	return buf, nil
}
func int32ToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(v.(int32)))
	return buf, nil
}
func uInt32ToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v.(uint32))
	return buf, nil
}
func int64ToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(v.(int64)))
	return buf, nil
}
func uInt64ToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v.(uint64))
	return buf, nil
}
func float32ToBytes(v interface{}) ([]byte, error) {
	bits := math.Float32bits(v.(float32))
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes, nil
}
func float64ToBytes(v interface{}) ([]byte, error) {
	bits := math.Float64bits(v.(float64))
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes, nil
}
func stringToBytes(v interface{}) ([]byte, error) {
	return []byte(v.(string)), nil
}
func sliceByteToBytes(v interface{}) ([]byte, error) {
	return v.([]byte), nil
}
func sliceInt32ToBytes(v interface{}) ([]byte, error) {
	d := v.([]int32)
	data := []byte{}
	for _, v := range d {
		var buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, uint32(v))
		data = append(data, buf[0:]...)
	}
	return data, nil
}
func sliceUInt32ToBytes(v interface{}) ([]byte, error) {
	d := v.([]uint32)
	data := []byte{}
	for _, v := range d {
		var buf = make([]byte, 4)
		binary.BigEndian.PutUint32(buf, v)
		data = append(data, buf[0:]...)
	}
	return data, nil
}
func sliceUInt64ToBytes(v interface{}) ([]byte, error) {
	d := v.([]uint64)
	data := []byte{}
	for _, v := range d {
		var buf = make([]byte, 8)
		binary.BigEndian.PutUint64(buf, v)
		data = append(data, buf[0:]...)
	}
	return data, nil
}
func sliceStringToBytes(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}
func sliceInterfaceToBytes(v interface{}) ([]byte, error) {
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
func mapStringInterfaceToBytes(v interface{}) ([]byte, error) {
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
func mapInt32InterfaceToBytes(v interface{}) ([]byte, error) {
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
func mapUInt32InterfaceToBytes(v interface{}) ([]byte, error) {
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
func jsonStructMarshal(v interface{}) ([]byte, error) {
	return json.Marshal(v)
}

func protoStructMarshal(v interface{}) ([]byte, error) {
	return proto.Marshal(v.(proto.Message))
}

func errorCodeToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(v.(core.ErrorCode)))
	return buf, nil
}
func customRouteToBytes(v interface{}) ([]byte, error) {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(v.(core.CustomRoute)))
	return buf, nil
}

//反序列化--------------------------------------------------------------------------------------------------------------
func bytesToNull(dataType reflect.Type, buf []byte) (interface{}, error) {
	return nil, nil
}
func bytesTobyte(dataType reflect.Type, buf []byte) (interface{}, error) {
	var data byte = buf[0]
	return data, nil
}
func bytesToBool(dataType reflect.Type, buf []byte) (interface{}, error) {
	var data bool = buf[0] != 0
	return data, nil
}
func bytesToInt8(dataType reflect.Type, buf []byte) (interface{}, error) {
	var data int8 = int8(buf[0])
	return data, nil
}
func bytesToInt16(dataType reflect.Type, buf []byte) (interface{}, error) {
	return int16(binary.BigEndian.Uint16(buf)), nil
}
func bytesToUInt16(dataType reflect.Type, buf []byte) (interface{}, error) {
	return binary.BigEndian.Uint16(buf), nil
}
func bytesToInt(dataType reflect.Type, buf []byte) (interface{}, error) {
	return int(binary.BigEndian.Uint32(buf)), nil
}
func bytesToInt32(dataType reflect.Type, buf []byte) (interface{}, error) {
	return int32(binary.BigEndian.Uint32(buf)), nil
}
func bytesToUInt32(dataType reflect.Type, buf []byte) (interface{}, error) {
	return binary.BigEndian.Uint32(buf), nil
}
func bytesToInt64(dataType reflect.Type, buf []byte) (interface{}, error) {
	return int64(binary.BigEndian.Uint64(buf)), nil
}
func bytesToUInt64(dataType reflect.Type, buf []byte) (interface{}, error) {
	return binary.BigEndian.Uint64(buf), nil
}
func bytesToFloat32(dataType reflect.Type, buf []byte) (interface{}, error) {
	bits := binary.LittleEndian.Uint32(buf)
	return math.Float32frombits(bits), nil
}
func bytesToFloat64(dataType reflect.Type, buf []byte) (interface{}, error) {
	bits := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(bits), nil
}
func bytesToString(dataType reflect.Type, buf []byte) (interface{}, error) {
	bits := string(buf)
	return bits, nil
}
func bytesToSliceInt32(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make([]int32, len(buf)/4)
	for i, _ := range data {
		data[i] = int32(binary.BigEndian.Uint32(buf[i*4:]))
	}
	return data, nil
}
func bytesToSliceUInt32(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make([]uint32, len(buf)/4)
	for i, _ := range data {
		data[i] = binary.BigEndian.Uint32(buf[i*4:])
	}
	return data, nil
}
func bytesToSliceUInt64(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make([]uint64, len(buf)/8)
	for i, _ := range data {
		data[i] = binary.BigEndian.Uint64(buf[i*8:])
	}
	return data, nil
}
func bytesToSliceString(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := []string{}
	err := json.Unmarshal(buf, &data)
	return data, err
}
func bytesToSliceInterface(dataType reflect.Type, buf []byte) (interface{}, error) {
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
func bytesToMapString(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make(map[string]string)
	err := json.Unmarshal(buf, &data)
	return data, err
}
func int32ToMapString(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make(map[int32]string)
	err := json.Unmarshal(buf, &data)
	return data, err
}
func bytesToSliceByte(dataType reflect.Type, buf []byte) (interface{}, error) {
	return buf, nil
}
func bytesToMapStringRpc(dataType reflect.Type, buf []byte) (interface{}, error) {
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
func bytesToMapInt32Rpc(dataType reflect.Type, buf []byte) (interface{}, error) {
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
func bytesToMapUInt32Rpc(dataType reflect.Type, buf []byte) (interface{}, error) {
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
func jsonStructUnmarshal(dataType reflect.Type, buf []byte) (interface{}, error) {
	msg := reflect.New(dataType.Elem()).Interface()
	err := json.Unmarshal(buf, msg)
	return msg, err
}

func protoStructUnmarshal(dataType reflect.Type, buf []byte) (interface{}, error) {
	msg := reflect.New(dataType.Elem()).Interface()
	err := proto.UnmarshalMerge(buf, msg.(proto.Message))
	return msg, err
}

func bytesToErrorCode(dataType reflect.Type, buf []byte) (interface{}, error) {
	return core.ErrorCode(binary.BigEndian.Uint32(buf)), nil
}
func bytesToCustomRoute(dataType reflect.Type, buf []byte) (interface{}, error) {
	return core.CustomRoute(binary.BigEndian.Uint32(buf)), nil
}
func bytesToMapUnit16SliceUInt16Rpc(dataType reflect.Type, buf []byte) (interface{}, error) {
	data := make(map[uint16][]uint16)
	err := json.Unmarshal(buf, &data)
	if err != nil {
		return nil, err
	}
	return data, err
}
