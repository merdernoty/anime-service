
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем go mod и sum файлы
COPY go.mod go.sum ./

# Загружаем зависимости
RUN go mod download

# Копируем исходный код
COPY . .

# Собираем приложение
RUN go build -o out ./cmd/api

# Финальный образ
FROM alpine:latest

WORKDIR /root/

# Копируем собранное приложение из builder
COPY --from=builder /app/out .

# Если у вас есть какие-то статические файлы или конфигурации
# COPY --from=builder /app/some-config-or-static-files ./

# Указываем порт, который будет использовать приложение
EXPOSE 8080

# Команда запуска приложения
CMD ["./out"]
