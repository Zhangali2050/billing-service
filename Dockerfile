# Этап сборки
FROM golang:1.21-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o billing-service ./cmd/main.go

# Финальный минимальный образ
FROM alpine:3.18

WORKDIR /app

# Копируем только бинарник
COPY --from=builder /app/billing-service .

# Добавим tzdata и сертификаты для HTTPS-запросов
RUN apk --no-cache add ca-certificates tzdata

# Убедимся, что tzdata установлен и выставлена зона (опционально)
ENV TZ=Asia/Almaty

EXPOSE 8080

CMD ["./billing-service"]
