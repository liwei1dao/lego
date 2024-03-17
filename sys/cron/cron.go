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

func (this *Cron) Close() {
	this.cron.Stop()
}

func (this *Cron) AddFunc(spec string, cmd func()) (id EntryID, err error) {
	var eid tcron.EntryID
	if eid, err = this.cron.AddFunc(spec, cmd); err != nil {
		id = EntryID(eid)
	}
	return
}

func (this *Cron) Remove(id EntryID) {
	this.cron.Remove(tcron.EntryID(id))
}
