package codec_test

import (
	"fmt"
	"reflect"
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

type test struct {
	Name string
	F1   *field
}

type field struct {
	Name string
}

func Test_reflect(t *testing.T) {
	test1 := &test{Name: "test1"}
	test2 := &test{Name: "test2", F1: &field{Name: "field2"}}

	v1 := reflect.ValueOf(test1)
	v2 := reflect.ValueOf(test2).Elem()
	v1.Elem().Set(v2)
	fmt.Printf("v1:%v", test1)
	f3 := &field{Name: "field3"}
	v3 := reflect.ValueOf(f3).Elem()
	v1.Elem().FieldByName("F1").Elem().Set(v3)
	fmt.Printf("F1:%v", test1.F1)
}
