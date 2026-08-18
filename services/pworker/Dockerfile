# Шаг 1: Сборка бинарника
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Собираем воркер
RUN go build -o /app/worker-bin .

# Шаг 2: Финальный образ
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/worker-bin .
CMD ["./worker-bin"]
