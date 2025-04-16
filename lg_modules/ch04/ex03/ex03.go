package main

import "fmt"

func main() {
	var total int
	for i := range 10 {
		total := total + i // := assigns new variable that shadows total variable in main() scope
		fmt.Println(total)
	}
	fmt.Println(total)
}
