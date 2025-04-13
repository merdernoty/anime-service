# --- Stage 1: Build ---
    FROM golang:1.21-alpine AS builder

    WORKDIR /app
    
    COPY . .
    
    RUN go build -o out .
    
    # --- Stage 2: Runtime ---
    FROM alpine:latest
    
    WORKDIR /app
    
    # Устанавливаем зависимости (если нужно)
    RUN apk --no-cache add ca-certificates tzdata postgresql-client curl
    
    # Копируем только бинарник
    COPY --from=builder /app/out .
    
    CMD ["./out"]
    