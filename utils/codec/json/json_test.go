package json_test

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/utils/codec/json"
)

type TestData struct {
	Name string
	Age  int
	List []string
	Map  map[string]interface{}
}

//测试api_getlist
func Test_Json_Write(t *testing.T) {
	ret, err := json.MarshalMap(&TestData{
		Name: "liwei",
		Age:  10,
		List: []string{"123", "456", "789"},
		Map: map[string]interface{}{
			"aa": 123,
			"b":  "123123",
		},
	})
	fmt.Printf("ret:%v  err:%v", ret, err)
	ret, err = json.MarshalMap(&TestData{
		Name: "asdasd",
		Age:  10,
		List: []string{"12asd3", "45sdaa6", "asdasd"},
		Map: map[string]interface{}{
			"asd":    586,
			"asdasd": "asd1231",
		},
	})
	fmt.Printf("ret:%v  err:%v", ret, err)
}
