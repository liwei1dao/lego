package ssynth_test

import (
	"context"
	"fmt"
	"os"
	"testing"

	"github.com/liwei1dao/lego/sys/sdks/xunfei/ssynth"
)

func Test_sys(t *testing.T) {
	var (
		sys ssynth.ISys
		err error
	)
	if sys, err = ssynth.NewSys(
		ssynth.SetHostUrl("wss://tts-api.xfyun.cn/v2/tts"),
		ssynth.SetAppid("xxxxxx"),
		ssynth.SetApiKey("xxxxxxxxxxxxxxxxx"),
		ssynth.SetApiSecret("xxxxxxxxxxxxxxxxx"),
	); err != nil {
		fmt.Printf("livego init err:%v \n", err)
		return
	}
	audioFile, err := os.OpenFile("./test.mp3", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		fmt.Printf("openFile err:%v \n", err)
		return
	}
	defer audioFile.Close()
	if err = sys.TxtToVoice(context.Background(), "你好呀,我的亲,你在干嘛呢？", audioFile); err != nil {
		fmt.Printf("TxtToVoice err:%v \n", err)
		return
	}
}
