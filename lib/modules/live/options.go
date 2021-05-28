package live

import (
	"github.com/liwei1dao/lego/utils/mapstructure"
)

type (
	IOptions interface {
		GetRtmpAddr() string
		GetAPIAddr() string
		GetHlsAddr() string
		GetHttpFlvAddr() string
		GetHlsKeepAfterEnd() bool
		GetRtmpNoAuth() bool
		GetCacheAddr() string
		GetCacheDB() int
		GetFLVArchive() bool
		GetFLVDir() string
		GetReadTimeout() int
		GetWriteTimeout() int
		GetGopNum() int
		GetJWTAlgorithm() string
		GetJWTJWTSecret() string
		GetHls() bool
		GetFlv() bool
		GetApi() bool
		CheckAppName(appname string) bool
		GetStaticPushUrlList(appname string) ([]string, bool)
	}
	Options struct {
		RtmpAddr        string
		APIAddr         string
		HlsAddr         string
		HttpFlvAddr     string
		HlsKeepAfterEnd bool
		RtmpNoAuth      bool
		CacheAddr       string
		CacheDB         int
		FLVArchive      bool
		FLVDir          string
		ReadTimeout     int
		WriteTimeout    int
		GopNum          int
		JWTAlgorithm    string
		JWTSecret       string
		Appname         string
		Live            bool
		Hls             bool
		Flv             bool
		Api             bool
		StaticPush      []string
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

func (this *Options) GetHlsAddr() string {
	return this.HlsAddr
}

func (this *Options) GetHttpFlvAddr() string {
	return this.HttpFlvAddr
}

func (this *Options) GetRtmpAddr() string {
	return this.RtmpAddr
}

func (this *Options) GetAPIAddr() string {
	return this.APIAddr
}

func (this *Options) GetHlsKeepAfterEnd() bool {
	return this.HlsKeepAfterEnd
}

func (this *Options) GetCacheAddr() string {
	return this.CacheAddr
}
func (this *Options) GetCacheDB() int {
	return this.CacheDB
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
func (this *Options) GetGopNum() int {
	return this.GopNum
}

func (this *Options) GetJWTAlgorithm() string {
	return this.JWTAlgorithm
}
func (this *Options) GetJWTJWTSecret() string {
	return this.JWTSecret
}

func (this *Options) GetHls() bool {
	return this.Hls
}
func (this *Options) GetFlv() bool {
	return this.Flv
}

func (this *Options) GetApi() bool {
	return this.Api
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
