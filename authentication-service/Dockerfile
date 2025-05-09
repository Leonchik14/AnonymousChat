# Используем официальный образ Golang для сборки
FROM golang:1.24-alpine AS builder

# Устанавливаем необходимые зависимости
RUN apk add --no-cache git

# Создаём рабочую директорию в контейнере
WORKDIR /app

# Копируем файлы с описанием модулей и скачиваем зависимости
COPY go.mod go.sum ./
RUN go mod download

# Копируем всё приложение
COPY . .

# Собираем статически скомпилированный бинарник
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -o /auth-service ./cmd

# ------------------------------------------------------------
# Финальный лёгкий образ
FROM alpine:latest

# Для HTTPS-запросов
RUN apk add --no-cache ca-certificates

# Копируем собранный бинарник
COPY --from=builder /auth-service /usr/local/bin/auth-service

WORKDIR /app

# Экспонируем порт
EXPOSE 8083

# Запуск сервиса
ENTRYPOINT [ "auth-service" ]
