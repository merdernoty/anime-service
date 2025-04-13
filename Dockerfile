FROM golang:1.24-alpine
WORKDIR /app

COPY go.mod go.sum ./

RUN go mod download

COPY . .

RUN go build -ldflags="-s -w" -o out ./cmd/api
CMD ["./out"]
