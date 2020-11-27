package cron

import (
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
	defsys Icron
)

func OnInit(config map[string]interface{}) (err error) {
	if defsys, err = newCron(newOptionsByConfig(config)); err == nil {
		Start()
	}
	return
}

func NewSys(option ...Option) (sys Icron, err error) {
	if sys, err = newCron(option...); err == nil {
		Start()
	}
	return
}

func Start() {
	defsys.Start()
}

func Stop() {
	defsys.Stop()
}

func AddFunc(spec string, cmd func()) (tcron.EntryID, error) {
	return defsys.AddFunc(spec, cmd)
}
