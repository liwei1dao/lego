package sreco_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/liwei1dao/lego/sys/sdks/xunfei/sreco"
)

func Test_sys(t *testing.T) {
	var (
		sys    sreco.ISys
		result string
		err    error
	)
	if sys, err = sreco.NewSys(
		sreco.SetHostUrl("wss://iat-api.xfyun.cn/v2/iat"),
		sreco.SetAppid("xxxx"),
		sreco.SetApiKey("xxxxxxxxxxxxxxxxx"),
		sreco.SetApiSecret("xxxxxxxxxxxxxxxxxxxxxxxxxxx"),
	); err != nil {
		fmt.Printf("livego init err:%v \n", err)
		return
	}
	audioFile, err := os.Open("./16k_10.pcm")
	if err != nil {
		panic(err)
	}
	defer audioFile.Close()
	if result, err = sys.VoiceToTxt(context.Background(), audioFile); err != nil {
		fmt.Printf("VoiceToTxt err:%v \n", err)
		return
	} else {
		fmt.Printf("VoiceToTxt result:%s \n", result)
	}
}
