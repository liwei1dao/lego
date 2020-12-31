package email

type (
	IEmail interface {
		SendEmail(temail string, title, content string) error
	}
)

var (
	defsys IEmail
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys IEmail, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func SendEmail(temail string, title, content string) error {
	return defsys.SendEmail(temail, title, content)
}
