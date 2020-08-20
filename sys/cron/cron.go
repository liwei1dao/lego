package cron

import (
	tcron "github.com/robfig/cron/v3"
)

func newCron(opt ...Option) (*Cron, error) {
	cron := new(Cron)
	cron.opts = newOptions(opt...)
	err := cron.init()
	return cron, err
}

type Cron struct {
	cron *tcron.Cron
	opts Options
}

func (this *Cron) init() (err error) {
	this.cron = tcron.New()
	return
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
