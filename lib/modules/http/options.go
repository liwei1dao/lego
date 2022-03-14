package http

import (
	"fmt"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	IOptions interface {
		GetListenPort() int
		GettCertPath() string
		GetKeyPath() string
		GetCors() bool
	}
	Options struct {
		ListenPort int    `json:"-"`
		CertPath   string `json:"-"`
		KeyPath    string `json:"-"`
		Cors       bool   `json:"-"` //是否跨域
	}
)

func (this *Options) LoadConfig(settings map[string]interface{}) (err error) {
	if settings != nil {
		if err = mapstructure.Decode(settings, this); err != nil {
			return
		}
	}
	if this.ListenPort == 0 {
		err = fmt.Errorf("module http missing necessary configuration:%v", settings)
	}
	return
}

func (this *Options) GetListenPort() int {
	return this.ListenPort
}
func (this *Options) GettCertPath() string {
	return this.CertPath
}
func (this *Options) GetKeyPath() string {
	return this.KeyPath
}
func (this *Options) GetCors() bool {
	return this.Cors
}
