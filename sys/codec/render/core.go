package render

import "github.com/liwei1dao/lego/sys/codec/core"

const maxDepth = 10000

var hexDigits []byte
var valueTypes []core.ValueType

func init() {
	hexDigits = make([]byte, 256)
	for i := 0; i < len(hexDigits); i++ {
		hexDigits[i] = 255
	}
	for i := '0'; i <= '9'; i++ {
		hexDigits[i] = byte(i - '0')
	}
	for i := 'a'; i <= 'f'; i++ {
		hexDigits[i] = byte((i - 'a') + 10)
	}
	for i := 'A'; i <= 'F'; i++ {
		hexDigits[i] = byte((i - 'A') + 10)
	}
	valueTypes = make([]core.ValueType, 256)
	for i := 0; i < len(valueTypes); i++ {
		valueTypes[i] = core.InvalidValue
	}
	valueTypes['"'] = core.StringValue
	valueTypes['-'] = core.NumberValue
	valueTypes['0'] = core.NumberValue
	valueTypes['1'] = core.NumberValue
	valueTypes['2'] = core.NumberValue
	valueTypes['3'] = core.NumberValue
	valueTypes['4'] = core.NumberValue
	valueTypes['5'] = core.NumberValue
	valueTypes['6'] = core.NumberValue
	valueTypes['7'] = core.NumberValue
	valueTypes['8'] = core.NumberValue
	valueTypes['9'] = core.NumberValue
	valueTypes['t'] = core.BoolValue
	valueTypes['f'] = core.BoolValue
	valueTypes['n'] = core.NilValue
	valueTypes['['] = core.ArrayValue
	valueTypes['{'] = core.ObjectValue
}
