package main

import (
	"fmt"
	"math/rand"
)

func main() {
	var randNumSlice []int = make([]int, 100)
	for i := range 100 {
		randNumSlice[i] = rand.Intn(100)
	}
	fmt.Println(randNumSlice)
}
