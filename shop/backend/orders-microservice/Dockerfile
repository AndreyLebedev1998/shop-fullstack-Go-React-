FROM golang:1.26.1-alpine

WORKDIR /orders-microservice

COPY go.mod go.sum ./
RUN go mod download

# ⬇️ Устанавливаем air
RUN go install github.com/air-verse/air@latest

COPY . .

# Air запускает hot-reload
CMD ["sh", "-c", "air -c .air.toml"]
