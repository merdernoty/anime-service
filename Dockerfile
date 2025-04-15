# Используйте официальный образ Go
FROM golang:1.24-alpine

# Установите рабочую директорию
WORKDIR /app

# Скопируйте go mod и sum файлы
COPY go.mod go.sum ./

# Скачайте зависимости
RUN go mod download

# Скопируйте исходный код
COPY . .

# Соберите приложение
RUN go build -v -o main ./cmd/api/main.go

# Откройте порт
EXPOSE 8080

# Запустите приложение
CMD ["./main"]