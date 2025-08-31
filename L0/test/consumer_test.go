package test

import (
	"encoding/json"
	"fmt"
	"os"
	"testing"

	"github.com/ChursinAlexUnder/wbtech-golang-course/L0/internal/database"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
)

type testInputData struct {
	order    []byte
	expected bool
}

func TestIsValidDataFromKafka(t *testing.T) {
	var (
		orderJson   []byte
		orderStruct database.Orders
		err         error
		testTable   []testInputData = []testInputData{}
	)

	// Первый тестовый пример (несоответствие полей ограничению бд)
	orderJson, err = os.ReadFile("../model.json")
	if err != nil {
		fmt.Printf("Ошибка чтения данных из файла model.json при unit тестировании: %v\n", err)
		return
	}
	err = json.Unmarshal(orderJson, &orderStruct)
	if err != nil {
		fmt.Printf("Ошибка форматирования данных из json в струкруру из файла model.json при unit тестировании: %v\n", err)
		return
	}
	orderStruct.Locale = "12345678910"
	orderStruct.Payment.Currency = "12345678910"
	orderStruct.Track_number = ""
	orderStruct.Delivery_uid = uuid.Nil
	orderJson, err = json.Marshal(orderStruct)
	if err != nil {
		fmt.Printf("Ошибка форматирования данных из струкруры в json при unit тестировании: %v\n", err)
		return
	}
	testTable = append(testTable, testInputData{orderJson, false})

	// Второй тестовый пример (json корректен)
	orderJson, err = os.ReadFile("../model.json")
	if err != nil {
		fmt.Printf("Ошибка чтения данных из файла model.json при unit тестировании: %v\n", err)
		return
	}
	testTable = append(testTable, testInputData{orderJson, true})

	for index, test := range testTable {
		result := database.IsValidDataFromKafka(test.order)

		t.Logf("Вызван тест №%d с результатом %t\n", index+1, result)

		assert.Equal(t, test.expected, result,
			fmt.Sprintf("Неправильный результат! Ожидалось %t, результат %t", test.expected, result))
	}

}
