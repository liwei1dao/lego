package sms_test

import (
	"testing"

	"github.com/liwei1dao/lego/sys/sdks/aliyun/sms"
)

func Test_aliyun(t *testing.T) {
	sdk, _ := sms.NewSys(sms.Setkey("你的阿里云key"), sms.KeySecret("你的阿里云keysecret"), sms.SetSignName("你的SignName"))
	sdk.SendCaptcha("13076983350", "52679")
}
