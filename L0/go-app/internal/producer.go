package internal

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

// Функция для создания topic по переданным параметрам
func createCustomTopic(topic string, partitions, replicationFactor int) error {
	conn, err := kafka.Dial("tcp", "kafka:9093")
	if err != nil {
		return err
	}
	defer conn.Close()

	controller, err := conn.Controller()
	if err != nil {
		return err
	}
	var controllerConn *kafka.Conn
	controllerConn, err = kafka.Dial("tcp", net.JoinHostPort(controller.Host, strconv.Itoa(controller.Port)))
	if err != nil {
		return err
	}
	defer controllerConn.Close()

	topicConfigs := []kafka.TopicConfig{
		{
			Topic:             topic,
			NumPartitions:     partitions,
			ReplicationFactor: replicationFactor,
		},
	}

	err = controllerConn.CreateTopics(topicConfigs...)
	if err != nil {
		return err
	}
	return nil
}

// Функция для получения списка имеющихся topic-ов
func takeListTopics() (map[string]struct{}, error) {
	conn, err := kafka.Dial("tcp", "kafka:9093")
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	partitions, err := conn.ReadPartitions()
	if err != nil {
		return nil, err
	}

	m := map[string]struct{}{}

	for _, p := range partitions {
		m[p.Topic] = struct{}{}
	}
	return m, nil
}

// Основная функция работы producer
func Producer(ctx context.Context, topic string, partitions, replicationFactor int) {
	var (
		topics map[string]struct{}
		order  []byte
		err    error
	)

	// Проверка на наличие topic с нужным именем. Если такового нет, то создаем
	for haveTopic := false; !haveTopic; {
		topics, err = takeListTopics()
		if err != nil {
			fmt.Printf("Ошибка чтения списка topic-ов: %v\n", err)
		} else {
			if _, ok := topics[topic]; ok {
				haveTopic = true
			} else {
				// Создание кастомного topic
				err = createCustomTopic(topic, partitions, replicationFactor)
				if err != nil {
					fmt.Printf("Ошибка добавления нового topic: %v\n", err)
				}
			}
		}
	}

	// Инициализация producer с настройками адреса брокера и имени топика
	writer := &kafka.Writer{
		Addr:                   kafka.TCP("kafka:9093"),
		Topic:                  "orders",
		AllowAutoTopicCreation: true,
	}
	defer writer.Close()
	fmt.Println("Producer успешно запущен!")

	// Берём данные из файла
	order, err = os.ReadFile("model.json")
	if err != nil {
		fmt.Printf("Ошибка чтения данных из файла model.json: %v\n", err)
		return
	}

	// Отправка сообщений брокеру
	for {
		err = writer.WriteMessages(ctx, kafka.Message{
			Value: order,
		})

		if err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		} else {
			fmt.Println("Сообщение успешно отправлено!")
		}

		// Пауза между отправлениями
		time.Sleep(10 * time.Second)
	}
}
