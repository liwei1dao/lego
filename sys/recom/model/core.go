package model

import (
	"github.com/liwei1dao/lego/sys/recom/core"
)

type (
	IModel interface {
		SetParams(params Params)
		GetParams() Params
		Predict(userId, itemId uint32) float64
		Fit(trainSet core.DataSetInterface)
	}
)

func Top(items map[uint32]bool, userId uint32, n int, exclude *core.MarginalSubSet, model IModel) ([]uint32, []float64) {
	// Get top-n list
	itemsHeap := core.NewMaxHeap(n)
	for itemId := range items {
		if !exclude.Contain(itemId) {
			itemsHeap.Add(itemId, model.Predict(userId, itemId))
		}
	}
	elem, scores := itemsHeap.ToSorted()
	recommends := make([]uint32, len(elem))
	for i := range recommends {
		recommends[i] = elem[i].(uint32)
	}
	return recommends, scores
}

func Items(dataSet ...core.DataSetInterface) map[uint32]bool {
	items := make(map[uint32]bool)
	for _, core := range dataSet {
		for i := 0; i < core.ItemCount(); i++ {
			itemId := core.ItemIndexer().ToID(i)
			items[itemId] = true
		}
	}
	return items
}
