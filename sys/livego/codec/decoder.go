package codec

import (
	"fmt"
	"io"
)

func NewTrait() *Trait {
	return &Trait{}
}

type Trait struct {
	Type           string
	Externalizable bool
	Dynamic        bool
	Properties     []string
}

type ExternalHandler func(*Decoder, io.Reader) (interface{}, error)

type Decoder struct {
	refCache         []interface{}
	stringRefs       []string
	objectRefs       []interface{}
	traitRefs        []Trait
	externalHandlers map[string]ExternalHandler
}

func (d *Decoder) DecodeBatch(r io.Reader, ver Version) (ret []interface{}, err error) {
	var v interface{}
	for {
		v, err = d.Decode(r, ver)
		if err != nil {
			break
		}
		ret = append(ret, v)
	}
	return
}

func (d *Decoder) Decode(r io.Reader, ver Version) (interface{}, error) {
	switch ver {
	case 0:
		return d.DecodeAmf0(r)
	case 3:
		return d.DecodeAmf3(r)
	}

	return nil, fmt.Errorf("decode amf: unsupported version %d", ver)
}
