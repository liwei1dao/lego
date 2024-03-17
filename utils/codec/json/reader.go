package json

import (
	"encoding/json"
	"fmt"

	"io"
	"math/big"
	"reflect"
	"unicode/utf16"

	"github.com/liwei1dao/lego/utils/codec"
	"github.com/liwei1dao/lego/utils/codec/codecore"
	"github.com/liwei1dao/lego/utils/codec/utils"

	"github.com/modern-go/reflect2"
)

type JsonReader struct {
	config           *codecore.Config
	reader           io.Reader
	buf              []byte
	head             int
	tail             int
	depth            int
	captureStartedAt int
	captured         []byte
	err              error
}

func (this *JsonReader) Config() *codecore.Config {
	return this.config
}
func (this *JsonReader) GetReader(buf []byte, r io.Reader) codecore.IReader {
	return GetReader(buf, r)
}
func (this *JsonReader) PutReader(r codecore.IReader) {
	PutReader(r)
}
func (this *JsonReader) GetWriter() codecore.IWriter {
	return GetWriter(nil)
}
func (this *JsonReader) PutWriter(w codecore.IWriter) {
	PutWriter(w)
}
func (this *JsonReader) ReadVal(obj interface{}) {
	depth := this.depth
	cacheKey := reflect2.RTypeOf(obj)
	decoder := codec.GetDecoder(cacheKey)
	if decoder == nil {
		typ := reflect2.TypeOf(obj)
		if typ == nil || typ.Kind() != reflect.Ptr {
			this.reportError("ReadVal", "can only unmarshal into pointer")
			return
		}
		decoder = codec.DecoderOf(typ, this.config)
	}
	ptr := reflect2.PtrOf(obj)
	if ptr == nil {
		this.reportError("ReadVal", "can not read into nil pointer")
		return
	}
	decoder.Decode(ptr, this)
	if this.depth != depth {
		this.reportError("ReadVal", "unexpected mismatched nesting")
		return
	}
}
func (this *JsonReader) WhatIsNext() codecore.ValueType {
	valueType := valueTypes[this.nextToken()]
	this.unreadByte()
	return valueType
}
func (this *JsonReader) Read() interface{} {
	valueType := this.WhatIsNext()
	switch valueType {
	case codecore.StringValue:
		return this.ReadString()
	case codecore.NumberValue:
		if this.config.UseNumber {
			ret, n, err := utils.ReadNumberAsString(this.buf[this.head:])
			if err != nil {
				this.err = err
			}
			this.head += n
			return json.Number(ret)
		}
		return this.ReadFloat64()
	case codecore.NilValue:
		this.skipFourBytes('n', 'u', 'l', 'l')
		return nil
	case codecore.BoolValue:
		return this.ReadBool()
	case codecore.ArrayValue:
		arr := []interface{}{}
		this.ReadArrayCB(func(extra codecore.IReader) bool {
			var elem interface{}
			extra.ReadVal(&elem)
			arr = append(arr, elem)
			return true
		})
		return arr
	case codecore.ObjectValue:
		obj := map[string]interface{}{}
		this.ReadMapCB(func(extra codecore.IReader, field string) bool {
			var elem interface{}
			this.ReadVal(&elem)
			obj[field] = elem
			return true
		})
		return obj
	default:
		this.reportError("Read", fmt.Sprintf("unexpected value type: %v", valueType))
		return nil
	}
}
func (this *JsonReader) ReadNil() (ret bool) {
	c := this.nextToken()
	if c == 'n' {
		this.skipThreeBytes('u', 'l', 'l') // null
		return true
	}
	this.unreadByte()
	return false
}
func (this *JsonReader) ReadArrayStart() (ret bool) {
	c := this.nextToken()
	if c == '[' {
		return true
	}
	this.reportError("ReadArrayStart", `expect [ but found `+string([]byte{c}))
	return
}
func (this *JsonReader) CheckNextIsArrayEnd() (ret bool) {
	c := this.nextToken()
	if c == ']' {
		return true
	}
	this.unreadByte()
	return
}
func (this *JsonReader) ReadArrayEnd() (ret bool) {
	c := this.nextToken()
	if c == ']' {
		return true
	}
	this.reportError("ReadArrayEnd", `expect ] but found `+string([]byte{c}))
	return
}
func (this *JsonReader) ReadObjectStart() (ret bool) {
	c := this.nextToken()
	if c == '{' {
		return this.incrementDepth()
	}
	this.reportError("ReadObjectStart", `expect { but found `+string([]byte{c}))
	return
}
func (this *JsonReader) CheckNextIsObjectEnd() (ret bool) {
	c := this.nextToken()
	if c == '}' {
		return this.decrementDepth()
	}
	this.unreadByte()
	return
}
func (this *JsonReader) ReadObjectEnd() (ret bool) {
	c := this.nextToken()
	if c == '}' {
		return this.decrementDepth()
	}
	this.reportError("ReadObjectEnd", `expect } but found `+string([]byte{c}))
	return
}
func (this *JsonReader) ReadMemberSplit() (ret bool) {
	c := this.nextToken()
	if c == ',' {
		return true
	}
	this.unreadByte()
	return
}
func (this *JsonReader) ReadKVSplit() (ret bool) {
	c := this.nextToken()
	if c == ':' {
		return true
	}
	this.reportError("ReadKVSplit", `expect : but found `+string([]byte{c}))
	return
}
func (this *JsonReader) ReadKeyStart() (ret bool) {
	c := this.nextToken()
	if c == '"' {
		return true
	}
	this.reportError("ReadKeyStart", `expect " but found `+string([]byte{c}))
	return
}
func (this *JsonReader) ReadKeyEnd() (ret bool) {
	c := this.nextToken()
	if c == '"' {
		return true
	}
	this.reportError("ReadKeyEnd", `expect " but found `+string([]byte{c}))
	return
}

func (this *JsonReader) SkipAndReturnBytes() []byte {
	this.nextToken()
	this.unreadByte()
	this.startCapture(this.head)
	this.Skip()
	return this.stopCapture()
}

func (this *JsonReader) Skip() {
	c := this.nextToken()
	switch c {
	case '"':
		this.skipString()
	case 'n':
		this.skipThreeBytes('u', 'l', 'l') // null
	case 't':
		this.skipThreeBytes('r', 'u', 'e') // true
	case 'f':
		this.skipFourBytes('a', 'l', 's', 'e') // false
	case '0':
		this.unreadByte()
		this.ReadFloat32()
	case '-', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		this.skipNumber()
	case '[':
		this.skipArray()
	case '{':
		this.skipObject()
	default:
		this.reportError("Skip", fmt.Sprintf("do not know how to skip: %v", c))
		return
	}
}
func (this *JsonReader) ReadBool() (ret bool) {
	c := this.nextToken()
	if c == 't' {
		this.skipThreeBytes('r', 'u', 'e')
		return true
	}
	if c == 'f' {
		this.skipFourBytes('a', 'l', 's', 'e')
		return false
	}
	this.reportError("ReadBool", "expect t or f, but found "+string([]byte{c}))
	return
}
func (this *JsonReader) ReadInt8() (ret int8) {
	var (
		n   int
		err error
	)
	if ret, n, err = utils.ReadInt8ForString(this.buf[this.head:]); err != nil {
		this.reportError("ReadInt8", err.Error())
		return
	}
	this.head += n
	return
}
func (this *JsonReader) ReadInt16() (ret int16) {
	var (
		n   int
		err error
	)
	if ret, n, err = utils.ReadInt16ForString(this.buf[this.head:]); err != nil {
		this.reportError("ReadInt32", err.Error())
		return
	}
	this.head += n
	return
}
func (this *JsonReader) ReadInt32() (ret int32) {
	var (
		n   int
		err error
	)
	if ret, n, err = utils.ReadInt32ForString(this.buf[this.head:]); err != nil {
		this.reportError("ReadInt32", err.Error())
		return
	}
	this.head += n
	return
}
func (this *JsonReader) ReadInt64() (ret int64) {
	var (
		n   int
		err error
	)
	if ret, n, err = utils.ReadInt64ForString(this.buf[this.head:]); err != nil {
		this.reportError("ReadInt64", err.Error())
		return
	}
	this.head += n
	return
}
func (this *JsonReader) ReadUint8() (ret uint8) {
	var (
		n   int
		err error
	)
	if ret, n, err = utils.ReadUint8ForString(this.buf[this.head:]); err != nil {
		this.reportError("ReadUint8", err.Error())
		return
	}
	this.head += n
	return
}
func (this *JsonReader) ReadUint16() (ret uint16) {
	var (
		n   int
		err error
	)
	if ret, n, err = utils.ReadUint16ForString(this.buf[this.head:]); err != nil {
		this.reportError("ReadUint16", err.Error())
		return
	}
	this.head += n
	return
}
func (this *JsonReader) ReadUint32() (ret uint32) {
	var (
		n   int
		err error
	)
	if ret, n, err = utils.ReadUint32ForString(this.buf[this.head:]); err != nil {
		this.reportError("ReadUint32", err.Error())
		return
	}
	this.head += n
	return
}
func (this *JsonReader) ReadUint64() (ret uint64) {
	var (
		n   int
		err error
	)
	if ret, n, err = utils.ReadUint64ForString(this.buf[this.head:]); err != nil {
		this.reportError("ReadUint64", err.Error())
		return
	}
	this.head += n
	return
}
func (this *JsonReader) ReadFloat32() (ret float32) {
	var (
		n   int
		err error
	)
	if ret, n, err = utils.ReadFloat32ForString(this.buf[this.head:]); err != nil {
		this.reportError("ReadFloat32", err.Error())
		return
	}
	this.head += n
	return
}
func (this *JsonReader) ReadFloat64() (ret float64) {
	var (
		n   int
		err error
	)
	if ret, n, err = utils.ReadFloat64ForString(this.buf[this.head:]); err != nil {
		this.reportError("ReadFloat64", err.Error())
		return
	}
	this.head += n
	return
}
func (this *JsonReader) ReadString() (ret string) {
	c := this.nextToken()
	if c == '"' {
		for i := this.head; i < this.tail; i++ {
			c := this.buf[i]
			if c == '"' {
				ret = string(this.buf[this.head:i])
				this.head = i + 1
				return ret
			} else if c == '\\' {
				break
			} else if c < ' ' {
				this.reportError("ReadString",
					fmt.Sprintf(`invalid control character found: %d`, c))
				return
			}
		}
		return this.readStringSlowPath()
	} else if c == 'n' {
		this.skipThreeBytes('u', 'l', 'l')
		return ""
	}
	this.reportError("ReadString", `expects " or n, but found `+string([]byte{c}))
	return
}
func (this *JsonReader) ResetBytes(d []byte, r io.Reader) {
	this.buf = d
	this.reader = r
	this.head = 0
	this.tail = len(d)
	if this.reader != nil {
		if !this.loadMore() {
			this.err = io.EOF
		}
	}
}
func (this *JsonReader) Error() error {
	return this.err
}
func (this *JsonReader) SetErr(err error) {
	this.err = err
}

//-----------------------------------------------------------------------------------------------------------------------------------
func (this *JsonReader) readByte() (ret byte) {
	if this.head == this.tail {
		if this.loadMore() {
			ret = this.buf[this.head]
			this.head++
			return ret
		}
		return 0
	}
	ret = this.buf[this.head]
	this.head++
	return ret
}
func (this *JsonReader) readStringSlowPath() (ret string) {
	var str []byte
	var c byte
	for this.err == nil {
		c = this.readByte()
		if c == '"' {
			return string(str)
		}
		if c == '\\' {
			c = this.readByte()
			str = this.readEscapedChar(c, str)
		} else {
			str = append(str, c)
		}
	}
	this.reportError("readStringSlowPath", "unexpected end of input")
	return
}
func (this *JsonReader) readEscapedChar(c byte, str []byte) []byte {
	switch c {
	case 'u':
		r := this.readU4()
		if utf16.IsSurrogate(r) {
			c = this.readByte()
			if this.err != nil {
				return nil
			}
			if c != '\\' {
				this.unreadByte()
				str = utils.AppendRune(str, r)
				return str
			}
			c = this.readByte()
			if this.err != nil {
				return nil
			}
			if c != 'u' {
				str = utils.AppendRune(str, r)
				return this.readEscapedChar(c, str)
			}
			r2 := this.readU4()
			if this.err != nil {
				return nil
			}
			combined := utf16.DecodeRune(r, r2)
			if combined == '\uFFFD' {
				str = utils.AppendRune(str, r)
				str = utils.AppendRune(str, r2)
			} else {
				str = utils.AppendRune(str, combined)
			}
		} else {
			str = utils.AppendRune(str, r)
		}
	case '"':
		str = append(str, '"')
	case '\\':
		str = append(str, '\\')
	case '/':
		str = append(str, '/')
	case 'b':
		str = append(str, '\b')
	case 'f':
		str = append(str, '\f')
	case 'n':
		str = append(str, '\n')
	case 'r':
		str = append(str, '\r')
	case 't':
		str = append(str, '\t')
	default:
		this.reportError("readEscapedChar",
			`invalid escape char after \`)
		return nil
	}
	return str
}
func (this *JsonReader) readU4() (ret rune) {
	for i := 0; i < 4; i++ {
		c := this.readByte()
		if this.err != nil {
			return
		}
		if c >= '0' && c <= '9' {
			ret = ret*16 + rune(c-'0')
		} else if c >= 'a' && c <= 'f' {
			ret = ret*16 + rune(c-'a'+10)
		} else if c >= 'A' && c <= 'F' {
			ret = ret*16 + rune(c-'A'+10)
		} else {
			this.reportError("readU4", "expects 0~9 or a~f, but found "+string([]byte{c}))
			return
		}
	}
	return ret
}
func (iter *JsonReader) loadMore() bool {
	if iter.reader == nil {
		if iter.Error == nil {
			iter.head = iter.tail
			iter.err = io.EOF
		}
		return false
	}
	if iter.captured != nil {
		iter.captured = append(iter.captured,
			iter.buf[iter.captureStartedAt:iter.tail]...)
		iter.captureStartedAt = 0
	}
	for {
		n, err := iter.reader.Read(iter.buf)
		if n == 0 {
			if err != nil {
				if iter.err == nil {
					iter.err = err
				}
				return false
			}
		} else {
			iter.head = 0
			iter.tail = n
			return true
		}
	}
}

func (this *JsonReader) nextToken() byte {
	for {
		for i := this.head; i < this.tail; i++ {
			c := this.buf[i]
			switch c {
			case ' ', '\n', '\t', '\r':
				continue
			}
			this.head = i + 1
			return c
		}
		if !this.loadMore() {
			return 0
		}
	}
}
func (this *JsonReader) unreadByte() {
	if this.err != nil {
		return
	}
	this.head--
	return
}
func (this *JsonReader) startCapture(captureStartedAt int) {
	this.startCaptureTo(make([]byte, 0, 32), captureStartedAt)
}
func (this *JsonReader) startCaptureTo(buf []byte, captureStartedAt int) {
	if this.captured != nil {
		panic("already in capture mode")
	}
	this.captureStartedAt = captureStartedAt
	this.captured = buf
}
func (this *JsonReader) stopCapture() []byte {
	if this.captured == nil {
		panic("not in capture mode")
	}
	captured := this.captured
	remaining := this.buf[this.captureStartedAt:this.head]
	this.captureStartedAt = -1
	this.captured = nil
	return append(captured, remaining...)
}
func (this *JsonReader) skipNumber() {
	if !this.trySkipNumber() {
		this.unreadByte()
		if this.err != nil && this.err != io.EOF {
			return
		}
		this.ReadFloat64()
		if this.err != nil && this.err != io.EOF {
			this.err = nil
			this.ReadBigFloat()
		}
	}
}
func (this *JsonReader) ReadBigFloat() (ret *big.Float) {
	var (
		n   int
		err error
	)
	if ret, n, err = utils.ReadBigFloatForString(this.buf[this.head:]); err != nil {
		this.reportError("ReadBigFloat", err.Error())
		return
	}
	this.head += n
	return
}
func (iter *JsonReader) trySkipNumber() bool {
	dotFound := false
	for i := iter.head; i < iter.tail; i++ {
		c := iter.buf[i]
		switch c {
		case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		case '.':
			if dotFound {
				iter.reportError("validateNumber", `more than one dot found in number`)
				return true // already failed
			}
			if i+1 == iter.tail {
				return false
			}
			c = iter.buf[i+1]
			switch c {
			case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
			default:
				iter.reportError("validateNumber", `missing digit after dot`)
				return true // already failed
			}
			dotFound = true
		default:
			switch c {
			case ',', ']', '}', ' ', '\t', '\n', '\r':
				if iter.head == i {
					return false // if - without following digits
				}
				iter.head = i
				return true // must be valid
			}
			return false // may be invalid
		}
	}
	return false
}
func (this *JsonReader) skipString() {
	if !this.trySkipString() {
		this.unreadByte()
		this.ReadString()
	}
}
func (this *JsonReader) trySkipString() bool {
	for i := this.head; i < this.tail; i++ {
		c := this.buf[i]
		if c == '"' {
			this.head = i + 1
			return true // valid
		} else if c == '\\' {
			return false
		} else if c < ' ' {
			this.reportError("trySkipString",
				fmt.Sprintf(`invalid control character found: %d`, c))
			return true // already failed
		}
	}
	return false
}
func (this *JsonReader) skipObject() {
	this.unreadByte()
	this.ReadObjectCB(func(extra codecore.IReader, field string) bool {
		extra.Skip()
		return true
	})
}
func (this *JsonReader) skipArray() {
	this.unreadByte()
	this.ReadArrayCB(func(extra codecore.IReader) bool {
		extra.Skip()
		return true
	})
}
func (this *JsonReader) skipThreeBytes(b1, b2, b3 byte) {
	if this.readByte() != b1 {
		this.reportError("skipThreeBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3})))
		return
	}
	if this.readByte() != b2 {
		this.reportError("skipThreeBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3})))
		return
	}
	if this.readByte() != b3 {
		this.reportError("skipThreeBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3})))
		return
	}
}
func (this *JsonReader) skipFourBytes(b1, b2, b3, b4 byte) {
	if this.readByte() != b1 {
		this.reportError("skipFourBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3, b4})))
		return
	}
	if this.readByte() != b2 {
		this.reportError("skipFourBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3, b4})))
		return
	}
	if this.readByte() != b3 {
		this.reportError("skipFourBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3, b4})))
		return
	}
	if this.readByte() != b4 {
		this.reportError("skipFourBytes", fmt.Sprintf("expect %s", string([]byte{b1, b2, b3, b4})))
		return
	}
}
func (this *JsonReader) reportError(operation string, msg string) {
	if this.err != nil {
		if this.err != io.EOF {
			return
		}
	}
	peekStart := this.head - 10
	if peekStart < 0 {
		peekStart = 0
	}
	peekEnd := this.head + 10
	if peekEnd > this.tail {
		peekEnd = this.tail
	}
	parsing := string(this.buf[peekStart:peekEnd])
	contextStart := this.head - 50
	if contextStart < 0 {
		contextStart = 0
	}
	contextEnd := this.head + 50
	if contextEnd > this.tail {
		contextEnd = this.tail
	}
	context := string(this.buf[contextStart:contextEnd])
	this.err = fmt.Errorf("%s: %s, error found in #%v byte of ...|%s|..., bigger context ...|%s|...",
		operation, msg, this.head-peekStart, parsing, context)
}
func (this *JsonReader) incrementDepth() (success bool) {
	this.depth++
	if this.depth <= maxDepth {
		return true
	}
	this.reportError("incrementDepth", "exceeded max depth")
	return false
}
func (this *JsonReader) decrementDepth() (success bool) {
	this.depth--
	if this.depth >= 0 {
		return true
	}
	this.reportError("decrementDepth", "unexpected negative nesting")
	return false
}
func (iter *JsonReader) ReadObjectCB(callback func(codecore.IReader, string) bool) bool {
	c := iter.nextToken()
	var field string
	if c == '{' {
		if !iter.incrementDepth() {
			return false
		}
		c = iter.nextToken()
		if c == '"' {
			iter.unreadByte()
			field = iter.ReadString()
			c = iter.nextToken()
			if c != ':' {
				iter.reportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
			}
			if !callback(iter, field) {
				iter.decrementDepth()
				return false
			}
			c = iter.nextToken()
			for c == ',' {
				field = iter.ReadString()
				c = iter.nextToken()
				if c != ':' {
					iter.reportError("ReadObject", "expect : after object field, but found "+string([]byte{c}))
				}
				if !callback(iter, field) {
					iter.decrementDepth()
					return false
				}
				c = iter.nextToken()
			}
			if c != '}' {
				iter.reportError("ReadObjectCB", `object not ended with }`)
				iter.decrementDepth()
				return false
			}
			return iter.decrementDepth()
		}
		if c == '}' {
			return iter.decrementDepth()
		}
		iter.reportError("ReadObjectCB", `expect " after {, but found `+string([]byte{c}))
		iter.decrementDepth()
		return false
	}
	if c == 'n' {
		iter.skipThreeBytes('u', 'l', 'l')
		return true // null
	}
	iter.reportError("ReadObjectCB", `expect { or n, but found `+string([]byte{c}))
	return false
}

func (iter *JsonReader) ReadMapCB(callback func(codecore.IReader, string) bool) bool {
	c := iter.nextToken()
	if c == '{' {
		if !iter.incrementDepth() {
			return false
		}
		c = iter.nextToken()
		if c == '"' {
			iter.unreadByte()
			field := iter.ReadString()
			if iter.nextToken() != ':' {
				iter.reportError("ReadMapCB", "expect : after object field, but found "+string([]byte{c}))
				iter.decrementDepth()
				return false
			}
			if !callback(iter, field) {
				iter.decrementDepth()
				return false
			}
			c = iter.nextToken()
			for c == ',' {
				field = iter.ReadString()
				if iter.nextToken() != ':' {
					iter.reportError("ReadMapCB", "expect : after object field, but found "+string([]byte{c}))
					iter.decrementDepth()
					return false
				}
				if !callback(iter, field) {
					iter.decrementDepth()
					return false
				}
				c = iter.nextToken()
			}
			if c != '}' {
				iter.reportError("ReadMapCB", `object not ended with }`)
				iter.decrementDepth()
				return false
			}
			return iter.decrementDepth()
		}
		if c == '}' {
			return iter.decrementDepth()
		}
		iter.reportError("ReadMapCB", `expect " after {, but found `+string([]byte{c}))
		iter.decrementDepth()
		return false
	}
	if c == 'n' {
		iter.skipThreeBytes('u', 'l', 'l')
		return true // null
	}
	iter.reportError("ReadMapCB", `expect { or n, but found `+string([]byte{c}))
	return false
}

func (iter *JsonReader) ReadArrayCB(callback func(codecore.IReader) bool) (ret bool) {
	c := iter.nextToken()
	if c == '[' {
		if !iter.incrementDepth() {
			return false
		}
		c = iter.nextToken()
		if c != ']' {
			iter.unreadByte()
			if !callback(iter) {
				iter.decrementDepth()
				return false
			}
			c = iter.nextToken()
			for c == ',' {
				if !callback(iter) {
					iter.decrementDepth()
					return false
				}
				c = iter.nextToken()
			}
			if c != ']' {
				iter.reportError("ReadArrayCB", "expect ] in the end, but found "+string([]byte{c}))
				iter.decrementDepth()
				return false
			}
			return iter.decrementDepth()
		}
		return iter.decrementDepth()
	}
	if c == 'n' {
		iter.skipThreeBytes('u', 'l', 'l')
		return true // null
	}
	iter.reportError("ReadArrayCB", "expect [ or n, but found "+string([]byte{c}))
	return false
}
