package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/segmentio/kafka-go"
)

func main() {
	ctx := context.Background()

	// Инициализация producer с настройками адреса брокера и топика
	writer := kafka.NewWriter(kafka.WriterConfig{
		Brokers: []string{"kafka:9093"},
		Topic:   "orders",
	})
	defer writer.Close()
	fmt.Println("Продюсер успешно запущен!")

	// Отправка сообщений
	var (
		order []byte
		err   error
	)
	// Берём данные из файла
	order, err = os.ReadFile("model.json")
	if err != nil {
		fmt.Printf("Ошибка чтения данных из файла model.json: %v\n", err)
		return
	}
	for {
		// Отправка сообщения брокеру
		err := writer.WriteMessages(ctx, kafka.Message{
			Value: order,
		})
		if err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
			return
		}
		time.Sleep(20 * time.Second)
	}
}
