package main

import (
	"fmt"

	"golang.org/x/tour/tree"
)

// Walk walks the tree t sending all values
// from the tree to the channel ch.
func Walk(t *tree.Tree, ch chan int) {

	if t == nil {
		return
	}

	Walk(t.Left, ch)
	ch <- t.Value
	Walk(t.Right, ch)
}

// Same determines whether the trees
// t1 and t2 contain the same values.
func Same(t1, t2 *tree.Tree) bool {
	ch1 := make(chan int)
	ch2 := make(chan int)
	go func() {
		defer close(ch1)
		Walk(t1, ch1)
	}()
	go func() {
		defer close(ch2)
		Walk(t2, ch2)
	}()
	for {
		v1, ok1 := <-ch1
		v2, ok2 := <-ch2
		if v1 != v2 {
			return false
		}
		if !ok1 && !ok2 {
			break
		}
	}
	return true
}

func main() {
	ch := make(chan int)
	tree1 := tree.New(1)
	tree2 := tree.New(2)
	go Walk(tree1, ch)
	for i := 0; i < 10; i++ {
		fmt.Println(<-ch)
	}
	myBool := Same(tree1, tree2)
	fmt.Println(myBool)
}
