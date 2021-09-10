package colly

import (
	tcolly "github.com/gocolly/colly/v2"
)

func newSys(options Options) (sys *Colly, err error) {
	sys = &Colly{options: options}

	return
}

type Colly struct {
	options Options
}

func (this *Colly) NewCollector() (c *tcolly.Collector) {
	c = tcolly.NewCollector()
	return
}
