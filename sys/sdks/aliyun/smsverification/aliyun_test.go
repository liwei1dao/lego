package smsverification

import (
	"testing"
)

func Test_aliyun(t *testing.T) {
	sdk := NewALiYun("你的阿里云key", "你的阿里云keysecret", "你的SignName")
	sdk.SendCaptcha("13076983350", "52679")
}
