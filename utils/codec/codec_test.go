package codec

import (
	"encoding/json"
	"fmt"
	"testing"
)

type TestData struct {
	Fild_1 string
	Fild_3 int
	Fild_4 float32
}

func Test_Encoder(t *testing.T) {
	encoder := &Encoder{}
	// data, err := encoder.EncoderToMap(map[string]interface{}{"liwei": 106, "sasd": "2564"})
	// fmt.Printf("EncoderToMap data1:%v err:%v", data, err)
	data, err := encoder.EncoderToMap([]interface{}{"liwei", 106, "sasd", "2564"})
	fmt.Printf("EncoderToMap data1:%v err:%v", data, err)
	// data, err := encoder.EncoderToMap(&TestData{Fild_1: "liwei1dao", Fild_3: 25, Fild_4: 3.54})
	// fmt.Printf("EncoderToMap data2:%v err:%v", data, err)
}

func Test_Decoder(t *testing.T) {
	decoder := &Decoder{}
	data := &TestData{}
	err := decoder.DecoderMapString(map[string]string{"Fild_1": "liwei1dao"}, data)
	fmt.Printf("DecoderMap data1:%v err:%v", data, err)
}

func Test_Slice_Decoder(t *testing.T) {
	decoder := &Decoder{DefDecoder: json.Unmarshal}
	encoder := &Encoder{DefEncoder: json.Marshal}
	data := []*TestData{{Fild_1: "1dao", Fild_3: 10, Fild_4: 3.5}, {Fild_1: "2dao", Fild_3: 20, Fild_4: 6.5}}
	datastr, err := encoder.EncoderToSliceString(data)
	fmt.Printf("EncoderToSliceString datastr:%v err:%v", datastr, err)
	if err != nil {
		return
	}
	data1 := make([]*TestData, 0)
	err = decoder.DecoderSliceString(datastr, data1)
	fmt.Printf("DecoderMap data1:%v err:%v", data1, err)
}

func Test_Slice_Type(t *testing.T) {
	decoder := &Decoder{DefDecoder: json.Unmarshal}
	encoder := &Encoder{DefEncoder: json.Marshal}
	data := []*TestData{{Fild_1: "1dao", Fild_3: 10, Fild_4: 3.5}, {Fild_1: "2dao", Fild_3: 20, Fild_4: 6.5}}
	datastr, err := encoder.EncoderToSliceString(data)
	fmt.Printf("EncoderToSliceString datastr:%v err:%v", datastr, err)
	if err != nil {
		return
	}
	data1 := make([]*TestData, 0)
	err = decoder.DecoderSliceString(datastr, &data1)
	fmt.Printf("DecoderMap data1:%v err:%v", data1, err)
}

func Test_EncoderToMapString(t *testing.T) {
	encoder := &Encoder{DefEncoder: json.Marshal}
	data := &TestData{Fild_1: "1dao", Fild_3: 10, Fild_4: 3.5}
	_map, err := encoder.EncoderToMapString(data)
	fmt.Printf("DecoderMap map:%v err:%v", _map, err)
}
