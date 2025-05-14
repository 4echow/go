package main

import "golang.org/x/exp/constraints"

type Number interface {
	constraints.Float | constraints.Integer
}

// Add adds two Integer numbers and returns their sum
//
// More info on the addition operation on [Mathisfun]
// [Mathisfun]: https://www.mathsisfun.com/numbers/addition.html
func Add[T1, T2 Number](a, b T1) T2 {
	return T2(a + b)
}
