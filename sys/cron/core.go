package cron

import (
	"github.com/liwei1dao/lego/core"
	tcron "github.com/robfig/cron/v3"
)

type (
	Icron interface {
		Start()
		Stop()
		AddFunc(spec string, cmd func()) (tcron.EntryID, error)
	}
)

var (
	defcron Icron
)

func OnInit(s core.IService, opt ...Option) (err error) {
	if defcron, err = newCron(opt...); err == nil {
		Start()
	}
	return
}

func Start() {
	defcron.Start()
}

func Stop() {
	defcron.Stop()
}

func AddFunc(spec string, cmd func()) (tcron.EntryID, error) {
	return defcron.AddFunc(spec, cmd)
}
