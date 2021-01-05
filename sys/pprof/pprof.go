package pprof

import (
	"fmt"
	"net/http"
	_ "net/http/pprof"
	"os"
	"time"

	"github.com/liwei1dao/lego/sys/log"

	"github.com/wolfogre/go-pprof-practice/animal"
)

func newSys(options Options) (sys *Pprof, err error) {
	sys = &Pprof{options: options}
	go sys.Start()
	return
}

type Pprof struct {
	options Options
}

func (this *Pprof) Start() {
	go func() {
		if err := http.ListenAndServe(fmt.Sprintf(":%d", this.options.ListenPort), nil); err != nil {
			log.Fatalf("Start Pprof Fatalf err:%v", err)
		}
		os.Exit(0)
	}()

	for {
		for _, v := range animal.AllAnimals {
			v.Live()
		}
		time.Sleep(time.Second)
	}
}
