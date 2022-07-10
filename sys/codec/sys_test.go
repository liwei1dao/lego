package codec_test

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/sys/codec"
	"github.com/liwei1dao/lego/sys/log"
)

type TestData struct {
	Name  string
	Value int
}

func Test_sys(t *testing.T) {
	if err := log.OnInit(nil); err != nil {
		fmt.Printf("log init err:%v", err)
		return
	}
	if sys, err := codec.NewSys(); err != nil {
		fmt.Printf("gin init err:%v", err)
	} else {
		d, err := sys.MarshalJson(&TestData{Name: "liwe1idao", Value: 10})
		fmt.Printf("codec Marshal d:%s err:%v", d, err)
	}
}
