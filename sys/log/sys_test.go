package log_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/liwei1dao/lego/sys/log"
)

type TestData struct {
	Name string
	Age  int32
}

func (this *TestData) Log() {
	sys.Error("妈妈咪呀!")
}

var sys log.ISys

func TestMain(m *testing.M) {
	var err error
	if sys, err = log.NewSys(
		log.SetFileName("log.log"),
		log.SetIsDebug(true),
		log.SetEncoder(log.TextEncoder),
	); err != nil {
		fmt.Println(err)
		return
	}
	defer os.Exit(m.Run())
}
func Test_sys(t *testing.T) {
	data := &TestData{}
	data.Log()
}

//性能测试
func Benchmark_Ability(b *testing.B) {
	for i := 0; i < b.N; i++ { //use b.N for looping
		sys.Error("妈妈咪呀!")
	}
}
