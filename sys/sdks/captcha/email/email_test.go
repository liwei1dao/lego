package sdk_email

import (
	"testing"
)

func Test_email(t *testing.T) {
	sdk := NewEmail("smtp.qq.com", "liwei1dao@foxmail.com", "tucfnwfgmnqleaaf", 587)
	sdk.SendCaptcha("2323654976@qq.com", "52678")
}
