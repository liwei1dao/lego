package recom

import (
	"sync"

	"github.com/liwei1dao/lego/sys/recom/core"
	"github.com/liwei1dao/lego/sys/recom/model"
)

func newSys(options Options) (sys *Recom, err error) {
	sys = new(Recom)
	sys.dataset = core.NewDataSet(options.ItemIdsScore)
	sys.model = model.NewSVD(model.Params{
		model.NFactors:   10,
		model.Reg:        0.01,
		model.Lr:         0.05,
		model.NEpochs:    100,
		model.InitMean:   0,
		model.InitStdDev: 0.001,
	})
	sys.wg.Add(1)
	go sys.Fit()
	return
}

type Recom struct {
	dataset core.DataSetInterface
	model   model.IModel
	wg      sync.WaitGroup
}

func (this *Recom) Fit() {
	this.model.Fit(this.dataset)
	this.wg.Done()
}

func (this *Recom) Wait() {
	this.wg.Wait()
}

func (this *Recom) RecommendItems(uId uint32, howmany int) (itemIds []uint32) {
	excludeItems := this.dataset.User(uId)
	itemIds, _ = model.Top(model.Items(this.dataset), uId, howmany, excludeItems, this.model)
	return
}
