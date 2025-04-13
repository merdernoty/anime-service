FROM alpine:latest
RUN apk --no-cache add ca-certificates tzdata postgresql-client curl

WORKDIR /app
COPY out .
CMD ["./out"]