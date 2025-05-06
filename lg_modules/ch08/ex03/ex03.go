package main

import (
	"errors"
	"fmt"
)

type Node[T comparable] struct {
	Value T
	Next  *Node[T]
}

type LList[T comparable] struct {
	Head *Node[T]
	Tail *Node[T]
}

func (l *LList[T]) Add(v T) {
	n := &Node[T]{Value: v}
	if l.Head == nil {
		l.Head = n
		l.Tail = l.Head
		return
	}
	l.Tail.Next = n
	l.Tail = l.Tail.Next
}

func (l *LList[T]) Insert(v T, idx int) error {
	if idx < 0 {
		return errors.New("Index cannot be negative")
	}

	n := &Node[T]{Value: v}
	currNode := l.Head
	var listIdx int
	for currNode != nil {
		listIdx++
		currNode = currNode.Next
	}
	if idx > listIdx {
		return errors.New("Index out of bounds")
	}
	if idx == listIdx {
		l.Tail.Next = n
		l.Tail = l.Tail.Next
		return nil
	}
	if idx == 0 && listIdx == 0 {
		l.Head = n
		l.Tail = l.Head
		return nil
	}
	if idx == 0 && listIdx != 0 {
		n.Next = l.Head
		l.Head = n
		return nil
	}
	currNode2 := l.Head
	var counter int
	for currNode2 != nil {
		if counter == idx-1 {
			n.Next = currNode2.Next
			currNode2.Next = n
			return nil
		}
		currNode2 = currNode2.Next
		counter++
	}
	return nil
}

func (l *LList[T]) Index(v T) int {
	currNode := l.Head
	var idx int
	for currNode != nil {
		if currNode.Value == v {
			return idx
		}
		idx++
		currNode = currNode.Next
	}
	return -1
}

func main() {
	myList := &LList[int]{}
	myList.Add(1)
	myList.Add(2)
	myList.Add(3)
	fmt.Println(myList.Index(1), myList.Index(2), myList.Index(3))
	err := myList.Insert(4, 0)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(myList.Index(1), myList.Index(2), myList.Index(3), myList.Index(4))
	err2 := myList.Insert(5, 4)
	if err2 != nil {
		fmt.Println(err2)
	}
	fmt.Println(myList.Index(1), myList.Index(2), myList.Index(3), myList.Index(4), myList.Index(5))
}
