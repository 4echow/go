package main

type Team struct {
	name        string
	playerNames []string
}

type League struct {
	Teams map[string]Team
	Wins  map[string]int
}

func main() {
	return
}
