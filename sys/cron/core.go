package cron

import (
	tcron "github.com/robfig/cron/v3"
)

type (
	ISys interface {
		Start()
		Stop()
		AddFunc(spec string, cmd func()) (tcron.EntryID, error)
		Remove(id tcron.EntryID)
	}
)

var (
	defsys ISys
)

func OnInit(config map[string]interface{}, option ...Option) (err error) {
	if defsys, err = newSys(newOptions(config, option...)); err == nil {
		Start()
	}
	return
}

func NewSys(option ...Option) (sys ISys, err error) {
	if sys, err = newSys(newOptionsByOption(option...)); err == nil {
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

func Remove(id tcron.EntryID) {
	defsys.Remove(id)
}
