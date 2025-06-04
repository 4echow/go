package main

import (
	"context"
	"fmt"
	"math/rand"
	"time"
)

func get1234(ctx context.Context) (int, int, string) {
	var sum int = 0
	const magicnumber int = 1234
	i := 0
	for {
		randInt := rand.Intn(100_000_000)
		select {
		case <-ctx.Done():
			return sum, i, "timeout."
		default:
			sum += randInt
			if randInt == magicnumber {
				return sum, i, "1234 generated."
			}
			i++
		}
	}
}

func main() {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()
	sum, i, cause := get1234(ctx)
	fmt.Printf("Sum: %v, iterations: %v, Finish with cause: %s\n", sum, i, cause)
}
