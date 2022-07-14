package core

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
