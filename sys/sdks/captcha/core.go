package sdk_aliyun

import (
	"fmt"
	"github.com/liwei1dao/lego/core"
)

type (
	CaptchaService interface {
		SendCaptcha(target, captcha string) error //发送验证码
	}
)

var opt *Options

func OnInit(s core.IService, opts ...Option) (err error) {
	opt = newOptions(opts...)
	return
}

func SendModuleCaptcha(module, captcha string) error {
	if opt != nil && opt.Aliyun != nil {
		return opt.Aliyun.SendCaptcha(module, captcha)
	} else {
		return fmt.Errorf("Sdk 未初始化")
	}
}

func SendEMAILCaptcha(email, captcha string) error {
	if opt != nil && opt.Aliyun != nil {
		return opt.Email.SendCaptcha(email, captcha)
	} else {
		return fmt.Errorf("Sdk 未初始化")
	}
}
