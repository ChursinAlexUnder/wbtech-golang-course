package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

		// Для горутин
		wg sync.WaitGroup
	)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Подключение к локальной базе данных
	pool, err := database.InitDB(ctx)
	if err != nil {
		log.Fatalf("Не удалось подключиться к бд: %v\n", err)
		return
	}
	defer pool.Close()

	cache := expirable.NewLRU[uuid.UUID, database.Orders](1000, nil, 0)

	orders, err = database.SelectOrdersForCache(ctx, pool)
	if err != nil {
		log.Fatalf("Не удалось получить актуальные заказы для загрузки кеша из бд: %v\n", err)
		return
	}

	for _, order := range orders {
		cache.Add(order.Order_uid, order)
	}
	log.Printf("Кеш успешно заполнен!\nТекущее количество записей в кеше %d\n", cache.Len())

	wg.Add(1)

	// Запуск producer в горутине
	go func() {
		defer wg.Done()
		internal.Producer(ctx, topic, partitions, replicationFactor)
	}()

	// Запуск consumer
	wg.Add(1)
	go func() {
		defer wg.Done()
		internal.Consumer(ctx, pool, cache)
	}()

	// Запускаем сервер
	router := router.SetupRouter(ctx, pool, cache)
	srv := &http.Server{
		Addr:    ":8081",
		Handler: router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Не удалось запустить сервер: %v\n", err)
		}
	}()

	log.Println("Приложение успешно запустилось")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	log.Println("Приложение плавно выключается")
	cancel()

	if err := srv.Shutdown(context.Background()); err != nil {
		log.Fatalf("Ошибка плавного выключения приложения: %v", err)
	}

	done := make(chan struct{})
	go func() {
		wg.Wait()
		close(done)
	}()

	select {
	case <-done:
		// всё завершилось нормально
	case <-time.After(5 * time.Second):
		log.Println("Timeout ожидания фоновых горутин, принудительное завершение")
	}

	log.Println("Приложение выключилось")
}
