package amf

import "io"

type Encoder struct {
}

type ExternalHandler func(*Decoder, io.Reader) (interface{}, error)
type Decoder struct {
	refCache         []interface{}
	stringRefs       []string
	objectRefs       []interface{}
	traitRefs        []Trait
	externalHandlers map[string]ExternalHandler
}

type Trait struct {
	Type           string
	Externalizable bool
	Dynamic        bool
	Properties     []string
}
