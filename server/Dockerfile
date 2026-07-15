# Шаг 1: Сборка бинарника (остается как был)
FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/server-bin ./server

# Шаг 2: Финальный образ (ОБНОВЛЯЕМ)
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/server-bin .

# ВАЖНО: Копируем папку с HTML-шаблонами в финальный контейнер
COPY --from=builder /app/internal/network/web/templates ./internal/network/web/templates

EXPOSE 80
CMD ["./server-bin"]
