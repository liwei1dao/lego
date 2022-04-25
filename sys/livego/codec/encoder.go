package codec

import (
	"fmt"
	"io"
)

type Encoder struct {
}

func (e *Encoder) Encode(w io.Writer, val interface{}, ver Version) (int, error) {
	switch ver {
	case AMF0:
		return e.EncodeAmf0(w, val)
	case AMF3:
		return e.EncodeAmf3(w, val)
	}

	return 0, fmt.Errorf("encode amf: unsupported version %d", ver)
}
