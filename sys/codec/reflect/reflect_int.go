package reflect

import (
	"errors"
	"math/big"
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

//Read-----------------------------------------------------------------------------------------------------------------------------------
func ReadBigInt(buf []byte) (ret *big.Int, n int, err error) {
	var str string
	if str, n, err = readNumberAsString(buf); err == nil {
		return
	}
	ret = big.NewInt(0)
	var success bool
	ret, success = ret.SetString(str, 10)
	if !success {
		err = errors.New("ReadBigInt invalid big int")
	}
	return
}

//Write--------------------------------------------------------------------------------------------------------------------------------
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

// WriteUint8 write uint8 to stream
func WriteUint8(buf *[]byte, val uint8) (err error) {
	*buf = writeFirstBuf(*buf, digits[val])
	return
}

// WriteInt8 write int8 to stream
func WriteInt8(buf *[]byte, nval int8) (err error) {
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
func WriteUint16(buf *[]byte, val uint16) (err error) {
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
func WriteInt16(buf *[]byte, nval int16) (err error) {
	var val uint16
	if nval < 0 {
		val = uint16(-nval)
		*buf = append(*buf, '-')
	} else {
		val = uint16(nval)
	}
	err = WriteUint16(buf, val)
	return
}

// WriteUint32 write uint32 to stream
func WriteUint32(buf *[]byte, val uint32) (err error) {
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
func WriteInt32(buf *[]byte, nval int32) (err error) {
	var val uint32
	if nval < 0 {
		val = uint32(-nval)
		*buf = append(*buf, '-')
	} else {
		val = uint32(nval)
	}
	err = WriteUint32(buf, val)
	return
}

// WriteUint64 write uint64 to stream
func WriteUint64(buf *[]byte, val uint64) (err error) {
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
func WriteInt64(buf *[]byte, nval int64) (err error) {
	var val uint64
	if nval < 0 {
		val = uint64(-nval)
		*buf = append(*buf, '-')
	} else {
		val = uint64(nval)
	}
	err = WriteUint64(buf, val)
	return
}

// WriteInt write int to stream
func WriteInt(buf *[]byte, val int) (err error) {
	WriteInt64(buf, int64(val))
	return
}

// WriteUint write uint to stream
func WriteUint(buf *[]byte, val uint) (err error) {
	WriteUint64(buf, uint64(val))
	return
}
