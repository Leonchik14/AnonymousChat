# Используем официальный образ Golang для сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем необходимые зависимости
RUN apk add --no-cache git

# Создаём рабочую директорию
WORKDIR /app

# Копируем файлы зависимостей и скачиваем их
COPY go.mod go.sum ./
RUN go mod download

# Копируем весь код
COPY . .
COPY .env .env

# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /chat-service ./cmd

# ------------------------------------------------------------
# Финальный образ
FROM alpine:latest

# Устанавливаем сертификаты
RUN apk add --no-cache ca-certificates

# Копируем бинарник
COPY --from=builder /chat-service /usr/local/bin/chat-service

WORKDIR /app

# Экспонируем порт
EXPOSE 8082
EXPOSE 50051

# Запуск сервиса
ENTRYPOINT [ "chat-service" ]
