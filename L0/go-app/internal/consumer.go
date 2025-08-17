package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	"unicode/utf8"

	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/go-app/database"
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
		utf8.RuneCountInString(orderStruct.Delivery.Email) > 254 {
		return false
	}
	for _, item := range orderStruct.Items {
		if utf8.RuneCountInString(item.Size) > 10 ||
			utf8.RuneCountInString(item.Brand) > 150 {
			return false
		}
	}
	return true
}

func Consumer(ctx context.Context, pool *pgxpool.Pool) {
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
	fmt.Println("Consumer успешно запущен!")

	// Чтение сообщений
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			fmt.Printf("Ошибка принятия сообщения: %v\n", err)
		} else {
			if IsValidDataFromKafka(msg.Value) {
				err = json.Unmarshal(msg.Value, &order)
				if err != nil {
					fmt.Printf("Ошибка обработки в струкруру сообщения: %v\n", err)
				} else {
					// Вставляем в бд
					err = database.InsertOrder(ctx, pool, order)
					if err != nil {
						fmt.Printf("Ошибка вставки полученных данных из kafka в бд: %v\n", err)
					} else {
						fmt.Printf("Новая запись успешно вставлена! Её order_uid: %s\n", order.Order_uid)
					}
				}
				// Коммитим оффсет вручную после обработки
				err = reader.CommitMessages(ctx, msg)
				if err != nil {
					fmt.Printf("Ошибка коммита сообщения: %v\n", err)
				}
			} else {
				fmt.Println("Пришедшие данные из kafka невалидны!")
			}
		}
	}
}
