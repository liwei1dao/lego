package codec

import (
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
