FROM golang:1.25-alpine

WORKDIR /shop-microservice

COPY go.mod go.sum ./
RUN go mod download

# ⬇️ Устанавливаем air
RUN go install github.com/air-verse/air@latest

# ⬇️ Устанавливаем goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

# Air запускает hot-reload
CMD ["sh", "-c", "goose -dir ./migrations postgres \"$DB_CONN\" up && air -c .air.toml"]
