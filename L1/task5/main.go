package main

import (
	"fmt"
	"log"
	"sync"
	"time"
)

func main() {
	var N int
	log.Println("Введите число секунд N: ")
	_, err := fmt.Scan(&N)
	if err != nil || N <= 0 {
		log.Println("Ошибка ввода!")
		return
	}

	ch := make(chan int)
	quit := make(chan struct{})
	var wg sync.WaitGroup

	wg.Add(1)
	go func() {
		defer wg.Done()
		i := 1
		for {
			select {
			case <-quit:
				close(ch)
				return
			case ch <- i:
				i++
				time.Sleep(time.Second)
			}
		}
	}()

	wg.Add(1)
	go func() {
		defer wg.Done()
		for v := range ch {
			log.Println(v)
		}
	}()

	go func() {
		<-time.After(time.Duration(N) * time.Second)
		close(quit)
	}()

	wg.Wait()
	log.Println("Программа завершена.")
}
