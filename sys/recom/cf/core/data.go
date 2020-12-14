package core

import "github.com/liwei1dao/lego/sys/recom/cf/base"

type DataSetInterface interface {
	Count() int
	UserCount() int
	ItemCount() int
	UserIndexer() *base.Indexer
	ItemIndexer() *base.Indexer
	Users() []*base.MarginalSubSet
	Items() []*base.MarginalSubSet
	SubSet(subset []int) DataSetInterface
	Item(itemId string) *base.MarginalSubSet
	User(userId string) *base.MarginalSubSet
	UserByIndex(userIndex int) *base.MarginalSubSet
	ItemByIndex(itemIndex int) *base.MarginalSubSet
}

func NewDataSet(userIDs, itemIDs []string, ratings []float64) *DataSet {
	set := new(DataSet)
	set.ratings = ratings
	set.userIndexer = base.NewIndexer()
	set.itemIndexer = base.NewIndexer()
	set.featureIndexer = base.NewStringIndexer()
	set.userIndices = make([]int, set.Count())
	set.itemIndices = make([]int, set.Count())
	for i := 0; i < set.Count(); i++ {
		userId := userIDs[i]
		itemId := itemIDs[i]
		set.userIndexer.Add(userId)
		set.itemIndexer.Add(itemId)
		set.userIndices[i] = set.userIndexer.ToIndex(userId)
		set.itemIndices[i] = set.itemIndexer.ToIndex(itemId)
	}
	userSubsetIndices := base.NewMatrixInt(set.UserCount(), 0)
	itemSubsetIndices := base.NewMatrixInt(set.ItemCount(), 0)
	for i := 0; i < set.Count(); i++ {
		userId := userIDs[i]
		itemId := itemIDs[i]
		userIndex := set.userIndexer.ToIndex(userId)
		itemIndex := set.itemIndexer.ToIndex(itemId)
		userSubsetIndices[userIndex] = append(userSubsetIndices[userIndex], i)
		itemSubsetIndices[itemIndex] = append(itemSubsetIndices[itemIndex], i)
	}
	set.users = make([]*base.MarginalSubSet, set.UserCount())
	set.items = make([]*base.MarginalSubSet, set.ItemCount())
	for u := range set.users {
		set.users[u] = base.NewMarginalSubSet(set.itemIndexer, set.itemIndices, ratings, userSubsetIndices[u])
	}
	for i := range set.items {
		set.items[i] = base.NewMarginalSubSet(set.userIndexer, set.userIndices, ratings, itemSubsetIndices[i])
	}
	return set
}

type DataSet struct {
	ratings        []float64
	userIndices    []int
	itemIndices    []int
	userIndexer    *base.Indexer
	itemIndexer    *base.Indexer
	featureIndexer *base.StringIndexer
	users          []*base.MarginalSubSet
	items          []*base.MarginalSubSet
	userFeatures   []*base.SparseVector
	itemFeatures   []*base.SparseVector
}

func (set *DataSet) Count() int {
	if set == nil {
		return 0
	}
	return len(set.ratings)
}
func (set *DataSet) UserIndexer() *base.Indexer {
	return set.userIndexer
}
func (set *DataSet) ItemIndexer() *base.Indexer {
	return set.itemIndexer
}
func (set *DataSet) UserCount() int {
	return set.UserIndexer().Len()
}
func (set *DataSet) ItemCount() int {
	return set.ItemIndexer().Len()
}

func (set *DataSet) SubSet(subset []int) DataSetInterface {
	return NewSubSet(set, subset)
}

func (set *DataSet) Users() []*base.MarginalSubSet {
	return set.users
}

func (set *DataSet) Items() []*base.MarginalSubSet {
	return set.items
}

func (set *DataSet) User(userId string) *base.MarginalSubSet {
	userIndex := set.userIndexer.ToIndex(userId)
	return set.UserByIndex(userIndex)
}

func (set *DataSet) Item(itemId string) *base.MarginalSubSet {
	itemIndex := set.itemIndexer.ToIndex(itemId)
	return set.ItemByIndex(itemIndex)
}

func (set *DataSet) UserByIndex(userIndex int) *base.MarginalSubSet {
	return set.users[userIndex]
}

func (set *DataSet) ItemByIndex(itemIndex int) *base.MarginalSubSet {
	return set.items[itemIndex]
}

func NewSubSet(dataSet *DataSet, subset []int) DataSetInterface {
	set := new(SubSet)
	set.DataSet = dataSet
	set.subset = subset
	// Index users' ratings and items' ratings
	userSubsetIndices := base.NewMatrixInt(set.UserCount(), 0)
	itemSubsetIndices := base.NewMatrixInt(set.ItemCount(), 0)
	for _, i := range set.subset {
		userIndex := set.userIndices[i]
		itemIndex := set.itemIndices[i]
		userSubsetIndices[userIndex] = append(userSubsetIndices[userIndex], i)
		itemSubsetIndices[itemIndex] = append(itemSubsetIndices[itemIndex], i)
	}
	set.users = make([]*base.MarginalSubSet, set.UserCount())
	set.items = make([]*base.MarginalSubSet, set.ItemCount())
	for u := range set.users {
		set.users[u] = base.NewMarginalSubSet(set.itemIndexer, set.itemIndices, set.ratings, userSubsetIndices[u])
	}
	for i := range set.items {
		set.items[i] = base.NewMarginalSubSet(set.userIndexer, set.userIndices, set.ratings, itemSubsetIndices[i])
	}
	return set
}

type SubSet struct {
	*DataSet                        // the existed dataset.
	subset   []int                  // indices of subset's samples
	users    []*base.MarginalSubSet // subsets of users' ratings
	items    []*base.MarginalSubSet // subsets of items' ratings
}
