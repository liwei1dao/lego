package http

import (
	"fmt"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	IHttpOptions interface {
		GetHttpAddr() string
		GettCertPath() string
		GetKeyPath() string
	}
	HttpOptions struct {
		HttpAddr string
		CertPath string
		KeyPath  string
	}
)

func (this *HttpOptions) LoadConfig(settings map[string]interface{}) (err error) {
	if settings != nil {
		mapstructure.Decode(settings, &this)
	}
	if this.HttpAddr == "" {
		err = fmt.Errorf("module http missing necessary configuration:%v", settings)
	}
	return
}

func (this *HttpOptions) GetHttpAddr() string {
	return this.HttpAddr
}
func (this *HttpOptions) GettCertPath() string {
	return this.CertPath
}
func (this *HttpOptions) GetKeyPath() string {
	return this.KeyPath
}
