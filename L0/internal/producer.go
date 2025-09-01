package internal

import (
	"context"
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/internal/database"
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
	// Инициализация writer
	writer := &kafka.Writer{
		Addr:                   kafka.TCP("kafka:9093"),
		Topic:                  topic,
		RequiredAcks:           -1,
		MaxAttempts:            10,
		BatchSize:              100,
		WriteTimeout:           10 * time.Second,
		Balancer:               &kafka.RoundRobin{},
		AllowAutoTopicCreation: true,
	}
	defer func() {
		if err := writer.Close(); err != nil {
			log.Printf("Ошибка закрытия writer: %v", err)
		}
	}()

	log.Println("Producer успешно запущен!")

	// Загружаем шаблон один раз
	orderJson, err := os.ReadFile("./api/model.json")
	if err != nil {
		log.Printf("Ошибка чтения данных из файла model.json: %v\n", err)
		return
	}
	var orderStruct database.Orders
	if err := json.Unmarshal(orderJson, &orderStruct); err != nil {
		log.Printf("Ошибка форматирования данных из json в структуру: %v\n", err)
		return
	}

	rand.Seed(time.Now().UnixNano())

	for {
		// Проверяем отмену контекста перед каждой итерацией
		select {
		case <-ctx.Done():
			log.Println("Producer: контекст отменён, завершение работы")
			return
		default:
		}

		msgStruct := orderStruct
		msgStruct.Order_uid = uuid.New()
		msgStruct.Payment.Transaction = msgStruct.Order_uid
		msgStruct.Delivery_uid = uuid.New()
		msgStruct.Delivery.Uid = msgStruct.Delivery_uid

		// Защитная проверка длины Track_number
		if len(msgStruct.Track_number) > 2 {
			randomIndex := rand.Intn(len(msgStruct.Track_number))
			randomNumber := rune(rand.Intn(26) + 65)
			trackNumberRune := []rune(msgStruct.Track_number)
			trackNumberRune[randomIndex] = randomNumber
			msgStruct.Track_number = string(trackNumberRune)
		}

		for i := range msgStruct.Items {
			msgStruct.Items[i].Rid = uuid.New()
			msgStruct.Items[i].Track_number = msgStruct.Track_number
		}

		data, err := json.Marshal(msgStruct)
		if err != nil {
			log.Printf("Ошибка маршалинга сообщения: %v\n", err)
			goto sleepPeriod
		}

		if err := writer.WriteMessages(ctx, kafka.Message{Value: data}); err != nil {
			if errors.Is(err, context.Canceled) || ctx.Err() != nil {
				log.Println("Producer: запись прервана контекстом, завершение работы")
				return
			}
			log.Printf("Ошибка отправки сообщения: %v\n", err)
		} else {
			log.Println("Сообщение успешно отправлено!")
		}

	sleepPeriod:
		select {
		case <-ctx.Done():
			log.Println("Producer: контекст отменён, завершение работы")
			return
		case <-time.After(20 * time.Second):
		}
	}
}
