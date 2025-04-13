FROM golang:1.24-alpine
WORKDIR /app
# Сначала копируем только go.mod и go.sum
COPY go.mod go.sum ./
# Затем download dependencies
RUN go mod download
# Затем копируем остальной код
COPY . .
# Компиляция
RUN go build -o out ./cmd/api
CMD ["./out"]
