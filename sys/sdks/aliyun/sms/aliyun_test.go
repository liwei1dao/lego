package sms

import (
	"testing"
)

func Test_aliyun(t *testing.T) {
	sdk, _ := NewSys(Setkey("你的阿里云key"), KeySecret("你的阿里云keysecret"), SetSignName("你的SignName"))
	sdk.SendCaptcha("13076983350", "52679")
}
