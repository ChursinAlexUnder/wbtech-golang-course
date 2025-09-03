package main

import "fmt"

func goroutinOutput(c chan int) {
	for {
		data := <-c
		fmt.Println(data)
	}
}

func main() {
	var (
		n, num int
	)
	fmt.Scan(&n)
	c := make(chan int, n)
	for i := 0; i < n; i++ {
		go goroutinOutput(c)
	}
	for {
		fmt.Scan(&num)
		for i := 0; i < n; i++ {
			c <- num
		}
	}
}
