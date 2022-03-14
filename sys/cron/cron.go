package cron

import (
	tcron "github.com/robfig/cron/v3"
)

func newSys(options Options) (sys *Cron, err error) {
	parser := tcron.NewParser(tcron.Second | tcron.Minute | tcron.Hour | tcron.Dom | tcron.Month | tcron.Dow | tcron.Descriptor)
	sys = &Cron{options: options, cron: tcron.New(tcron.WithParser(parser))}
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

func (this *Cron) Remove(id tcron.EntryID) {
	this.cron.Remove(id)
}
