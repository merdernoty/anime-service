FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
# Проверим, какие файлы есть в репозитории
RUN find . -type f | grep "\.go$" || echo "No Go files found"
# Попробуем скомпилировать из директории cmd/api
RUN if [ -d cmd/api ]; then cd cmd/api && go build -o /app/out; else echo "cmd/api directory not found"; fi

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata
COPY --from=builder /app/out .
CMD ["./out"]