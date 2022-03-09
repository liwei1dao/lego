package convert

import (
	"encoding/binary"
	"math"
)

func ByteToBytes(v byte) []byte {
	return []byte{v}
}

func BoolToBytes(v bool) []byte {
	var buf = make([]byte, 1)
	if v {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	return buf
}
func Int8ToBytes(v int8) []byte {
	return []byte{(byte(v))}
}
func Int16ToBytes(v int16) []byte {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(v))
	return buf
}
func UInt16ToBytes(v uint16) []byte {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, v)
	return buf
}
func IntToBytes(v int) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(v))
	return buf
}
func Int32ToBytes(v int32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(v))
	return buf
}
func UInt32ToBytes(v uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	return buf
}
func Int64ToBytes(v int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(v))
	return buf
}
func UInt64ToBytes(v uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v)
	return buf
}
func Float32ToBytes(v float32) []byte {
	bits := math.Float32bits(v)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}
func Float64ToBytes(v float64) []byte {
	bits := math.Float64bits(v)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
func StringToBytes(v string) []byte {
	return []byte(v)
}

func BytesTobyte(buf []byte) byte {
	var data byte = buf[0]
	return data
}
func BytesToBool(buf []byte) bool {
	var data bool = buf[0] != 0
	return data
}
func BytesToInt8(buf []byte) int8 {
	var data int8 = int8(buf[0])
	return data
}
func BytesToInt16(buf []byte) int16 {
	return int16(binary.BigEndian.Uint16(buf))
}
func BytesToUInt16(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf)
}
func BytesToInt(buf []byte) int {
	return int(binary.BigEndian.Uint32(buf))
}
func BytesToInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}
func BytesToUInt32(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}
func BytesToUInt64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}
func BytesToFloat32(buf []byte) float32 {
	bits := binary.LittleEndian.Uint32(buf)
	return math.Float32frombits(bits)
}
func BytesToFloat64(buf []byte) float64 {
	bits := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(bits)
}
func BytesToString(buf []byte) string {
	bits := string(buf)
	return bits
}
