# wb_test_0

Инициализация и запуск postgres, запуск nats-streaming сервера
```
docker-compose up -d
```

Запуск программы для запуска http-сервера для работой с заказами и подписки на события nats-streaming сервера
```
go run cmd/cons/main.go
```

Запуск программы для генерации событий в nats-streaming сервер
```
go run cmd/pub/main.go
```

Запуск фронтенда
```
npm run frontend/ start
```

Запуск тестов
```
go test ./...
```
