package main

import (
	"fmt"
	"sync"
)

func main() {
	sizes := make(chan int)
	var wg sync.WaitGroup
	nums := 10
	for num := 0; num < nums; num++ {
		wg.Add(1)
		go func(num int) {
			defer wg.Done()
			sizes <- num
		}(num)
	}

	go func() {
		wg.Wait()
		close(sizes)
	}()

	for rt := range sizes {
		fmt.Println(rt)
	}
}
