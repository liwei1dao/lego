package model

import (
	"github.com/liwei1dao/lego/sys/recom/data"
)

type (
	IModel interface {
		SetParams(params Params)
		GetParams() Params
		Predict(userId, itemId uint32) float64
		Fit(trainSet data.DataSetInterface)
	}
)

func Top(items map[uint32]bool, userId uint32, n int, exclude *data.MarginalSubSet, model IModel) ([]uint32, []float64) {
	// Get top-n list
	itemsHeap := data.NewMaxHeap(n)
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

func Items(dataSet ...data.DataSetInterface) map[uint32]bool {
	items := make(map[uint32]bool)
	for _, data := range dataSet {
		for i := 0; i < data.ItemCount(); i++ {
			itemId := data.ItemIndexer().ToID(i)
			items[itemId] = true
		}
	}
	return items
}
