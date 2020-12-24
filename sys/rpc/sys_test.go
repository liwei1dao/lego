package rpc

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/golang/protobuf/proto"
	"github.com/liwei1dao/lego/sys/rpc/core"
)

func Test_sys(t *testing.T) {
	data := &core.ResultInfo{}
	dtype := reflect.TypeOf(data)
	fmt.Printf("data type:%v \n", dtype.Kind())
	PrintfType(data)
}

func PrintfType(i interface{}) {
	switch i.(type) {
	case proto.Message:
		fmt.Printf("data type:proto.Message \n")
		break
	}
}
