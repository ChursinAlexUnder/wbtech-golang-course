# L0
## Описание работы
Producer-скрипт при запуске сервиса берет данные из json файла и раз в 20 секунд отправляет их брокеру. 
Запущенный consumer, подключенный к нужному topic-у считывает их и если данные валидны, то вставляет в бд и добавляет в кеш. 
Получение данных на странице реализуется из кеша, если там данных нет, то из бд (кеш заполняется при запуске сервиса).
## Как запустить:
В каталоге **L0** команда:
```
docker compose up --build
```
## Схема работы сервисов
![Схема работы сервисов](/L0/assets/Схема%20работы.png)
## Эндпоинты
### Golang приложение
- [localhost:8081/order](url) - главная страница сайта
- [localhost:8081/order/<order_uid>](url) - информация о заказе в web интерфейсе
- [localhost:8081/api/<order_uid>](url) - json заказа с сервера
### Kafka-ui
- [localhost:8080/](url) - главная страница
### Swagger
- [localhost:8081/swagger/index.html](url) - главная страница с доп. информацией о проекте (методы, модели данных и т. д.)
## Структура проекта L0
```
L0/
├── api/
│   └── model.json - файл с json, который producer считывает для отправки в kafka
├── assets/
│   └── favicon.ico - иконка для вкладки в браузере
├── build/
│   └── package/
│       └── Dockerfile - docker файл для создания образа go-сервиса
├── cmd/
│   └── go-app/
│       └── main.go - точка входа приложения
├── deployments/
│   └── docker-compose.yml - compose файл для создания и запуска всех образов
├── docs/ - папка для работы swagger (создается автоматически с помощью библиотеки)
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
├── internal/
│   └── controller/
│       └── controller.go - файл с настройкой эндпоинтов и их описанием для swagger
├── database/
│   ├── connection.go - подключение к базе данных
│   ├── inserter.go - вставка данных из kafka в бд
│   ├── selecter.go - взятие данных из бд
│   └── structure.go - структуры для взаимодействия с данными
├── router/
│   └── router.go - файл с вызовом эндпоинтов
├── consumer.go - вся логика consumer-а + валидация данных от брокера
├── producer.go вся логика producer-а
├── schema/ - миграции
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
├── test/
│   └── consumer_test.go - unit-тест
├── web/
│   ├── index.html - страница сайта
│   ├── script.js - вся логика клиента
│   └── style.css - стили для страницы
├── .gitignore
├── go.mod - файл зависимости
└── go.sum - файл зависимости
```
---
