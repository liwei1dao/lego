package core

import "sort"

const NotId = -1

func NewIndexer() *Indexer {
	set := new(Indexer)
	set.Indices = make(map[uint32]int)
	set.Ids = make([]uint32, 0)
	return set
}

type Indexer struct {
	Indices map[uint32]int
	Ids     []uint32
}

func (set *Indexer) Len() int {
	return len(set.Ids)
}
func (set *Indexer) Add(id uint32) {
	if _, exist := set.Indices[id]; !exist {
		set.Indices[id] = len(set.Ids)
		set.Ids = append(set.Ids, id)
	}
}
func (set *Indexer) ToIndex(id uint32) int {
	if denseId, exist := set.Indices[id]; exist {
		return denseId
	}
	return NotId
}
func (set *Indexer) ToID(index int) uint32 {
	return set.Ids[index]
}

func NewMarginalSubSet(indexer *Indexer, indices []int, values []uint8, subset []int) *MarginalSubSet {
	set := new(MarginalSubSet)
	set.Indexer = indexer
	set.Indices = indices
	set.Values = values
	set.SubSet = subset
	sort.Sort(set)
	return set
}

type MarginalSubSet struct {
	Indexer *Indexer
	Indices []int
	Values  []uint8
	SubSet  []int
}

func (set *MarginalSubSet) Len() int {
	return len(set.SubSet)
}

func (set *MarginalSubSet) Swap(i, j int) {
	set.SubSet[i], set.SubSet[j] = set.SubSet[j], set.SubSet[i]
}

func (set *MarginalSubSet) Less(i, j int) bool {
	return set.GetID(i) < set.GetID(j)
}

func (set *MarginalSubSet) GetIndex(i int) int {
	return set.Indices[set.SubSet[i]]
}

func (set *MarginalSubSet) GetID(i int) uint32 {
	index := set.GetIndex(i)
	return set.Indexer.ToID(index)
}
