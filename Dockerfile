# Используем официальный Go-образ для сборки приложения
FROM golang:1.23-alpine as builder

# Устанавливаем рабочую директорию в контейнере
WORKDIR /app

# Копируем файлы проекта в контейнер
COPY . .

# Загружаем зависимости и компилируем приложение
RUN go mod tidy
RUN go build -o main ./cmd/main.go

# Используем минимальный образ для финального контейнера
FROM alpine:latest

# Устанавливаем необходимые зависимости
RUN apk add --no-cache ca-certificates

# Копируем скомпилированный файл из первого контейнера
COPY --from=builder /app/main /app/

# Копируем миграции в контейнер
COPY migration /app/migration

# Копируем .env файл
COPY .env /app/

# Указываем рабочую директорию
WORKDIR /app

# Открываем порт для приложения
EXPOSE 8080

# Запускаем приложение
CMD ["./main"]
