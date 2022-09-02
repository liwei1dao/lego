package translation

import "context"

/*
	Google 翻译系统
*/

type (
	ISys interface {
		Translation_Text_Base(ctx context.Context, original string, from, to string) (result string, err error)
		Translation_Text_GoogleV3(ctx context.Context, original string, from, to string) (result string, err error)
		Translation_Voice(ctx context.Context, original []byte, from string) (result string, err error)
	}
)

var defsys ISys

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	defsys, err = newSys(newOptions(config, option...))
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	sys, err = newSys(newOptionsByOption(option...))
	return
}

func Translation_Text_Base(ctx context.Context, original string, from, to string) (result string, err error) {
	return defsys.Translation_Text_Base(ctx, original, from, to)
}

func Translation_Text_GoogleV3(ctx context.Context, original string, from, to string) (result string, err error) {
	return defsys.Translation_Text_GoogleV3(ctx, original, from, to)
}

func Translation_Voice(ctx context.Context, original []byte, from string) (result string, err error) {
	return defsys.Translation_Voice(ctx, original, from)
}
