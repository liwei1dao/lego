package utils

import (
	"errors"
	"math"
	"strconv"
)

var digits []uint32
var intDigits []int8

const uint32SafeToMultiply10 = uint32(0xffffffff)/10 - 1
const uint64SafeToMultiple10 = uint64(0xffffffffffffffff)/10 - 1
const maxFloat64 = 1<<53 - 1

func init() {
	digits = make([]uint32, 1000)
	for i := uint32(0); i < 1000; i++ {
		digits[i] = (((i / 100) + '0') << 16) + ((((i / 10) % 10) + '0') << 8) + i%10 + '0'
		if i < 10 {
			digits[i] += 2 << 24
		} else if i < 100 {
			digits[i] += 1 << 24
		}
	}

	intDigits = make([]int8, 256)
	for i := 0; i < len(intDigits); i++ {
		intDigits[i] = invalidCharForNumber
	}
	for i := int8('0'); i <= int8('9'); i++ {
		intDigits[i] = i - int8('0')
	}
}

//Read Int For String--------------------------------------------------------------------------------------------------------------------------------
func ReadInt8ForString(buf []byte) (ret int8, n int, err error) {
	var (
		val uint32
	)
	c := buf[0]
	if c == '-' {
		val, n, err = ReadUint32ForString(buf[1:])
		n++
		if val > math.MaxInt8+1 {
			err = errors.New("ReadInt8ForString overflow: " + strconv.FormatInt(int64(val), 10))
			return
		}
		ret = -int8(val)
		return
	}
	val, n, err = ReadUint32ForString(buf)
	if val > math.MaxInt8 {
		err = errors.New("ReadInt8ForString overflow: " + strconv.FormatInt(int64(val), 10))
		return
	}
	ret = int8(val)
	return
}

func ReadInt16ForString(buf []byte) (ret int16, n int, err error) {
	var (
		val uint32
	)
	c := buf[0]
	if c == '-' {
		val, n, err = ReadUint32ForString(buf[1:])
		n++
		if val > math.MaxInt16+1 {
			err = errors.New("ReadInt16ForString overflow: " + strconv.FormatInt(int64(val), 10))
			return
		}
		ret = -int16(val)
		return
	}
	val, n, err = ReadUint32ForString(buf)
	if val > math.MaxInt16 {
		err = errors.New("ReadInt16ForString overflow: " + strconv.FormatInt(int64(val), 10))
		return
	}
	ret = int16(val)
	return
}

func ReadInt32ForString(buf []byte) (ret int32, n int, err error) {
	var (
		val uint32
	)
	c := buf[0]
	if c == '-' {
		val, n, err = ReadUint32ForString(buf[1:])
		n++
		if val > math.MaxInt32+1 {
			err = errors.New("ReadInt32ForString overflow: " + strconv.FormatInt(int64(val), 10))
			return
		}
		ret = -int32(val)
		return
	}
	val, n, err = ReadUint32ForString(buf)
	if val > math.MaxInt32 {
		err = errors.New("ReadInt32ForString overflow: " + strconv.FormatInt(int64(val), 10))
		return
	}
	ret = int32(val)
	return
}

func ReadInt64ForString(buf []byte) (ret int64, n int, err error) {
	var (
		val uint64
	)
	c := buf[0]
	if c == '-' {
		val, n, err = ReadUint64ForString(buf[1:])
		n++
		if val > math.MaxInt64+1 {
			err = errors.New("ReadInt64ForString overflow: " + strconv.FormatInt(int64(val), 10))
			return
		}
		ret = -int64(val)
		return
	}
	val, n, err = ReadUint64ForString(buf)
	if val > math.MaxInt64 {
		err = errors.New("ReadInt64ForString overflow: " + strconv.FormatInt(int64(val), 10))
		return
	}
	ret = int64(val)
	return
}

func ReadUint8ForString(buf []byte) (ret uint8, n int, err error) {
	var (
		val uint32
	)
	val, n, err = ReadUint32ForString(buf)
	if val > math.MaxUint8 {
		err = errors.New("ReadUInt8ForString overflow: " + strconv.FormatInt(int64(val), 10))
		return
	}
	ret = uint8(val)
	return
}

func ReadUint16ForString(buf []byte) (ret uint16, n int, err error) {
	var (
		val uint32
	)
	val, n, err = ReadUint32ForString(buf)
	if val > math.MaxUint16 {
		err = errors.New("ReadUInt16ForString overflow: " + strconv.FormatInt(int64(val), 10))
		return
	}
	ret = uint16(val)
	return
}

func ReadUint32ForString(buf []byte) (ret uint32, n int, err error) {
	ind := intDigits[buf[0]]
	n = 1
	if ind == 0 {
		err = assertInteger(buf[1:])
		return
	}
	if ind == invalidCharForNumber {
		err = errors.New("ReadUint32ForString unexpected character: " + string([]byte{byte(ind)}))
		return
	}
	ret = uint32(ind)
	if len(buf) > 11 {
		i := 1
		ind2 := intDigits[buf[i]]
		if ind2 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			return
		}
		i++
		ind3 := intDigits[buf[i]]
		if ind3 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			ret = ret*10 + uint32(ind2)
			return
		}
		i++
		ind4 := intDigits[buf[i]]
		if ind4 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			ret = ret*100 + uint32(ind2)*10 + uint32(ind3)
			return
		}
		i++
		ind5 := intDigits[buf[i]]
		if ind5 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			ret = ret*1000 + uint32(ind2)*100 + uint32(ind3)*10 + uint32(ind4)
			return
		}
		i++
		ind6 := intDigits[buf[i]]
		if ind6 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			ret = ret*10000 + uint32(ind2)*1000 + uint32(ind3)*100 + uint32(ind4)*10 + uint32(ind5)
			return
		}
		i++
		ind7 := intDigits[buf[i]]
		if ind7 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			ret = ret*100000 + uint32(ind2)*10000 + uint32(ind3)*1000 + uint32(ind4)*100 + uint32(ind5)*10 + uint32(ind6)
			return
		}
		i++
		ind8 := intDigits[buf[i]]
		if ind8 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			ret = ret*1000000 + uint32(ind2)*100000 + uint32(ind3)*10000 + uint32(ind4)*1000 + uint32(ind5)*100 + uint32(ind6)*10 + uint32(ind7)
			return
		}
		i++
		ind9 := intDigits[buf[i]]
		ret = ret*10000000 + uint32(ind2)*1000000 + uint32(ind3)*100000 + uint32(ind4)*10000 + uint32(ind5)*1000 + uint32(ind6)*100 + uint32(ind7)*10 + uint32(ind8)
		n = i
		if ind9 == invalidCharForNumber {
			err = assertInteger(buf[i:])
			return
		}
	}
	for i := n; i < len(buf); i++ {
		ind = intDigits[buf[i]]
		if ind == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			return
		}
		if ret > uint32SafeToMultiply10 {
			value2 := (ret << 3) + (ret << 1) + uint32(ind)
			if value2 < ret {
				err = errors.New("ReadUint32ForString overflow")
				return
			}
			n = i
			ret = value2
			continue
		}
		ret = (ret << 3) + (ret << 1) + uint32(ind)
	}
	n = len(buf)
	return
}

func ReadUint64ForString(buf []byte) (ret uint64, n int, err error) {
	ind := intDigits[buf[0]]
	n = 1
	if ind == 0 {
		err = assertInteger(buf[1:])
		return
	}
	if ind == invalidCharForNumber {
		err = errors.New("ReadUint64ForString unexpected character: " + string([]byte{byte(ind)}))
		return
	}
	ret = uint64(ind)
	if len(buf) > 11 {
		i := 1
		ind2 := intDigits[buf[i]]
		if ind2 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			return
		}
		i++
		ind3 := intDigits[buf[i]]
		if ind3 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			ret = ret*10 + uint64(ind2)
			return
		}
		i++
		ind4 := intDigits[buf[i]]
		if ind4 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			ret = ret*100 + uint64(ind2)*10 + uint64(ind3)
			return
		}
		i++
		ind5 := intDigits[buf[i]]
		if ind5 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			ret = ret*1000 + uint64(ind2)*100 + uint64(ind3)*10 + uint64(ind4)
			return
		}
		i++
		ind6 := intDigits[buf[i]]
		if ind6 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			ret = ret*10000 + uint64(ind2)*1000 + uint64(ind3)*100 + uint64(ind4)*10 + uint64(ind5)
			return
		}
		i++
		ind7 := intDigits[buf[i]]
		if ind7 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			ret = ret*100000 + uint64(ind2)*10000 + uint64(ind3)*1000 + uint64(ind4)*100 + uint64(ind5)*10 + uint64(ind6)
			return
		}
		i++
		ind8 := intDigits[buf[i]]
		if ind8 == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			ret = ret*1000000 + uint64(ind2)*100000 + uint64(ind3)*10000 + uint64(ind4)*1000 + uint64(ind5)*100 + uint64(ind6)*10 + uint64(ind7)
			return
		}
		i++
		ind9 := intDigits[buf[i]]
		ret = ret*10000000 + uint64(ind2)*1000000 + uint64(ind3)*100000 + uint64(ind4)*10000 + uint64(ind5)*1000 + uint64(ind6)*100 + uint64(ind7)*10 + uint64(ind8)
		n = i
		if ind9 == invalidCharForNumber {
			err = assertInteger(buf[i:])
			return
		}
	}
	for i := n; i < len(buf); i++ {
		ind = intDigits[buf[i]]
		if ind == invalidCharForNumber {
			n = i
			err = assertInteger(buf[i:])
			return
		}
		if ret > uint64SafeToMultiple10 {
			value2 := (ret << 3) + (ret << 1) + uint64(ind)
			if value2 < ret {
				err = errors.New("ReadUint64ForString overflow")
				return
			}
			ret = value2
			n = i
			continue
		}
		ret = (ret << 3) + (ret << 1) + uint64(ind)
	}
	n = len(buf)
	return
}

func assertInteger(buf []byte) (err error) {
	if len(buf) > 0 && buf[0] == '.' {
		err = errors.New("assertInteger can not decode float as int")
	}
	return
}

//Write Int To String--------------------------------------------------------------------------------------------------------------------------------
func WriteChar(buff *[]byte, c byte) {
	*buff = append(*buff, c)
}

// WriteUint8 write uint8 to stream
func WriteUint8ToString(buf *[]byte, val uint8) {
	*buf = writeFirstBuf(*buf, digits[val])
	return
}

// WriteInt8 write int8 to stream
func WriteInt8ToString(buf *[]byte, nval int8) {
	var val uint8
	if nval < 0 {
		val = uint8(-nval)
		*buf = append(*buf, '-')
	} else {
		val = uint8(nval)
	}
	*buf = writeFirstBuf(*buf, digits[val])
	return
}

// WriteUint16 write uint16 to stream
func WriteUint16ToString(buf *[]byte, val uint16) {
	q1 := val / 1000
	if q1 == 0 {
		*buf = writeFirstBuf(*buf, digits[val])
		return
	}
	r1 := val - q1*1000
	*buf = writeFirstBuf(*buf, digits[q1])
	*buf = writeBuf(*buf, digits[r1])
	return
}

// WriteInt16 write int16 to stream
func WriteInt16ToString(buf *[]byte, nval int16) {
	var val uint16
	if nval < 0 {
		val = uint16(-nval)
		*buf = append(*buf, '-')
	} else {
		val = uint16(nval)
	}
	WriteUint16ToString(buf, val)
	return
}

// WriteUint32 write uint32 to stream
func WriteUint32ToString(buf *[]byte, val uint32) {
	q1 := val / 1000
	if q1 == 0 {
		*buf = writeFirstBuf(*buf, digits[val])
		return
	}
	r1 := val - q1*1000
	q2 := q1 / 1000
	if q2 == 0 {
		*buf = writeFirstBuf(*buf, digits[q1])
		*buf = writeBuf(*buf, digits[r1])
		return
	}
	r2 := q1 - q2*1000
	q3 := q2 / 1000
	if q3 == 0 {
		*buf = writeFirstBuf(*buf, digits[q2])
	} else {
		r3 := q2 - q3*1000
		*buf = append(*buf, byte(q3+'0'))
		*buf = writeBuf(*buf, digits[r3])
	}
	*buf = writeBuf(*buf, digits[r2])
	*buf = writeBuf(*buf, digits[r1])
	return
}

// WriteInt32 write int32 to stream
func WriteInt32ToString(buf *[]byte, nval int32) {
	var val uint32
	if nval < 0 {
		val = uint32(-nval)
		*buf = append(*buf, '-')
	} else {
		val = uint32(nval)
	}
	WriteUint32ToString(buf, val)
	return
}

// WriteUint64 write uint64 to stream
func WriteUint64ToString(buf *[]byte, val uint64) {
	q1 := val / 1000
	if q1 == 0 {
		*buf = writeFirstBuf(*buf, digits[val])
		return
	}
	r1 := val - q1*1000
	q2 := q1 / 1000
	if q2 == 0 {
		*buf = writeFirstBuf(*buf, digits[q1])
		*buf = writeBuf(*buf, digits[r1])
		return
	}
	r2 := q1 - q2*1000
	q3 := q2 / 1000
	if q3 == 0 {
		*buf = writeFirstBuf(*buf, digits[q2])
		*buf = writeBuf(*buf, digits[r2])
		*buf = writeBuf(*buf, digits[r1])
		return
	}
	r3 := q2 - q3*1000
	q4 := q3 / 1000
	if q4 == 0 {
		*buf = writeFirstBuf(*buf, digits[q3])
		*buf = writeBuf(*buf, digits[r3])
		*buf = writeBuf(*buf, digits[r2])
		*buf = writeBuf(*buf, digits[r1])
		return
	}
	r4 := q3 - q4*1000
	q5 := q4 / 1000
	if q5 == 0 {
		*buf = writeFirstBuf(*buf, digits[q4])
		*buf = writeBuf(*buf, digits[r4])
		*buf = writeBuf(*buf, digits[r3])
		*buf = writeBuf(*buf, digits[r2])
		*buf = writeBuf(*buf, digits[r1])
		return
	}
	r5 := q4 - q5*1000
	q6 := q5 / 1000
	if q6 == 0 {
		*buf = writeFirstBuf(*buf, digits[q5])
	} else {
		*buf = writeFirstBuf(*buf, digits[q6])
		r6 := q5 - q6*1000
		*buf = writeBuf(*buf, digits[r6])
	}
	*buf = writeBuf(*buf, digits[r5])
	*buf = writeBuf(*buf, digits[r4])
	*buf = writeBuf(*buf, digits[r3])
	*buf = writeBuf(*buf, digits[r2])
	*buf = writeBuf(*buf, digits[r1])
	return
}

// WriteInt64 write int64 to stream
func WriteInt64ToString(buf *[]byte, nval int64) {
	var val uint64
	if nval < 0 {
		val = uint64(-nval)
		*buf = append(*buf, '-')
	} else {
		val = uint64(nval)
	}
	WriteUint64ToString(buf, val)
	return
}

// WriteInt write int to stream
func WriteIntToString(buf *[]byte, val int) {
	WriteInt64ToString(buf, int64(val))
	return
}

// WriteUint write uint to stream
func WriteUint(buf *[]byte, val uint) {
	WriteUint64ToString(buf, uint64(val))
	return
}

func writeFirstBuf(space []byte, v uint32) []byte {
	start := v >> 24
	if start == 0 {
		space = append(space, byte(v>>16), byte(v>>8))
	} else if start == 1 {
		space = append(space, byte(v>>8))
	}
	space = append(space, byte(v))
	return space
}

func writeBuf(buf []byte, v uint32) []byte {
	return append(buf, byte(v>>16), byte(v>>8), byte(v))
}
