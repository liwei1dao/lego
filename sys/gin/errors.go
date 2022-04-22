package gin

import (
	"fmt"
	"strings"
)

type ErrorType uint64

const (
	//当 Context.Bind() 失败时使用 ErrorTypeBind。
	ErrorTypeBind ErrorType = 1 << 63
	//当 Context.Render() 失败时使用 ErrorTypeRender。
	ErrorTypeRender ErrorType = 1 << 62
	//ErrorTypePrivate 表示私有错误。
	ErrorTypePrivate ErrorType = 1 << 0
	//ErrorTypePublic 表示公共错误。
	ErrorTypePublic ErrorType = 1 << 1
	//ErrorTypeAny 表示任何其他错误。
	ErrorTypeAny ErrorType = 1<<64 - 1
	//ErrorTypeNu 表示任何其他错误。
	ErrorTypeNu = 2
)

type Error struct {
	Err  error
	Type ErrorType
	Meta interface{}
}

func (msg *Error) IsType(flags ErrorType) bool {
	return (msg.Type & flags) > 0
}

type errorMsgs []*Error

func (this errorMsgs) ByType(typ ErrorType) errorMsgs {
	if len(this) == 0 {
		return nil
	}
	if typ == ErrorTypeAny {
		return this
	}
	var result errorMsgs
	for _, msg := range this {
		if msg.IsType(typ) {
			result = append(result, msg)
		}
	}
	return result
}

func (this errorMsgs) String() string {
	if len(this) == 0 {
		return ""
	}
	var buffer strings.Builder
	for i, msg := range this {
		fmt.Fprintf(&buffer, "Error #%02d: %s\n", i+1, msg.Err)
		if msg.Meta != nil {
			fmt.Fprintf(&buffer, "     Meta: %v\n", msg.Meta)
		}
	}
	return buffer.String()
}
