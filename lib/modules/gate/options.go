package gate

import (
	"time"

	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	IOptions interface {
		GetTcpAddr() string
		GetWSAddr() string
		GetWSOuttime() time.Duration
		GetWSHeartbeat() time.Duration
		GetWSCertFile() string
		GetWSKeyFile() string
	}
	Options struct {
		TcpAddr     string
		WSAddr      string
		WSOuttime   int
		WSHeartbeat int
		WSCertFile  string
		WSKeyFile   string
	}
)

func (this *Options) GetTcpAddr() string {
	return this.TcpAddr
}
func (this *Options) GetWSAddr() string {
	return this.WSAddr
}
func (this *Options) GetWSOuttime() time.Duration {
	return time.Duration(this.WSOuttime) * time.Second
}
func (this *Options) GetWSHeartbeat() time.Duration {
	return time.Duration(this.WSHeartbeat) * time.Second
}
func (this *Options) GetWSCertFile() string {
	return this.WSCertFile
}
func (this *Options) GetWSKeyFile() string {
	return this.WSKeyFile
}

func (this *Options) LoadConfig(settings map[string]interface{}) (err error) {
	if settings != nil {
		err = mapstructure.Decode(settings, this)
	}
	return
}
