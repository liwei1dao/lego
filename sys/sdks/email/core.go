package email

import (
	"github.com/liwei1dao/lego/core"
)

type (
	IEmail interface {
		SendCaptcha(email string, captcha string) error
	}
)

var (
	service core.IService
	email   IEmail
)

func OnInit(s core.IService, opt ...Option) (err error) {
	service = s
	email, err = newEmail(opt...)
	return
}

func SendCaptcha(_email string, _captcha string) error {
	return email.SendCaptcha(_email, _captcha)
}
