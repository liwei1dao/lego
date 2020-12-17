package model

import (
	"github.com/liwei1dao/lego/sys/recom/data"
)

type ModelBase struct {
	Params            Params
	UserIndexer       *data.Indexer        // Users' ID set
	ItemIndexer       *data.Indexer        // Items' ID set
	rng               data.RandomGenerator // Random generator
	randState         int64                // Random seed
	isSetParamsCalled bool
}

func (model *ModelBase) SetParams(params Params) {
	model.isSetParamsCalled = true
	model.Params = params
}

func (model *ModelBase) GetParams() Params {
	return model.Params
}

func (model *ModelBase) Init(trainSet data.DataSetInterface) {
	// Check Base.SetParams() called
	if model.isSetParamsCalled == false {
		panic("Base.SetParams() not called")
	}
	model.UserIndexer = trainSet.UserIndexer()
	model.ItemIndexer = trainSet.ItemIndexer()
	// Setup random state
	model.rng = data.NewRandomGenerator(model.randState)
}

func NewItemPop(params Params) *ItemPop {
	pop := new(ItemPop)
	pop.SetParams(params)
	return pop
}

type ItemPop struct {
	ModelBase
	Pop []float64
}

func (pop *ItemPop) Predict(userId, itemId uint32) float64 {
	// Return items' popularity
	denseItemId := pop.ItemIndexer.ToIndex(itemId)
	if denseItemId == data.NotId {
		return 0
	}
	return pop.Pop[denseItemId]
}

func (pop *ItemPop) Fit(set data.DataSetInterface) {
	pop.Init(set)
	// Get items' popularity
	pop.Pop = make([]float64, set.ItemCount())
	for i := 0; i < set.ItemCount(); i++ {
		pop.Pop[i] = float64(set.ItemByIndex(i).Len())
	}
}
