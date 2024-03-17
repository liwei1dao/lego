package nlp

import (
	"github.com/go-ego/gse"
)

func newSys(options *Options) (sys *NLP, err error) {
	sys = &NLP{
		option: options,
	}
	err = sys.segmenter.LoadDict()
	return
}

type NLP struct {
	option    *Options
	segmenter gse.Segmenter
}

func (this *NLP) Cut(text string) []string {
	return this.segmenter.Cut(text)
}

func (this *NLP) Pos(text string) []gse.SegPos {
	return this.segmenter.Pos(text)
}
