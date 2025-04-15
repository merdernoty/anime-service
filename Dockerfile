FROM golang:1.24-alpine
WORKDIR /app
COPY . .
RUN ls -la
RUN go version
RUN go mod tidy
RUN go build -v -o out ./cmd/api
CMD ["./out"]