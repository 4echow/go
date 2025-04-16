package main

import (
	"fmt"
	"math/rand"
)

func main() {
	var randNumSlice []int = make([]int, 100)
	for i := range 100 {
		randNumSlice[i] = rand.Intn(100)
		switch {
		case randNumSlice[i]%3 == 0 && randNumSlice[i]%2 == 0:
			fmt.Println("Six!")
		case randNumSlice[i]%2 == 0:
			fmt.Println("Two!")
		case randNumSlice[i]%3 == 0:
			fmt.Println("Three!")
		default:
			fmt.Println("Never mind")
		}
	}

}
