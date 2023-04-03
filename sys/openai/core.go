package openai

type (
	ISys interface {
		SendReq(content string) (results string, err error)
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

func SendReq(content string) (results string, err error) {
	return defsys.SendReq(content)
}
