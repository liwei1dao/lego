package utils

import (
	"errors"
	"fmt"
	"math"
	"math/big"
	"strconv"
	"strings"
	"unsafe"
)

var pow10 []uint64
var floatDigits []int8

const invalidCharForNumber = int8(-1)
const endOfNumber = int8(-2)
const dotInNumber = int8(-3)

func init() {
	pow10 = []uint64{1, 10, 100, 1000, 10000, 100000, 1000000}
	floatDigits = make([]int8, 256)
	for i := 0; i < len(floatDigits); i++ {
		floatDigits[i] = invalidCharForNumber
	}
	for i := int8('0'); i <= int8('9'); i++ {
		floatDigits[i] = i - int8('0')
	}
	floatDigits[','] = endOfNumber
	floatDigits[']'] = endOfNumber
	floatDigits['}'] = endOfNumber
	floatDigits[' '] = endOfNumber
	floatDigits['\t'] = endOfNumber
	floatDigits['\n'] = endOfNumber
	floatDigits['.'] = dotInNumber
}

//Read Float For String-------------------------------------------------------------------------------------------------------------------------------------------------------------
//从字符串中读取float 对象
func ReadBigFloatForString(buf []byte) (ret *big.Float, n int, err error) {
	var str string
	if str, n, err = readNumberAsString(buf); err != nil {
		return
	}
	prec := 64
	if len(str) > prec {
		prec = len(str)
	}
	ret, _, err = big.ParseFloat(str, 10, uint(prec), big.ToZero)
	return
}
func ReadFloat32ForString(buf []byte) (ret float32, n int, err error) {
	if buf[0] == '-' {
		ret, n, err = readPositiveFloat32(buf[1:])
		n++
		ret = -ret
		return
	}
	return readPositiveFloat32(buf)
}
func ReadFloat64ForString(buf []byte) (ret float64, n int, err error) {
	if buf[0] == '-' {
		ret, n, err = readPositiveFloat64(buf[1:])
		n++
		ret = -ret
		return
	}
	return readPositiveFloat64(buf)
}

//读取float32
func readPositiveFloat32(buf []byte) (ret float32, n int, err error) {
	i := 0
	// first char
	if len(buf) == 0 {
		return readFloat32SlowPath(buf)
	}
	c := buf[i]
	i++
	ind := floatDigits[c]
	switch ind {
	case invalidCharForNumber:
		return readFloat32SlowPath(buf)
	case endOfNumber:
		err = errors.New("readFloat32 empty number")
		return
	case dotInNumber:
		err = errors.New("readFloat32 leading dot is invalid")
		return
	case 0:
		if i == len(buf) {
			return readFloat32SlowPath(buf)
		}
		c = buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			err = errors.New("readFloat32 leading zero is invalid")
			return
		}
	}
	value := uint64(ind)
	// chars before dot
non_decimal_loop:
	for ; i < len(buf); i++ {
		c = buf[i]
		ind := floatDigits[c]
		switch ind {
		case invalidCharForNumber:
			return readFloat32SlowPath(buf)
		case endOfNumber:
			n = i
			ret = float32(value)
			return
		case dotInNumber:
			break non_decimal_loop
		}
		if value > uint64SafeToMultiple10 {
			return readFloat32SlowPath(buf)
		}
		value = (value << 3) + (value << 1) + uint64(ind) // value = value * 10 + ind;
	}
	// chars after dot
	if c == '.' {
		i++
		decimalPlaces := 0
		if i == len(buf) {
			return readFloat32SlowPath(buf)
		}
		for ; i < len(buf); i++ {
			c = buf[i]
			ind := floatDigits[c]
			switch ind {
			case endOfNumber:
				if decimalPlaces > 0 && decimalPlaces < len(pow10) {
					n = i
					ret = float32(float64(value) / float64(pow10[decimalPlaces]))
					return
				}
				// too many decimal places
				return readFloat32SlowPath(buf)
			case invalidCharForNumber, dotInNumber:
				return readFloat32SlowPath(buf)
			}
			decimalPlaces++
			if value > uint64SafeToMultiple10 {
				return readFloat32SlowPath(buf)
			}
			value = (value << 3) + (value << 1) + uint64(ind)
		}
	}
	return readFloat32SlowPath(buf)
}
func readFloat32SlowPath(buf []byte) (ret float32, n int, err error) {
	var str string
	if str, n, err = readNumberAsString(buf); err != nil {
		return
	}
	errMsg := validateFloat(str)
	if errMsg != "" {
		err = errors.New(errMsg)
		return
	}
	var val float64
	if val, err = strconv.ParseFloat(str, 32); err != nil {
		return
	}
	ret = float32(val)
	return
}

//读取float64
func readPositiveFloat64(buf []byte) (ret float64, n int, err error) {
	i := 0
	// first char
	if i == len(buf) {
		return readFloat64SlowPath(buf)
	}
	c := buf[i]
	i++
	ind := floatDigits[c]
	switch ind {
	case invalidCharForNumber:
		return readFloat64SlowPath(buf)
	case endOfNumber:
		err = errors.New("readFloat64 empty number")
		return
	case dotInNumber:
		err = errors.New("readFloat64 leading dot is invalid")
		return
	case 0:
		if i == len(buf) {
			return readFloat64SlowPath(buf)
		}
		c = buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			err = errors.New("readFloat64 leading zero is invalid")
			return
		}
	}
	value := uint64(ind)
	// chars before dot
non_decimal_loop:
	for ; i < len(buf); i++ {
		c = buf[i]
		ind := floatDigits[c]
		switch ind {
		case invalidCharForNumber:
			return readFloat64SlowPath(buf)
		case endOfNumber:
			n = i
			ret = float64(value)
			return
		case dotInNumber:
			break non_decimal_loop
		}
		if value > uint64SafeToMultiple10 {
			return readFloat64SlowPath(buf)
		}
		value = (value << 3) + (value << 1) + uint64(ind) // value = value * 10 + ind;
	}
	// chars after dot
	if c == '.' {
		i++
		decimalPlaces := 0
		if i == len(buf) {
			return readFloat64SlowPath(buf)
		}
		for ; i < len(buf); i++ {
			c = buf[i]
			ind := floatDigits[c]
			switch ind {
			case endOfNumber:
				if decimalPlaces > 0 && decimalPlaces < len(pow10) {
					n = i
					ret = float64(value) / float64(pow10[decimalPlaces])
					return
				}
				// too many decimal places
				return readFloat64SlowPath(buf)
			case invalidCharForNumber, dotInNumber:
				return readFloat64SlowPath(buf)
			}
			decimalPlaces++
			if value > uint64SafeToMultiple10 {
				return readFloat64SlowPath(buf)
			}
			value = (value << 3) + (value << 1) + uint64(ind)
			if value > maxFloat64 {
				return readFloat64SlowPath(buf)
			}
		}
	}
	return readFloat64SlowPath(buf)
}

func readFloat64SlowPath(buf []byte) (ret float64, n int, err error) {
	var str string
	if str, n, err = readNumberAsString(buf); err != nil {
		return
	}
	errMsg := validateFloat(str)
	if errMsg != "" {
		err = fmt.Errorf("readFloat64SlowPath:%v", errMsg)
		return
	}
	ret, err = strconv.ParseFloat(str, 64)
	return
}

func ReadNumberAsString(buf []byte) (ret string, n int, err error) {
	return readNumberAsString(buf)
}

//读取字符串数字
func readNumberAsString(buf []byte) (ret string, n int, err error) {
	strBuf := [16]byte{}
	str := strBuf[0:0]
load_loop:
	for i := 0; i < len(buf); i++ {
		c := buf[i]
		switch c {
		case '+', '-', '.', 'e', 'E', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			str = append(str, c)
			continue
		default:
			n = i
			break load_loop
		}
	}
	if len(str) == 0 {
		err = errors.New("readNumberAsString invalid number")
	}
	ret = *(*string)(unsafe.Pointer(&str))
	return
}

//字符串float 错误检验
func validateFloat(str string) string {
	if len(str) == 0 {
		return "empty number"
	}
	if str[0] == '-' {
		return "-- is not valid"
	}
	dotPos := strings.IndexByte(str, '.')
	if dotPos != -1 {
		if dotPos == len(str)-1 {
			return "dot can not be last character"
		}
		switch str[dotPos+1] {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		default:
			return "missing digit after dot"
		}
	}
	return ""
}

//Write Float To String-------------------------------------------------------------------------------------------------------------------------------------------------------------
func WriteFloat32ToString(buff *[]byte, val float32) (err error) {
	if math.IsInf(float64(val), 0) || math.IsNaN(float64(val)) {
		err = fmt.Errorf("unsupported value: %f", val)
		return
	}
	abs := math.Abs(float64(val))
	fmt := byte('f')
	if abs != 0 {
		if float32(abs) < 1e-6 || float32(abs) >= 1e21 {
			fmt = 'e'
		}
	}
	*buff = strconv.AppendFloat(*buff, float64(val), fmt, -1, 32)
	return
}

//float32 写入只有 6 位精度的流，但速度要快得多
func WriteFloat32LossyToString(buf *[]byte, val float32) (err error) {
	if math.IsInf(float64(val), 0) || math.IsNaN(float64(val)) {
		err = fmt.Errorf("unsupported value: %f", val)
		return
	}
	if val < 0 {
		WriteChar(buf, '-')
		val = -val
	}
	if val > 0x4ffffff {
		WriteFloat32ToString(buf, val)
		return
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(float64(val)*float64(exp) + 0.5)
	WriteUint64ToString(buf, lval/exp)
	fval := lval % exp
	if fval == 0 {
		return
	}
	WriteChar(buf, '.')
	for p := precision - 1; p > 0 && fval < pow10[p]; p-- {
		WriteChar(buf, '0')
	}
	WriteUint64ToString(buf, fval)
	temp := *buf
	for temp[len(temp)-1] == '0' {
		temp = temp[:len(temp)-1]
	}
	buf = &temp
	return
}

func WriteFloat64ToString(buf *[]byte, val float64) (err error) {
	if math.IsInf(val, 0) || math.IsNaN(val) {
		err = fmt.Errorf("unsupported value: %f", val)
		return
	}
	abs := math.Abs(val)
	fmt := byte('f')
	// Note: Must use float32 comparisons for underlying float32 value to get precise cutoffs right.
	if abs != 0 {
		if abs < 1e-6 || abs >= 1e21 {
			fmt = 'e'
		}
	}
	*buf = strconv.AppendFloat(*buf, float64(val), fmt, -1, 64)
	return
}

//将 float64 写入只有 6 位精度的流，但速度要快得多
func WriteFloat64LossyToString(buf *[]byte, val float64) (err error) {
	if math.IsInf(val, 0) || math.IsNaN(val) {
		err = fmt.Errorf("unsupported value: %f", val)
		return
	}
	if val < 0 {
		WriteChar(buf, '-')
		val = -val
	}
	if val > 0x4ffffff {
		WriteFloat64ToString(buf, val)
		return
	}
	precision := 6
	exp := uint64(1000000) // 6
	lval := uint64(val*float64(exp) + 0.5)
	WriteUint64ToString(buf, lval/exp)
	fval := lval % exp
	if fval == 0 {
		return
	}
	WriteChar(buf, '.')
	for p := precision - 1; p > 0 && fval < pow10[p]; p-- {
		WriteChar(buf, '0')
	}
	WriteUint64ToString(buf, fval)
	temp := *buf
	for temp[len(temp)-1] == '0' {
		temp = temp[:len(temp)-1]
	}
	buf = &temp
	return
}
