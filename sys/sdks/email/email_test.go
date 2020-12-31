package email

import (
	"testing"
)

func Test_email(t *testing.T) {
	sys, _ := NewSys(SetServerhost("smtp.qq.com"), SetFromemail("liwei1dao@foxmail.com"), SetFompasswd("tucfnwfgmnqleaaf"), SetServerport(587))
	sys.SendMail("2323654976@qq.com", "验证码", "52678")
}
