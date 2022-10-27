package sitemap

type (
	xmlns int8
	base  struct {
		xmlns xmlns
	}
)

const (
	ImageXmlNS xmlns = 1
	VideoXmlNS xmlns = 2
	NewsXmlNS  xmlns = 4
)

func (b *base) setNs(xmlns xmlns) {
	b.xmlns = b.xmlns | xmlns
}
