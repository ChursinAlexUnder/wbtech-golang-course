package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"
)

func main() {
	jobs := make(chan int)
	done := make(chan struct{})
	var wg sync.WaitGroup
	var once sync.Once

	for id := 1; id <= 3; id++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for j := range jobs {
				fmt.Printf("Воркер %d: получает job %d\n", id, j)
				time.Sleep(300 * time.Millisecond)
			}
			fmt.Printf("Воркер %d: Завершается\n", id)
		}(id)
	}

	go func() {
		defer close(jobs)
		i := 1
		for {
			select {
			case <-done:
				fmt.Println("Остановка")
				return
			case jobs <- i:
				i++
				time.Sleep(100 * time.Millisecond)
			}
		}
	}()

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt)

	<-sig
	fmt.Println("\nЗавершение программы...")

	once.Do(func() { close(done) })

	wg.Wait()
	fmt.Println("Программа завершилась")
}
