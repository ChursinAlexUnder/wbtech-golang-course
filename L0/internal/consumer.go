package internal

import (
	"context"
	"encoding/json"
	"log"
	"time"
	"unicode/utf8"

	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/internal/database"
	"github.com/google/uuid"
	"github.com/hashicorp/golang-lru/v2/expirable"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/segmentio/kafka-go"
)

func IsValidDataFromKafka(orderJson []byte) bool {
	var (
		orderStruct database.Orders
		err         error
	)
	err = json.Unmarshal(orderJson, &orderStruct)
	if err != nil {
		return false
	}
	if utf8.RuneCountInString(orderStruct.Locale) > 10 ||
		utf8.RuneCountInString(orderStruct.Payment.Currency) > 10 ||
		utf8.RuneCountInString(orderStruct.Payment.Provider) > 50 ||
		utf8.RuneCountInString(orderStruct.Payment.Bank) > 50 ||
		utf8.RuneCountInString(orderStruct.Delivery.Phone) > 30 ||
		utf8.RuneCountInString(orderStruct.Delivery.Email) > 254 ||
		orderStruct.Payment.Transaction != orderStruct.Order_uid ||
		orderStruct.Delivery.Uid != orderStruct.Delivery_uid {
		return false
	}
	for _, item := range orderStruct.Items {
		if utf8.RuneCountInString(item.Size) > 10 ||
			utf8.RuneCountInString(item.Brand) > 150 ||
			item.Track_number != orderStruct.Track_number {
			return false
		}
	}
	return true
}

func Consumer(ctx context.Context, pool *pgxpool.Pool, cache *expirable.LRU[uuid.UUID, database.Orders]) {
	var order database.Orders

	dialer := &kafka.Dialer{
		Timeout:   20 * time.Second,
		DualStack: true,
	}

	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:        []string{"kafka:9093"},
		Topic:          "orders",
		GroupID:        "go-app",
		CommitInterval: 0,
		Dialer:         dialer,
	})
	log.Println("Consumer успешно запущен!")

	// Чтение сообщений
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			log.Printf("Ошибка принятия сообщения: %v\n", err)
		} else {
			if IsValidDataFromKafka(msg.Value) {
				err = json.Unmarshal(msg.Value, &order)
				if err != nil {
					log.Printf("Ошибка обработки в струкруру сообщения: %v\n", err)
				} else {
					// Вставляем в бд
					err = database.InsertOrder(ctx, pool, order)
					if err != nil {
						log.Printf("Ошибка вставки полученных данных из kafka в бд: %v\n", err)
					} else {
						log.Printf("Новая запись успешно вставлена в бд! Её order_uid: %s\n", order.Order_uid)
						// Добавление новой записи в кеш
						cache.Add(order.Order_uid, order)
						log.Printf("Новая запись успешно добавлена в кеш! Её order_uid: %s\n", order.Order_uid)
					}
				}
				// Коммитим оффсет вручную после обработки
				err = reader.CommitMessages(ctx, msg)
				if err != nil {
					log.Printf("Ошибка коммита сообщения: %v\n", err)
				}
			} else {
				log.Println("Пришедшие данные из kafka невалидны!")
			}
		}
	}
}
