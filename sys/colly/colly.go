package colly

import (
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/proxy"
)

func newSys(options Options) (sys *Colly, err error) {
	var (
		py colly.ProxyFunc
	)
	sys = &Colly{
		options: options,
		colly: colly.NewCollector(
			colly.AllowURLRevisit(),
		),
	}
	//限速
	// sys.colly.Limit(&colly.LimitRule{
	// 	DomainGlob:  "www.douban.com",
	// 	Parallelism: 1,
	// 	Delay:       2 * time.Second,
	// })
	if len(options.UserAgent) > 0 {
		sys.colly.UserAgent = options.UserAgent
	}
	if options.Proxy != nil && len(options.Proxy) > 0 {
		if py, err = proxy.RoundRobinProxySwitcher(options.Proxy...); err == nil {
			sys.colly.SetProxyFunc(py)
		}
	}
	return
}

type Colly struct {
	options Options
	colly   *colly.Collector
}

func (this *Colly) OnResponse(f colly.ResponseCallback) {
	this.colly.OnResponse(f)
}

func (this *Colly) OnScraped(f colly.ScrapedCallback) {
	this.colly.OnScraped(f)
}

func (this *Colly) Post(url string, requestData map[string]string) error {
	return this.colly.Post(url, requestData)
}
func (this *Colly) PostRaw(url string, requestData []byte) error {
	return this.colly.PostRaw(url, requestData)
}
func (this *Colly) Visit(url string) error {
	return this.colly.Visit(url)
}
