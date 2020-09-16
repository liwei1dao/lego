package token

import (
	"github.com/liwei1dao/lego/core"
)

type (
	IToken interface {
		Start() (err error)
		Stop() (err error)
	}
)

var (
	tokne IToken
)

func OnInit(s core.IService, opt ...Option) (err error) {
	tokne, err = newEthPay(opt...)
	return
}

func Start() (err error) {
	return tokne.Start()
}

func Stop() (err error) {
	return tokne.Stop()
}
