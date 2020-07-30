package email

import (
	"testing"
)

func Test_email(t *testing.T) {
	sdk, _ := newEmail(SetServerhost("smtp.qq.com"), SetFromemail("liwei1dao@foxmail.com"), SetFompasswd("tucfnwfgmnqleaaf"), SetServerport(587))
	sdk.SendCaptcha("2323654976@qq.com", "52678")
}
