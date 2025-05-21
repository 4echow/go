package main

import (
	"fmt"
	"sync"
)

func main() {
	ch := make(chan int, 20)
	var wg sync.WaitGroup
	var wg2 sync.WaitGroup
	wg.Add(2)
	wg2.Add(1)
	go func() {
		defer wg.Done()
		for i := range 10 {
			ch <- i
		}
	}()
	go func() {
		defer wg.Done()
		for i := range 10 {
			ch <- i
		}
	}()
	go func() {
		wg.Wait()
		close(ch)
	}()
	go func() {
		defer wg2.Done()
		for v := range ch {
			fmt.Println(v)
		}
	}()
	wg2.Wait()
}
