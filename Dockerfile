# Используем официальный образ Go
FROM golang:1.24-alpine AS builder

# Устанавливаем рабочую директорию
WORKDIR /app

# Копируем зависимости
COPY go.mod ./
COPY go.sum ./
RUN go mod download

# Копируем всё остальное
COPY . ./

# Собираем приложение
RUN go build -o main ./cmd/main.go

# Финальный образ
FROM alpine:3.20

# Устанавливаем необходимые пакеты (для BusyBox, curl и т.д.)
RUN apk --no-cache add ca-certificates tzdata

# Рабочая директория
WORKDIR /root/

# Копируем бинарник и статические файлы из builder
COPY --from=builder /app/main .
COPY --from=builder /app/static ./static
COPY --from=builder /app/wait-for.sh ./wait-for.sh

# Делаем скрипт исполняемым
RUN chmod +x ./wait-for.sh

# Открываем порт
EXPOSE 8080

# Команда запуска
CMD ["./wait-for.sh", "db:5432", "--", "./main"]
