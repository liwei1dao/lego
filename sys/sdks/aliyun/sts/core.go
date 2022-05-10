package sts

type (
	Authorization struct {
		Expiration      string
		AccessKeyId     string
		AccessKeySecret string
		SecurityToken   string
	}
	ISys interface {
		//roleArn:角色ARN。
		AssumeRole(roleArn, roleSessionName string) (auth *Authorization, err error)
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

func AssumeRole(roleArn, roleSessionName string) (auth *Authorization, err error) {
	return defsys.AssumeRole(roleArn, roleSessionName)
}
