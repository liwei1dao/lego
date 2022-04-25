package core

type (
	IEngine interface {
		Timeout() int
		GopNum() int
		FLVDir() string
		FLVArchive() bool
		RTMPNoAuth() bool
		UseHlsHttps() bool
		HLSKeepAfterEnd() bool
		CheckAppName(appname string) bool
		GetStaticPushUrlList(appname string) ([]string, bool)
		GetKey(channel string) (newKey string, err error)
		GetChannel(key string) (channel string, err error)
	}
)
