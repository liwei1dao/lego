package model

import (
	"github.com/liwei1dao/lego/sys/recom/cf/base"
	"github.com/liwei1dao/lego/sys/recom/cf/core"
)

type Base struct {
	Params      base.Params          // Hyper-parameters
	UserIndexer *base.Indexer        // Users' ID set
	ItemIndexer *base.Indexer        // Items' ID set
	rng         base.RandomGenerator // Random generator
	randState   int64                // Random seed
	// Tracker
	isSetParamsCalled bool // Check whether SetParams called
}

func (model *Base) SetParams(params base.Params) {
	model.isSetParamsCalled = true
	model.Params = params
	model.randState = model.Params.GetInt64(base.RandomState, 0)
}
func (model *Base) GetParams() base.Params {
	return model.Params
}

func (model *Base) Init(trainSet core.DataSetInterface) {
	// Check Base.SetParams() called
	if model.isSetParamsCalled == false {
		panic("Base.SetParams() not called")
	}
	// Setup ID set
	model.UserIndexer = trainSet.UserIndexer()
	model.ItemIndexer = trainSet.ItemIndexer()
	// Setup random state
	model.rng = base.NewRandomGenerator(model.randState)
}

func NewItemPop(params base.Params) *ItemPop {
	pop := new(ItemPop)
	pop.SetParams(params)
	return pop
}

type ItemPop struct {
	Base
	Pop []float64
}

func (pop *ItemPop) Fit(set core.DataSetInterface, options *base.RuntimeOptions) {
	pop.Init(set)
	// Get items' popularity
	pop.Pop = make([]float64, set.ItemCount())
	for i := 0; i < set.ItemCount(); i++ {
		pop.Pop[i] = float64(set.ItemByIndex(i).Len())
	}
}

// Predict by the ItemPop model.
func (pop *ItemPop) Predict(userId, itemId string) float64 {
	// Return items' popularity
	denseItemId := pop.ItemIndexer.ToIndex(itemId)
	if denseItemId == base.NotId {
		return 0
	}
	return pop.Pop[denseItemId]
}
