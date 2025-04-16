package main

import "fmt"

func prefixer(prefix string) func(string) string {
	return func(name string) string {
		return prefix + " " + name
	}
}

func main() {
	helloPrefix := prefixer("Hello")
	fmt.Println(helloPrefix("Bob"))
	fmt.Println(helloPrefix("Maria"))
}
