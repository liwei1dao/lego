package cron

import (
	tcron "github.com/robfig/cron/v3"
)

type (
	EntryID tcron.EntryID
	ISys    interface {
		Start()
		Close()
		AddFunc(spec string, cmd func()) (EntryID, error)
		Remove(id EntryID)
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

func Close() {
	defsys.Close()
}

func AddFunc(spec string, cmd func()) (EntryID, error) {
	return defsys.AddFunc(spec, cmd)
}

func Remove(id EntryID) {
	defsys.Remove(id)
}
