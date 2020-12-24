package rpc

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
)

func Test_sys(t *testing.T) {
	// data := &core.ResultInfo{}
	dtype := reflect.TypeOf(PrintfFunc)
	// fmt.Printf("data Align:%v \n", dtype.Align())
	// fmt.Printf("data FieldAlign:%v \n", dtype.FieldAlign())
	// fmt.Printf("data NumMethod:%v \n", dtype.NumMethod())
	// fmt.Printf("data Name:%v \n", dtype.Name())
	// fmt.Printf("data Kind:%v \n", dtype.Kind())
	// fmt.Printf("data Comparable:%v \n", dtype.Comparable())
	// // fmt.Printf("data Bits:%v \n", dtype.Bits())
	// // fmt.Printf("data ChanDir:%+v \n", dtype.ChanDir())
	// fmt.Printf("data IsVariadic:%v \n", dtype.IsVariadic())
	// // fmt.Printf("data Elem:%v \n", dtype.Elem().String())
	// // fmt.Printf("data Key:%v \n", dtype.Key())
	// fmt.Printf("data Len:%v \n", dtype.Len())
	// fmt.Printf("data NumField:%v \n", dtype.NumField())
	// fmt.Printf("data NumIn:%v \n", dtype.NumIn())
	// fmt.Printf("data NumOut:%v \n", dtype.NumOut())
	// PrintfType(data)
	numIn := dtype.NumIn()
	addIn := make([]reflect.Type, numIn)
	for i := 0; i < numIn; i++ {
		addIn[i] = dtype.In(i)
		fmt.Printf("func In:%d type:%v \n", i, addIn[i].String())
	}
}

func PrintfType(i interface{}) {
	switch i.(type) {
	case proto.Message:
		fmt.Printf("data type:proto.Message \n")
		break
	}
}

func PrintfFunc(a int, b bool) {

}
