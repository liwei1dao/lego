package codec_test

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/liwei1dao/lego/sys/codec"
	"github.com/liwei1dao/lego/sys/log"

	"github.com/modern-go/reflect2"
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
	if err := log.OnInit(nil); err != nil {
		fmt.Printf("log init err:%v", err)
		return
	}
	if sys, err := codec.NewSys(); err != nil {
		fmt.Printf("gin init err:%v", err)
	} else {
		data := []*Test1Data{{"liwe", 1}, {"liwe2", 2}}
		d, err := sys.MarshalJson(data)
		fmt.Printf("codec Marshal d:%s err:%v", d, err)
		data = []*Test1Data{}
		err = sys.UnmarshalJson(d, &data)
		fmt.Printf("codec UnmarshalJson data:%v err:%v", data, err)
	}
}
func Test_sys_json(t *testing.T) {
	if err := log.OnInit(nil); err != nil {
		fmt.Printf("log init err:%v", err)
		return
	}
	if sys, err := codec.NewSys(); err != nil {
		fmt.Printf("gin init err:%v", err)
	} else {
		d, err := sys.MarshalJson(&TestData{Name: "http://liwei1dao.com?asd=1&dd=1", Value: 10, Array: []interface{}{1, "dajiahao", &Test1Data{Name: "liwe1dao", Value: 123}}, Data: map[string]interface{}{"hah": 1, "asd": 999}})
		fmt.Printf("codec Marshal d:%s err:%v", d, err)
		data := &TestData{}
		err = sys.UnmarshalJson(d, data)
		fmt.Printf("codec UnmarshalJson data:%v err:%v", data, err)
	}
}

func Test_sys_mapjson(t *testing.T) {
	if err := log.OnInit(nil); err != nil {
		fmt.Printf("log init err:%v", err)
		return
	}
	if sys, err := codec.NewSys(); err != nil {
		fmt.Printf("gin init err:%v", err)
	} else {
		m := map[string]interface{}{"liwe": 123, "aasd": "123"}
		fmt.Printf("codec Marshal m:%s err:%v", m, err)
		d, err := sys.MarshalMapJson(&TestData{Name: "http://liwei1dao.com?asd=1&dd=1", Value: 10, Array: []interface{}{1, "dajiahao", &Test1Data{Name: "liwe1dao", Value: 123}}, Data: map[string]interface{}{"hah": 1, "asd": 999}})
		fmt.Printf("codec Marshal d:%s err:%v", d, err)
		data := &TestData{}
		err = sys.UnmarshalMapJson(d, data)
		fmt.Printf("codec UnmarshalJson data:%v err:%v", data, err)
	}
}

func Test_sys_reflect2(t *testing.T) {
	data := []*Test1Data{}
	ptr := reflect2.TypeOf(&data)
	kind := ptr.Kind()
	switch kind {
	case reflect.Interface:
		return
	case reflect.Struct:
		return
	case reflect.Array:
		return
	case reflect.Slice:
		return
	case reflect.Map:
		return
	case reflect.Ptr:
		ptrType := ptr.(*reflect2.UnsafePtrType)
		elemType := ptrType.Elem()
		kind = elemType.Kind()
		if kind == reflect.Slice {
			sliceelem := elemType.(*reflect2.UnsafeSliceType).Elem()
			sliceelemkind := sliceelem.Kind()
			if sliceelemkind == reflect.Ptr {
				return
			}
			return
		}
		return
	default:
		return
	}
}
