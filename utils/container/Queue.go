package container

import (
	"container/list"
	"fmt"
	"sync"
)

type Queue struct {
	sync.Mutex
	data *list.List
}

func NewQueue() *Queue {
	q := new(Queue)
	q.data = list.New()
	return q
}

func (q *Queue) Len() int {
	return q.data.Len()
}

func (q *Queue) Push(v interface{}) {
	defer q.Unlock()
	q.Lock()
	q.data.PushFront(v)
}

func (q *Queue) Pop() interface{} {
	defer q.Unlock()
	q.Lock()
	iter := q.data.Back()
	if iter == nil {
		return nil
	}
	v := iter.Value
	q.data.Remove(iter)
	return v
}

func (q *Queue) Dump() {
	for iter := q.data.Back(); iter != nil; iter = iter.Prev() {
		fmt.Println("item:", iter.Value)
	}
}
