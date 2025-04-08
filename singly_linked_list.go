package main 

import "fmt"

type Node[T any] struct {
	next *Node[T]
	val  T
}

type List[T any] struct {
	head *Node[T]
}

func (li *List[T]) InsertBegin(val T) {
	if li.head == nil {
		li.head = &Node[T]{nil, val}
	} else {
		new_data := &Node[T]{li.head, val}
		li.head = new_data
	}
}

func (li *List[T]) InsertEnd(val T) {
	if li.head == nil {
		li.head = &Node[T]{nil, val}
	} else {
		node := li.head
		for node.next != nil {
			node = node.next
		}
		new_node := &Node[T]{nil, val}
		node.next = new_node
	}
}

func (li List[T]) Traverse() {
	if li.head == nil {
		fmt.Println("List is empty")
	} else {
		for node := li.head; node != nil; node = node.next {
			fmt.Println(node.val)
		}
	}
}

func main() {
	lst := List[int]{}
	lst.InsertBegin(6)
	lst.InsertBegin(5)
	lst.InsertBegin(4)
	lst.InsertEnd(7)
	lst.Traverse()
}
