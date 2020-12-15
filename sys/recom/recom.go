package recom

func newSys(options Options) (sys *Recom, err error) {
	sys = &Recom{options: options}
	return
}

type Recom struct {
	options Options
}

func (this *Recom) RecommendItems(uId uint32, howmany int) (itemIds []uint32) {
	return
}
