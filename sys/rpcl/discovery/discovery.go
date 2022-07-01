package discovery

import (
	"time"

	"github.com/liwei1dao/lego/sys/rpcl/core"
)

type Discovery struct {
	sys   core.ISys
	store core.IStore
	dying chan struct{}
	done  chan struct{}
}

func (this *Discovery) Start() {
	if this.sys.UpdateInterval() > 0 {
		go func() {
			ticker := time.NewTicker(this.sys.UpdateInterval())
			defer ticker.Stop()
			defer this.store.Close()
			for {
				select {
				case <-this.dying:
					close(this.done)
					return
				case <-ticker.C:
				}

			}
		}()
	}
}
