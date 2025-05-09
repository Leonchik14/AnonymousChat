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
COPY .env /app/.env


# Собираем бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /matchmaking-service ./cmd

# ------------------------------------------------------------
# Финальный образ
FROM alpine:latest

# Устанавливаем сертификаты
RUN apk add --no-cache ca-certificates

# Копируем бинарник
COPY --from=builder /matchmaking-service /usr/local/bin/matchmaking-service

WORKDIR /app

# Экспонируем порт
EXPOSE 8081

# Запуск сервиса
ENTRYPOINT [ "matchmaking-service" ]
