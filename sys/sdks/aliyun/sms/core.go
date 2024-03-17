package sms

type (
	ISys interface {
		SendCaptcha(mobile string, captcha string) error
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func SendCaptcha(_email string, _captcha string) error {
	return defsys.SendCaptcha(_email, _captcha)
}
