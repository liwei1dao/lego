package log_test

import (
	"fmt"
	"testing"

	"github.com/liwei1dao/lego/sys/log"
)

type TestData struct {
	Name string
	Age  int32
}

func Test_sys(t *testing.T) {
	if sys, err := log.NewSys(log.SetFileName("log.log"), log.SetEncoder(log.TextEncoder)); err != nil {
		fmt.Println(err)
		return
	} else {
		sys.Debug("妈妈咪呀!", log.Field{"num", &TestData{Name: "lala", Age: 165}})
	}
}
