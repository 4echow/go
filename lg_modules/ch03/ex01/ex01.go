package main

import "fmt"

func main() {
	greetings := []string{"Hello", "Hola", "नमस्कार", "こんにちは", "Привіт"}
	a := greetings[:2]
	b := greetings[1:4]
	c := greetings[3:]
	fmt.Println(a)
	fmt.Println(b)
	fmt.Println(c)
}
