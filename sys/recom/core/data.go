package core

func NewDataSet(userIds, itemIds []uint32, ratings []uint8) (set *DataSet) {
	set = new(DataSet)
	set.count = len(ratings)
	set.userIndexer = NewIndexer()
	set.itemIndexer = NewIndexer()
	set.userIndices = make([]int, len(userIds))
	set.itemIndices = make([]int, len(itemIds))
	for i := 0; i < set.Count(); i++ {
		userId := userIds[i]
		itemId := itemIds[i]
		set.userIndexer.Add(userId)
		set.itemIndexer.Add(itemId)
		set.userIndices[i] = set.userIndexer.ToIndex(userId)
		set.itemIndices[i] = set.itemIndexer.ToIndex(itemId)
	}
	userSubsetIndices := NewMatrixInt(set.UserCount(), 0)
	itemSubsetIndices := NewMatrixInt(set.ItemCount(), 0)
	for i := 0; i < set.Count(); i++ {
		userId := userIds[i]
		itemId := itemIds[i]
		userIndex := set.userIndexer.ToIndex(userId)
		itemIndex := set.itemIndexer.ToIndex(itemId)
		userSubsetIndices[userIndex] = append(userSubsetIndices[userIndex], i)
		itemSubsetIndices[itemIndex] = append(itemSubsetIndices[itemIndex], i)
	}
	set.users = make([]*MarginalSubSet, set.UserCount())
	set.items = make([]*MarginalSubSet, set.ItemCount())
	for u := range set.users {
		set.users[u] = NewMarginalSubSet(set.itemIndexer, set.itemIndices, ratings, userSubsetIndices[u])
	}
	for i := range set.items {
		set.items[i] = NewMarginalSubSet(set.userIndexer, set.userIndices, ratings, itemSubsetIndices[i])
	}
	return
}

type (
	DataSetInterface interface {
		Count() int
		UserCount() int
		ItemCount() int
		UserIndexer() *Indexer
		ItemIndexer() *Indexer
		User(userId uint32) *MarginalSubSet
		Item(itemId uint32) *MarginalSubSet
		UserByIndex(userIndex int) *MarginalSubSet
		ItemByIndex(itemIndex int) *MarginalSubSet
	}
	DataSet struct {
		count       int
		userIndices []int
		itemIndices []int
		userIndexer *Indexer
		itemIndexer *Indexer
		users       []*MarginalSubSet
		items       []*MarginalSubSet
	}
)

func (set *DataSet) Count() int {
	return set.count
}
func (set *DataSet) UserIndexer() *Indexer {
	return set.userIndexer
}
func (set *DataSet) ItemIndexer() *Indexer {
	return set.itemIndexer
}

func (set *DataSet) UserCount() int {
	return set.UserIndexer().Len()
}

func (set *DataSet) ItemCount() int {
	return set.ItemIndexer().Len()
}

func (set *DataSet) UserByIndex(userIndex int) *MarginalSubSet {
	return set.users[userIndex]
}

func (set *DataSet) ItemByIndex(itemIndex int) *MarginalSubSet {
	return set.items[itemIndex]
}

func (set *DataSet) User(userId uint32) *MarginalSubSet {
	userIndex := set.userIndexer.ToIndex(userId)
	return set.UserByIndex(userIndex)
}

func (set *DataSet) Item(itemId uint32) *MarginalSubSet {
	itemIndex := set.itemIndexer.ToIndex(itemId)
	return set.ItemByIndex(itemIndex)
}
