package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/segmentio/kafka-go"
)

func Consumer(ctx context.Context) {
	var order Orders

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
			err = json.Unmarshal(msg.Value, &order)
			if err != nil {
				fmt.Printf("Ошибка обработки в струкруру сообщения: %v\n", err)
			} else {
				fmt.Println(order.Order_uid)
			}
			// Коммитим оффсет вручную после обработки
			err = reader.CommitMessages(ctx, msg)
			if err != nil {
				fmt.Printf("Ошибка коммита сообщения: %v\n", err)
			}
		}
	}
}
