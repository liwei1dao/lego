package base

import (
	"container/heap"
	"sort"

	"gonum.org/v1/gonum/floats"
)

const NotId = -1

func NewIndexer() *Indexer {
	set := new(Indexer)
	set.Indices = make(map[string]int)
	set.IDs = make([]string, 0)
	return set
}

type Indexer struct {
	Indices map[string]int
	IDs     []string
}

func (set *Indexer) Len() int {
	return len(set.IDs)
}
func (set *Indexer) Add(ID string) {
	if _, exist := set.Indices[ID]; !exist {
		set.Indices[ID] = len(set.IDs)
		set.IDs = append(set.IDs, ID)
	}
}
func (set *Indexer) ToIndex(ID string) int {
	if denseId, exist := set.Indices[ID]; exist {
		return denseId
	}
	return NotId
}
func (set *Indexer) ToID(index int) string {
	return set.IDs[index]
}

func NewStringIndexer() *StringIndexer {
	set := new(StringIndexer)
	set.Indices = make(map[string]int)
	set.Names = make([]string, 0)
	return set
}

type StringIndexer struct {
	Indices map[string]int
	Names   []string
}

func NewMarginalSubSet(indexer *Indexer, indices []int, values []float64, subset []int) *MarginalSubSet {
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
	Values  []float64
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
func (set *MarginalSubSet) GetID(i int) string {
	index := set.GetIndex(i)
	return set.Indexer.ToID(index)
}

func (set *MarginalSubSet) Contain(id string) bool {
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

type SparseVector struct {
	Indices []int
	Values  []float64
	Sorted  bool
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
