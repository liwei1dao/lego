package cron

import (
	tcron "github.com/robfig/cron/v3"
)

func newSys(options Options) (sys *Cron, err error) {
	sys = &Cron{options: options, cron: tcron.New()}
	return
}

type Cron struct {
	cron    *tcron.Cron
	options Options
}

func (this *Cron) Start() {
	this.cron.Start()
}

func (this *Cron) Stop() {
	this.cron.Stop()
}

func (this *Cron) AddFunc(spec string, cmd func()) (tcron.EntryID, error) {
	return this.cron.AddFunc(spec, cmd)
}
