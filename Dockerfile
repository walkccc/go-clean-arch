# Build stage
FROM golang:1.20-alpine3.17 AS builder

WORKDIR /app
COPY . .
RUN go build -o main cmd/main.go
RUN apk add curl
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.15.2/migrate.linux-amd64.tar.gz | tar xvz

# Run stage
FROM alpine:3.17
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/migrate ./migrate
COPY app.env .
COPY start.sh .
COPY internal/db/migration ./migration

EXPOSE 8080
CMD ["/app/main"]
ENTRYPOINT ["/app/start.sh"]
