package internal

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/go-app/database"
	"github.com/google/uuid"
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
		topics      map[string]struct{}
		orderJson   []byte
		orderStruct database.Orders
		err         error
	)

	rand.Seed(time.Now().UnixNano())

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
		RequiredAcks:           -1, // ожидание подтверждения от всех реплик
		MaxAttempts:            10,
		BatchSize:              100,
		WriteTimeout:           10 * time.Second,
		Balancer:               &kafka.RoundRobin{},
		AllowAutoTopicCreation: true,
	}
	defer writer.Close()
	fmt.Println("Producer успешно запущен!")

	// Берём данные из файла
	orderJson, err = os.ReadFile("model.json")
	if err != nil {
		fmt.Printf("Ошибка чтения данных из файла model.json: %v\n", err)
		return
	}

	// Форматируем в структуру для изменения и отправки уникальных сообщений
	err = json.Unmarshal(orderJson, &orderStruct)
	if err != nil {
		fmt.Printf("Ошибка форматирования данных из json в струкруру из файла model.json: %v\n", err)
		return
	}

	// Отправка сообщений брокеру
	for {
		// Создаем рандомные uuid для обеспечения уникальности каждой записи
		orderStruct.Order_uid = uuid.New()
		orderStruct.Payment.Transaction = orderStruct.Order_uid
		orderStruct.Delivery_uid = uuid.New()
		orderStruct.Delivery.Uid = orderStruct.Delivery_uid

		// Для track_number
		// Случайный номер символа от 1 до 14
		randomIndex := rune(rand.Intn(12) + 1)
		// Случайное число для символа английского алфавита от 65 до 122
		randomNumber := rune(rand.Intn(56) + 65)
		// Перевоплощение
		trackNumberRune := []rune(orderStruct.Track_number)
		trackNumberRune[randomIndex] = randomNumber
		orderStruct.Track_number = string(trackNumberRune)

		for index := range orderStruct.Items {
			orderStruct.Items[index].Rid = uuid.New()
			orderStruct.Items[index].Track_number = orderStruct.Track_number
		}
		orderJson, err = json.Marshal(orderStruct)
		if err != nil {
			fmt.Printf("Ошибка форматирования обновленных данных обратно из струкруры в json из файла model.json: %v\n", err)
			return
		}

		err = writer.WriteMessages(ctx, kafka.Message{
			Value: orderJson,
		})

		if err != nil {
			fmt.Printf("Ошибка отправки сообщения: %v\n", err)
		} else {
			fmt.Println("Сообщение успешно отправлено!")
		}

		// Пауза между отправлениями
		time.Sleep(20 * time.Second)
	}
}
