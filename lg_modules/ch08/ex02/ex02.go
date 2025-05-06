package main

import (
	"fmt"
	"strconv"
)

type Printable interface {
	fmt.Stringer
	~int | ~float64
}

type PrintableInt int

func (p PrintableInt) String() string {
	return strconv.Itoa(int(p))
}

type PrintableFloat float64

func (p PrintableFloat) String() string {
	return strconv.FormatFloat(float64(p), 'f', 2, 64)
}

func PrintGeneric[T Printable](p T) {
	fmt.Println(p)
}

func main() {
	var v1 PrintableInt
	var v2 PrintableFloat
	v1 = 30
	v2 = 42.42
	PrintGeneric(v1)
	PrintGeneric(v2)
}
