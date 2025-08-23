package main

import "fmt"

type Action struct {
	Human
	Name string
}

type Human struct {
	Name   string
	Age    int
	Height int
}

// Метод для структуры Human
func (human *Human) IsAdult() bool {
	return human.Age >= 18
}

func main() {
	var (
		a Action = Action{
			Human: Human{
				Name:   "Tolik",
				Age:    18,
				Height: 299,
			},
			Name: "Study",
		}
	)
	// Используем метод Human переменной типа Action
	fmt.Println(a.IsAdult())
}
