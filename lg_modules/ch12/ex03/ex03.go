package main

import (
	"fmt"
	"math"
	"sync"
)

func initMap() map[int]float64 {
	rootMap := make(map[int]float64, 100000)
	for i := range 100000 {
		rootMap[i] = math.Sqrt(float64(i))
	}

	return rootMap
}

var initRootMapCached func() map[int]float64 = sync.OnceValue(initMap)

func GetRoot(v int) float64 {
	rootMap := initRootMapCached()
	return rootMap[v]
}

func main() {
	for i := 0; i < 100000; i += 1000 {
		fmt.Println(i, GetRoot(i))
	}
}
