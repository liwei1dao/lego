package httppool

import (
	"net/http"
	"sync"
	"time"
)

// NewClientPool 返回一个新的ClientPool实例
func newSys(options Options) (sys *HttpPool, err error) {
	sys = &HttpPool{
		options: &options,
	}
	sys.once.Do(func() {
		transport := &http.Transport{
			MaxIdleConns:        options.MaxIdleConns,
			MaxIdleConnsPerHost: options.MaxIdleConnsPerHost,
			IdleConnTimeout:     time.Second * time.Duration(options.IdleConnTimeout),
		}

		sys.client = &http.Client{
			Transport: transport,
			Timeout:   time.Second * time.Duration(options.IdleConnTimeout),
		}
	})
	return
}

type HttpPool struct {
	options *Options
	client  *http.Client
	once    sync.Once
}

// 获取实例对象
func (this *HttpPool) GetClient() *http.Client {
	return this.client
}
