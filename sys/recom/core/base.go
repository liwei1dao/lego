package core

import (
	"container/heap"
	"sort"

	"gonum.org/v1/gonum/floats"
)

const NotId = -1

func NewMatrixInt(row, col int) [][]int {
	ret := make([][]int, row)
	for i := range ret {
		ret[i] = make([]int, col)
	}
	return ret
}
func NewMatrixfloat64(row, col int) [][]float64 {
	ret := make([][]float64, row)
	for i := range ret {
		ret[i] = make([]float64, col)
	}
	return ret
}

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
func (set *Indexer) Add(id uint32) int {
	if i, exist := set.Indices[id]; !exist {
		set.Indices[id] = len(set.Ids)
		set.Ids = append(set.Ids, id)
		return set.Indices[id]
	} else {
		return i
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

func NewMarginalSubSet(indexer *Indexer, subset []int) *MarginalSubSet {
	set := new(MarginalSubSet)
	set.Indexer = indexer
	set.SubSet = subset
	sort.Sort(set)
	return set
}

type MarginalSubSet struct {
	Indexer *Indexer
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
	return set.SubSet[i]
}

func (set *MarginalSubSet) GetID(i int) uint32 {
	index := set.GetIndex(i)
	return set.Indexer.ToID(index)
}

func (set *MarginalSubSet) Contain(id uint32) bool {
	// if id is out of range
	if set.Len() == 0 || id < set.GetID(0) || id > set.GetID(set.Len()-1) {
		return false
	}
	// binary search
	low, high := 0, set.Len()-1
	for low <= high {
		// in bound
		if set.GetID(low) == id || set.GetID(high) == id {
			return true
		}
		mid := (low + high) / 2
		// in mid
		if id == set.GetID(mid) {
			return true
		} else if id < set.GetID(mid) {
			low = low + 1
			high = mid - 1
		} else if id > set.GetID(mid) {
			low = mid + 1
			high = high - 1
		}
	}
	return false
}

func NewMaxHeap(k int) *MaxHeap {
	knnHeap := new(MaxHeap)
	knnHeap.Elem = make([]interface{}, 0)
	knnHeap.Score = make([]float64, 0)
	knnHeap.K = k
	return knnHeap
}

type MaxHeap struct {
	Elem  []interface{} // store elements
	Score []float64     // store scores
	K     int           // the size of heap
}

type _HeapItem struct {
	Elem  interface{}
	Score float64
}

func (maxHeap *MaxHeap) Less(i, j int) bool {
	return maxHeap.Score[i] < maxHeap.Score[j]
}

func (maxHeap *MaxHeap) Swap(i, j int) {
	maxHeap.Elem[i], maxHeap.Elem[j] = maxHeap.Elem[j], maxHeap.Elem[i]
	maxHeap.Score[i], maxHeap.Score[j] = maxHeap.Score[j], maxHeap.Score[i]
}

func (maxHeap *MaxHeap) Len() int {
	return len(maxHeap.Elem)
}

func (maxHeap *MaxHeap) Push(x interface{}) {
	item := x.(_HeapItem)
	maxHeap.Elem = append(maxHeap.Elem, item.Elem)
	maxHeap.Score = append(maxHeap.Score, item.Score)
}

func (maxHeap *MaxHeap) Pop() interface{} {
	// Extract the minimum
	n := maxHeap.Len()
	item := _HeapItem{
		Elem:  maxHeap.Elem[n-1],
		Score: maxHeap.Score[n-1],
	}
	// Remove last element
	maxHeap.Elem = maxHeap.Elem[0 : n-1]
	maxHeap.Score = maxHeap.Score[0 : n-1]
	// We never use returned item
	return item
}

func (maxHeap *MaxHeap) Add(elem interface{}, score float64) {
	// Insert item
	heap.Push(maxHeap, _HeapItem{elem, score})
	// Remove minimum
	if maxHeap.Len() > maxHeap.K {
		heap.Pop(maxHeap)
	}
}

func (maxHeap *MaxHeap) ToSorted() ([]interface{}, []float64) {
	// sort indices
	scores := make([]float64, maxHeap.Len())
	indices := make([]int, maxHeap.Len())
	copy(scores, maxHeap.Score)
	floats.Argsort(scores, indices)
	// make output
	sorted := make([]interface{}, maxHeap.Len())
	for i := range indices {
		sorted[i] = maxHeap.Elem[indices[maxHeap.Len()-1-i]]
		scores[i] = maxHeap.Score[indices[maxHeap.Len()-1-i]]
	}
	return sorted, scores
}
