package codec

import (
	"encoding/binary"
	"math"
	"unsafe"
)

const host32bit = ^uint(0)>>32 == 0

// string
func BytesToString(b []byte) string {
	return *(*string)(unsafe.Pointer(&b))
}
func StringToBytes(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(
		&struct {
			string
			Cap int
		}{s, len(s)},
	))
}

// int
func IntToBytes(v int) []byte {
	if host32bit {
		return Int32ToBytes(int32(v))
	} else {
		return Int64ToBytes(int64(v))
	}
}
func BytesToInt(buf []byte) int {
	if host32bit {
		return int(BytesToInt32(buf))
	} else {
		return int(BytesToInt64(buf))
	}
}

//int8
func Int8ToBytes(v int8) []byte {
	return []byte{(byte(v))}
}
func BytesToInt8(buf []byte) int8 {
	return int8(buf[0])
}

//int16
func Int16ToBytes(v int16) []byte {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, uint16(v))
	return buf
}
func BytesToInt16(buf []byte) int16 {
	return int16(binary.BigEndian.Uint16(buf))
}

//int32
func Int32ToBytes(v int32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, uint32(v))
	return buf
}
func BytesToInt32(buf []byte) int32 {
	return int32(binary.BigEndian.Uint32(buf))
}

//int64
func Int64ToBytes(v int64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, uint64(v))
	return buf
}
func BytesToInt64(buf []byte) int64 {
	return int64(binary.BigEndian.Uint64(buf))
}

//uint
func UIntToBytes(v uint) []byte {
	if host32bit {
		return Int32ToBytes(int32(v))
	} else {
		return Int64ToBytes(int64(v))
	}
}
func BytesToUInt(buf []byte) uint {
	if host32bit {
		return uint(BytesToUInt32(buf))
	} else {
		return uint(BytesToUInt64(buf))
	}
}

//uint8
func UInt8ToBytes(v uint8) []byte {
	return []byte{v}
}
func BytesToUInt8(buf []byte) uint8 {
	var data uint8 = uint8(buf[0])
	return data
}

//uint16
func UInt16ToBytes(v uint16) []byte {
	var buf = make([]byte, 2)
	binary.BigEndian.PutUint16(buf, v)
	return buf
}
func BytesToUInt16(buf []byte) uint16 {
	return binary.BigEndian.Uint16(buf)
}

//uint32
func UInt32ToBytes(v uint32) []byte {
	var buf = make([]byte, 4)
	binary.BigEndian.PutUint32(buf, v)
	return buf
}
func BytesToUInt32(buf []byte) uint32 {
	return binary.BigEndian.Uint32(buf)
}

//uint64
func BytesToUInt64(buf []byte) uint64 {
	return binary.BigEndian.Uint64(buf)
}
func UInt64ToBytes(v uint64) []byte {
	var buf = make([]byte, 8)
	binary.BigEndian.PutUint64(buf, v)
	return buf
}

//float32
func Float32ToBytes(v float32) []byte {
	bits := math.Float32bits(v)
	bytes := make([]byte, 4)
	binary.LittleEndian.PutUint32(bytes, bits)
	return bytes
}
func BytesToFloat32(buf []byte) float32 {
	bits := binary.LittleEndian.Uint32(buf)
	return math.Float32frombits(bits)
}

//float64
func Float64ToBytes(v float64) []byte {
	bits := math.Float64bits(v)
	bytes := make([]byte, 8)
	binary.LittleEndian.PutUint64(bytes, bits)
	return bytes
}
func BytesToFloat64(buf []byte) float64 {
	bits := binary.LittleEndian.Uint64(buf)
	return math.Float64frombits(bits)
}

//bool
func BoolToBytes(v bool) []byte {
	var buf = make([]byte, 1)
	if v {
		buf[0] = 1
	} else {
		buf[0] = 0
	}
	return buf
}
func BytesToBool(buf []byte) bool {
	var data bool = buf[0] != 0
	return data
}
