package registry

import (
	"time"
)

type FindServiceHandle func(snode ServiceNode)
type UpDataServiceHandle func(snode ServiceNode)
type LoseServiceHandle func(sId string)

type Option func(*Options)
type Options struct {
	Address          string
	Timeout          time.Duration
	RegisterInterval time.Duration
	RegisterTTL      time.Duration
	Find             FindServiceHandle
	UpData           UpDataServiceHandle
	Lose             LoseServiceHandle
	Tag              string //集群标签
}

func Address(v string) Option {
	return func(o *Options) {
		o.Address = v
	}
}

func RegisterInterval(v time.Duration) Option {
	return func(o *Options) {
		o.RegisterInterval = v
	}
}

func RegisterTTL(v time.Duration) Option {
	return func(o *Options) {
		o.RegisterTTL = v
	}
}

func FindHandle(v FindServiceHandle) Option {
	return func(o *Options) {
		o.Find = v
	}
}

func UpDataHandle(v UpDataServiceHandle) Option {
	return func(o *Options) {
		o.UpData = v
	}
}
func LoseHandle(v LoseServiceHandle) Option {
	return func(o *Options) {
		o.Lose = v
	}
}

func SetTag(v string) Option {
	return func(o *Options) {
		o.Tag = v
	}
}
func newOptions(opts ...Option) Options {
	opt := Options{
		RegisterInterval: time.Second * time.Duration(10),
		RegisterTTL:      time.Second * time.Duration(20),
		Tag:              "github.com/liwei1dao/lego",
	}
	for _, o := range opts {
		o(&opt)
	}
	return opt
}
