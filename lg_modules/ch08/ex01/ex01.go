package main

import "fmt"

type DoubleAble interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr |
		~float32 | ~float64
}

func DoubleNum[T DoubleAble](dn T) T {
	return dn * 2
}

func main() {
	var v1 int32
	var v2 float64
	v1 = 30
	v2 = 42.42
	r1 := DoubleNum(v1)
	r2 := DoubleNum(v2)
	fmt.Printf("value:%v, type:%T\n", r1, r1)
	fmt.Printf("value:%v, type:%T\n", r2, r2)
}
