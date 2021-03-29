package live

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	IOptions interface {
		GetRtmpAddr() string
		GetRtmpNoAuth() bool
		CheckAppName(appname string) bool
	}
	Application struct {
		Appname    string   `mapstructure:"appname"`
		Live       bool     `mapstructure:"live"`
		Hls        bool     `mapstructure:"hls"`
		Flv        bool     `mapstructure:"flv"`
		Api        bool     `mapstructure:"api"`
		StaticPush []string `mapstructure:"static_push"`
	}
	Applications []Application
	Options      struct {
		RtmpAddr   string
		RtmpNoAuth bool
		Server     Applications
	}
)

func (this *Options) LoadConfig(settings map[string]interface{}) (err error) {
	if settings != nil {
		if err = mapstructure.Decode(settings, this); err != nil {
			return
		}
	}
	return
}

func (this *Options) GetRtmpAddr() string {
	return this.RtmpAddr
}

func (this *Options) GetRtmpNoAuth() bool {
	return this.RtmpNoAuth
}

func (this *Options) CheckAppName(appname string) bool {
	for _, app := range this.Server {
		if app.Appname == appname {
			return app.Live
		}
	}
	return false
}
