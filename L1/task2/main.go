package main

import (
	"fmt"
	"sync"
)

func inSquare(number int, wg *sync.WaitGroup) {
	defer wg.Done()
	fmt.Println(number * number)
}

func main() {
	var (
		mas []int           = []int{2, 4, 6, 8, 10}
		wg  *sync.WaitGroup = new(sync.WaitGroup)
	)
	for _, value := range mas {
		wg.Add(1)
		go inSquare(value, wg)
	}
	wg.Wait()
}
