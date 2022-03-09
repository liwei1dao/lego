package email_test

import (
	"testing"

	"github.com/liwei1dao/lego/sys/sdks/email"
)

func Test_email(t *testing.T) {
	sys, _ := email.NewSys(email.SetServerhost("smtp.qq.com"), email.SetFromemail("liwei1dao@foxmail.com"), email.SetFompasswd("tucfnwfgmnqleaaf"), email.SetServerport(587))
	sys.SendMail("2323654976@qq.com", "验证码", "52678")
}
