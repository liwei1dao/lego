package codec_test

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/utils/codec/json"
)

type TestData struct {
	Name  string
	Value int
	Array []interface{}
	Data  map[string]interface{}
}
type Test1Data struct {
	Name  string
	Value int
}

func Test_sys_slice(t *testing.T) {
	data := []*Test1Data{{"liwe", 1}, {"liwe2", 2}}
	d, err := json.Marshal(data)
	fmt.Printf("codec Marshal d:%s err:%v", d, err)
	data = []*Test1Data{}
	err = json.Unmarshal(d, &data)
	fmt.Printf("codec UnmarshalJson data:%v err:%v", data, err)
}
func Test_sys_json(t *testing.T) {
	d, err := json.Marshal(&TestData{Name: "http://liwei1dao.com?asd=1&dd=1", Value: 10, Array: []interface{}{1, "dajiahao", &Test1Data{Name: "liwe1dao", Value: 123}}, Data: map[string]interface{}{"hah": 1, "asd": 999}})
	fmt.Printf("codec Marshal d:%s err:%v", d, err)
	data := &TestData{}
	err = json.Unmarshal(d, data)
	fmt.Printf("codec UnmarshalJson data:%v err:%v", data, err)
}

func Test_sys_mapjson(t *testing.T) {
	d, err := json.MarshalMap(&TestData{Name: "http://liwei1dao.com?asd=1&dd=1", Value: 10, Array: []interface{}{1, "dajiahao", &Test1Data{Name: "liwe1dao", Value: 123}}, Data: map[string]interface{}{"hah": 1, "asd": 999}})
	fmt.Printf("codec Marshal d:%s err:%v", d, err)
	data := &TestData{}
	err = json.UnmarshalMap(d, data)
	fmt.Printf("codec UnmarshalJson data:%v err:%v", data, err)

}
