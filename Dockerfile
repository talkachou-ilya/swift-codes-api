FROM golang:1.24-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o swift-codes-api ./cmd/main.go

EXPOSE 8080

RUN go test ./tests/... -v

CMD ["./swift-codes-api"]
