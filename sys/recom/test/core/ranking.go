package core

import "github.com/liwei1dao/lego/sys/recom/cf/base"

func Top(items map[string]bool, userId string, n int, exclude *base.MarginalSubSet, model ModelInterface) ([]string, []float64) {
	// Get top-n list
	itemsHeap := base.NewMaxHeap(n)
	for itemId := range items {
		if !exclude.Contain(itemId) {
			itemsHeap.Add(itemId, model.Predict(userId, itemId))
		}
	}
	elem, scores := itemsHeap.ToSorted()
	recommends := make([]string, len(elem))
	for i := range recommends {
		recommends[i] = elem[i].(string)
	}
	return recommends, scores
}

func Items(dataSet ...DataSetInterface) map[string]bool {
	items := make(map[string]bool)
	for _, data := range dataSet {
		for i := 0; i < data.ItemCount(); i++ {
			itemId := data.ItemIndexer().ToID(i)
			items[itemId] = true
		}
	}
	return items
}
