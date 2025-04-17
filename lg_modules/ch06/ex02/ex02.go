package main

import "fmt"

func UpdateSlice(strSlice []string, str string) {
	strSlice[len(strSlice)-1] = str
	fmt.Println(strSlice)
}

func GrowSlice(strSlice []string, str string) {
	strSlice = append(strSlice, str)
	fmt.Println(strSlice)
}

func main() {
	mySlice := []string{"A", "B", "C"}
	fmt.Println(mySlice)
	UpdateSlice(mySlice, "FOO")
	fmt.Println(mySlice)
	GrowSlice(mySlice, "BAR")
	fmt.Println(mySlice)
}
