package main

import (
	"context"
	"fmt"
	"time"

	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/go-app/database"
	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/go-app/internal"
	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/go-app/router"
)

func main() {
	// Если topic-а с необходимым названием для producer не существует,
	// то он создастся по этим параметрам
	var (
		topic             string = "orders"
		partitions        int    = 3
		replicationFactor int    = 1
	)
	ctx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	// Подключение к локальной базе данных
	pool, err := database.InitDB(ctx)
	if err != nil {
		fmt.Printf("Не удалось подключиться к бд: %v\n", err)
		return
	}
	defer pool.Close()

	// Запускаем сервер
	router := router.SetupRouter(pool, ctx)
	err = router.Run(":8081")
	if err != nil {
		fmt.Printf("Не удалось запустить сервер: %v\n", err)
		return
	}

	// Запуск producer в горутине
	go internal.Producer(ctx, topic, partitions, replicationFactor)

	// Запуск consumer
	internal.Consumer(ctx)

}
