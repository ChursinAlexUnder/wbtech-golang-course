package main

import (
	"context"
	"log"

	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/internal"
	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/internal/database"
	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/internal/router"
	"github.com/google/uuid"
	"github.com/hashicorp/golang-lru/v2/expirable"
)

// @title           Сервис получения информации о заказе
// @version         1.0
// @description     Данный сервис позволяет получить всю информацию о заказе по его order_uid

// @host      localhost:8081
// @BasePath  /

func main() {
	// Если topic-а с необходимым названием для producer не существует,
	// то он создастся по этим параметрам
	var (
		topic             string = "orders"
		partitions        int    = 3
		replicationFactor int    = 1

		// Для кеша
		orders []database.Orders
	)
	ctx := context.Background()

	// Подключение к локальной базе данных
	pool, err := database.InitDB(ctx)
	if err != nil {
		log.Printf("Не удалось подключиться к бд: %v\n", err)
		return
	}
	defer pool.Close()

	cache := expirable.NewLRU[uuid.UUID, database.Orders](1000, nil, 0)

	orders, err = database.SelectOrdersForCache(ctx, pool)
	if err != nil {
		log.Printf("Не удалось получить актуальные заказы для загрузки кеша из бд: %v\n", err)
		return
	}

	for _, order := range orders {
		cache.Add(order.Order_uid, order)
	}
	log.Printf("Кеш успешно заполнен!\nТекущее количество записей в кеше %d\n", cache.Len())

	// Запуск producer в горутине
	go internal.Producer(ctx, topic, partitions, replicationFactor)

	// Запуск consumer
	go internal.Consumer(ctx, pool, cache)

	// Запускаем сервер
	router := router.SetupRouter(ctx, pool, cache)
	err = router.Run(":8081")
	if err != nil {
		log.Printf("Не удалось запустить сервер: %v\n", err)
		return
	}
}
