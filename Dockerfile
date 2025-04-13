FROM golang:1.24-alpine
WORKDIR /app
COPY . .
# Добавьте вывод для отладки
RUN ls -la
RUN go version
RUN go mod tidy
# Используйте флаг -v для подробного вывода
RUN go build -v -o out ./cmd/api
CMD ["./out"]
