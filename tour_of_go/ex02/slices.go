package main

import "golang.org/x/tour/pic"

func Pic(dx, dy int) [][]uint8 {
	myslice := make([][]uint8, dy)
	for y := range myslice {
		myslice[y] = make([]uint8, dx)
		for x := range myslice[y] {
			myslice[y][x] = uint8((x + y) / 2 * (x ^ y))
		}
	}

	return myslice
}

func main() {
	pic.Show(Pic)
}
