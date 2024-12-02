# Используем образ golang для сборки
FROM golang:1.22.9 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

# Собираем приложение
RUN go build -o /app/games-bot ./.

# Запускаем приложение при запуске контейнера
CMD ["./games-bot"]

EXPOSE 8080