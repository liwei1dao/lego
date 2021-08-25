package http

import (
	"fmt"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	IOptions interface {
		GetListenProt() int
		GettCertPath() string
		GetKeyPath() string
	}
	Options struct {
		ListenProt int    `json:"-"`
		CertPath   string `json:"-"`
		KeyPath    string `json:"-"`
	}
)

func (this *Options) LoadConfig(settings map[string]interface{}) (err error) {
	if settings != nil {
		if err = mapstructure.Decode(settings, this); err != nil {
			return
		}
	}
	if this.ListenProt == 0 {
		err = fmt.Errorf("module http missing necessary configuration:%v", settings)
	}
	return
}

func (this *Options) GetListenProt() int {
	return this.ListenProt
}
func (this *Options) GettCertPath() string {
	return this.CertPath
}
func (this *Options) GetKeyPath() string {
	return this.KeyPath
}
