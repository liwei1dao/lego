package email

type (
	ISys interface {
		SendMail(title, content string, to ...string) error
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

func SendMail(title, content string, to ...string) error {
	return defsys.SendMail(title, content, to...)
}
