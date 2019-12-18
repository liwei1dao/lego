package event

import (
	"reflect"
)

type FunctionInfo struct {
	Function  reflect.Value
	Goroutine bool
}
