package main

import "fmt"

func main() {
	womanEmoji := "\U0001F469"
	manEmoji := "\U0001F468"
	message := "Hi " + womanEmoji + " and " + manEmoji
	womanRune := []rune(message)[3]
	fmt.Println(string(womanRune))
}
