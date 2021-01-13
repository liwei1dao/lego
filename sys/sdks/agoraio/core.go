package agoraio

/*
声网SDK https://www.agora.io/cn
*/

type (
	IAgoraio interface {
		CreateToken(uid uint32, channelName string) (token string, err error)
	}
)

var defsys IAgoraio

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	if defsys, err = newSys(newOptions(config, option...)); err == nil {

	}
	return
}

func NewSys(option ...Option) (sys IAgoraio, err error) {
	if sys, err = newSys(newOptionsByOption(option...)); err == nil {

	}
	return
}

func CreateToken(uid uint32, channelName string) (token string, err error) {
	return defsys.CreateToken(uid, channelName)
}
