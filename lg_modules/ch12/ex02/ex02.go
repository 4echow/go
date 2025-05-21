package main

import (
	"fmt"
	"sync"
)

func main() {
	ch1 := make(chan int, 10)
	ch2 := make(chan int, 10)
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup
	wg.Add(2)
	wg2.Add(1)

	go func() {
		defer wg.Done()
		for i := range 10 {
			ch1 <- i
		}
	}()

	go func() {
		defer wg.Done()
		for i := range 10 {
			ch2 <- i
		}
	}()

	go func() {
		wg.Wait()
		close(ch1)
		close(ch2)
	}()

	go func() {
		defer wg2.Done()
		count := 2
		for count > 0 {
			select {
			case v, ok := <-ch1:
				if !ok {
					ch1 = nil
					count--
				}
				fmt.Printf("%v channel writes: %v\n", ch1, v)
			case v, ok := <-ch2:
				if !ok {
					ch2 = nil
					count--
				}
				fmt.Printf("%v channel writes: %v\n", ch2, v)
			}
		}
	}()
	wg2.Wait()
}
