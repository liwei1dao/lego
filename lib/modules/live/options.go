package live

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	IOptions interface {
		GetRtmpAddr() string
		GetRtmpNoAuth() bool
		GetCacheAddr() string
		GetFLVArchive() bool
		GetFLVDir() string
		GetReadTimeout() int
		GetWriteTimeout() int
		GetGopNum() int
		GetHls() bool
		GetFlv() bool
		GetApi() bool
		CheckAppName(appname string) bool
		GetStaticPushUrlList(appname string) ([]string, bool)
	}
	Options struct {
		RtmpAddr     string
		RtmpNoAuth   bool
		CacheAddr    string
		FLVArchive   bool
		FLVDir       string
		ReadTimeout  int
		WriteTimeout int
		GopNum       int
		Appname      string
		Live         bool
		Hls          bool
		Flv          bool
		Api          bool
		StaticPush   []string
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
func (this *Options) GetCacheAddr() string {
	return this.CacheAddr
}
func (this *Options) GetRtmpNoAuth() bool {
	return this.RtmpNoAuth
}
func (this *Options) GetFLVArchive() bool {
	return this.FLVArchive
}
func (this *Options) GetFLVDir() string {
	return this.FLVDir
}

func (this *Options) GetReadTimeout() int {
	return this.ReadTimeout
}
func (this *Options) GetWriteTimeout() int {
	return this.WriteTimeout
}

func (this *Options) CheckAppName(appname string) bool {
	if this.Appname == appname {
		return this.Live
	}
	return false
}

func (this *Options) GetStaticPushUrlList(appname string) ([]string, bool) {
	if (this.Appname == appname) && this.Live {
		if len(this.StaticPush) > 0 {
			return this.StaticPush, true
		} else {
			return nil, false
		}
	}
	return nil, false
}
