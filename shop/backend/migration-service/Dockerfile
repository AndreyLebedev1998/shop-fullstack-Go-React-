FROM golang:1.26.1-alpine

WORKDIR /migration-microservice

COPY go.mod go.sum ./
RUN go mod download

# ⬇️ Устанавливаем goose
RUN go install github.com/pressly/goose/v3/cmd/goose@latest

COPY . .

# Air запускает hot-reload
CMD ["sh", "-c", "goose -dir ./migrations postgres \"$DB_CONN\" up"]
