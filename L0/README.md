Как запустить

В папке L0 команда docker compose up --build

Основные эндпоинты

Golang приложение

    localhost:8081/order - главная страница сайта
    localhost:8081/order/<order_uid> - информация о заказе в web интерфейсе
    localhost:8081/api/<order_uid> - json заказа с сервера

Kafka-ui

    localhost:8080/ - главная страница

Swagger

    localhost:8081/swagger/index.html - главная страница

