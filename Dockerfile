FROM golang:1.24-alpine AS builder
WORKDIR /app
COPY . .
RUN go build -o out .

FROM alpine:latest
WORKDIR /app
RUN apk --no-cache add ca-certificates tzdata postgresql-client curl
COPY --from=builder /app/out .
CMD ["./out"]