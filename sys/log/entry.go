package log

import (
	"fmt"
	"reflect"
	"strings"
	"sync"
	"time"

	"github.com/liwei1dao/lego/utils/pools"
)

var (
	_cePool = sync.Pool{New: func() interface{} {
		// Pre-allocate some space for cores.
		return &Entry{
			Data: make(Fields, 0, 6),
		}
	}}
)

func getEntry() *Entry {
	return _cePool.Get().(*Entry)
}
func putEntry(e *Entry) {
	e.reset()
	_cePool.Put(e)
}

//日志堆栈信息
type EntryCaller struct {
	Defined  bool
	PC       uintptr
	File     string
	Line     int
	Function string
	Stack    string
}

func (ec EntryCaller) String() string {
	return ec.FullPath()
}
func (ec EntryCaller) FullPath() string {
	if !ec.Defined {
		return "undefined"
	}
	buf := pools.BufferPoolGet()
	buf.AppendString(ec.File)
	buf.AppendByte(':')
	buf.AppendInt(int64(ec.Line))
	caller := buf.String()
	buf.Free()
	return caller
}
func (ec EntryCaller) TrimmedPath() string {
	if !ec.Defined {
		return "undefined"
	}
	idx := strings.LastIndexByte(ec.File, '/')
	if idx == -1 {
		return ec.FullPath()
	}
	idx = strings.LastIndexByte(ec.File[:idx], '/')
	if idx == -1 {
		return ec.FullPath()
	}
	buf := pools.BufferPoolGet()
	buf.AppendString(ec.File[idx+1:])
	buf.AppendByte(':')
	buf.AppendInt(int64(ec.Line))
	caller := buf.String()
	buf.Free()
	return caller
}

type Entry struct {
	Name    string
	Level   Loglevel
	Caller  EntryCaller
	Time    time.Time
	Message string
	Data    Fields
	Err     string
}

func (entry *Entry) reset() {
	entry.Message = ""
	entry.Caller.Defined = false
	entry.Caller.File = ""
	entry.Caller.Line = 0
	entry.Caller.Function = ""
	entry.Caller.Stack = ""
	entry.Err = ""
	entry.Data = entry.Data[:0]
}

func (entry *Entry) WithFields(fields ...Field) {
	if cap(entry.Data) < len(fields) {
		entry.Data = make(Fields, 0, cap(entry.Data)+len(fields))
	}
	fieldErr := entry.Err
	for _, v := range fields {
		isErrField := false
		if t := reflect.TypeOf(v.Value); t != nil {
			switch {
			case t.Kind() == reflect.Func, t.Kind() == reflect.Ptr && t.Elem().Kind() == reflect.Func:
				isErrField = true
			}
		}
		if isErrField {
			tmp := fmt.Sprintf("can not add field %q", v.Key)
			if fieldErr != "" {
				fieldErr = entry.Err + ", " + tmp
			} else {
				fieldErr = tmp
			}
		} else {
			entry.Data = append(entry.Data, v)
		}
	}
}
