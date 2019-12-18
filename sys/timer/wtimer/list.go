package wtimer

import (
	"fmt"
	"lego/sys/log"
)

//定义节点
type Node struct {
	data interface{}
	prev *Node
	next *Node
}

type LinkedList struct {
	head   *Node
	last   *Node
	length uint
}

func NewLinkedList() *LinkedList {
	var list *LinkedList = new(LinkedList)
	list.head = nil
	list.last = nil
	list.length = 0
	return list
}

/**
 * 获取表头
 */
func (this LinkedList) GetHead() *Node {
	return this.head
}

/**
 * 获取表尾
 */
func (this LinkedList) GetLast() *Node {
	return this.last
}

func (this LinkedList) Length() uint {
	return this.length
}

func (this *LinkedList) PushBack(node Node) *Node {
	node.next = nil
	if nil == this.head { //空表
		this.head = &node
		this.head.prev = nil
		this.last = this.head
	} else {
		node.prev = this.last
		this.last.next = &node
		this.last = this.last.next
	}
	log.Errorf("insert %d %d\n", this.length, this.last.Data)
	this.length++
	return this.last
}

func (this *LinkedList) erase(node *Node) {
	fmt.Println(node)
	if nil == node {
		return
	} else if nil == node.next && nil == node.next {
		return
	}
	if node == this.head && node == this.last {
		this.head = nil
		this.last = nil
		this.length = 0
	} else {
		if node == this.head {
			this.head = this.head.next
			if nil != this.head {
				this.head.prev = nil
			}
		} else if node == this.last {
			node.prev.next = nil
			this.last = node.prev
		} else {
			node.prev.next = node.next
			node.next.prev = node.prev
		}
	}
	this.length--
}

func Delete(node *Node) {
	if nil == node {
		return
	} else if nil == node.prev { //该元素处于表头，不删除，默认表头不存元素
		return
	} else if nil == node.next { //该元素处于表尾
		node.prev.next = nil
		node.prev = nil
	} else {
		node.next.prev = node.prev
		node.prev.next = node.next
		node.prev = nil
		node.next = nil
	}
}

func (this *Node) InsertHead(node Node) *Node { //从表头插入
	if nil == this || nil != this.prev { //为空，或者不是表头(表头的prev为空)
		return nil
	} else {
		if nil != this.next {
			this.next.prev = &node
			node.next = this.next
		}
		this.next = &node
		node.prev = this
	}
	return &node
}

func (this *Node) Next() (node *Node) {
	return this.next
}

func (this *Node) Prev() (node *Node) {
	return this.prev
}

func (this *Node) Data() (data interface{}) {
	return this.data
}

func (this *Node) SetData(data interface{}) {
	this.data = data
}
