package recom

func newSys(options Options) (sys *Recom, err error) {
	sys = &Recom{options: options}
	return
}

type Recom struct {
	options Options
}

