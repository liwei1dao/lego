package core

func NewDataSet(data map[uint32]map[uint32]float64) (set *DataSet) {
	set = new(DataSet)
	var (
		totalsorce           float64
		itemindex, userindex int
	)
	totalsorce = 0
	set.ratings = make([]int, 0)
	set.userIndexer = NewIndexer()
	set.itemIndexer = NewIndexer()

	for itemid, v := range data {
		set.itemIndexer.Add(itemid)
		for userid, _ := range v {
			set.userIndexer.Add(userid)
		}
	}

	userSubsetIndices := NewMatrixInt(set.UserCount(), 0)
	itemSubsetIndices := NewMatrixInt(set.ItemCount(), 0)
	set.scores = NewMatrixfloat64(set.ItemCount(), set.UserCount())
	for itemid, v := range data {
		for userid, sorce := range v {
			userindex = set.userIndexer.ToIndex(userid)
			itemindex = set.itemIndexer.ToIndex(itemid)
			set.scores[itemindex][userindex] = sorce
			userSubsetIndices[userindex] = append(userSubsetIndices[userindex], itemindex)
			itemSubsetIndices[itemindex] = append(itemSubsetIndices[itemindex], userindex)
			totalsorce += sorce
			set.ratings = append(set.ratings, userindex*set.ItemCount()+itemindex)
		}
	}

	set.globalmean = totalsorce / float64(len(set.ratings))
	set.users = make([]*MarginalSubSet, set.UserCount())
	set.items = make([]*MarginalSubSet, set.ItemCount())
	for u := range set.users {
		set.users[u] = NewMarginalSubSet(set.itemIndexer, userSubsetIndices[u])
	}
	for i := range set.items {
		set.items[i] = NewMarginalSubSet(set.userIndexer, itemSubsetIndices[i])
	}
	return
}

type (
	DataSetInterface interface {
		Count() int
		GlobalMean() float64
		UserCount() int
		ItemCount() int
		Users() []*MarginalSubSet
		Items() []*MarginalSubSet
		UserIndexer() *Indexer
		ItemIndexer() *Indexer
		User(userId uint32) *MarginalSubSet
		Item(itemId uint32) *MarginalSubSet
		UserByIndex(userIndex int) *MarginalSubSet
		ItemByIndex(itemIndex int) *MarginalSubSet
		GetWithIndex(i int) (int, int, float64)
	}
	DataSet struct {
		globalmean  float64
		ratings     []int
		scores      [][]float64
		userIndexer *Indexer
		itemIndexer *Indexer
		users       []*MarginalSubSet
		items       []*MarginalSubSet
	}
)

func (set *DataSet) Count() int {
	if set == nil {
		return 0
	}
	return len(set.ratings)
}

func (set *DataSet) GlobalMean() float64 {
	return set.globalmean
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

func (set *DataSet) Users() []*MarginalSubSet {
	return set.users
}

func (set *DataSet) Items() []*MarginalSubSet {
	return set.items
}

func (set *DataSet) GetWithIndex(i int) (userindex int, itemindex int, score float64) {
	itemindex = set.ratings[i] % set.ItemCount()
	userindex = set.ratings[i] / set.ItemCount()
	score = set.scores[itemindex][userindex]
	return
}
