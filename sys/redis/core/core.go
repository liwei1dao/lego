package core

import (
	"encoding"
	"time"
)

type (
	ICodec interface {
		Marshal(v interface{}) ([]byte, error)
		Unmarshal(data []byte, v interface{}) error
		MarshalMap(val interface{}) (ret map[string]string, err error)
		UnmarshalMap(data map[string]string, val interface{}) (err error)
		MarshalSlice(val interface{}) (ret []string, err error)
		UnmarshalSlice(data []string, val interface{}) (err error)
	}
)

func IsBaseType(v interface{}) bool {
	switch v.(type) {
	case nil, string, []byte, int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64, bool, time.Time, time.Duration, encoding.BinaryMarshaler:
		return true
	default:
		return false
	}
}
