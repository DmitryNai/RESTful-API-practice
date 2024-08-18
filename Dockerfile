FROM golang:1.23-alpine

RUN apk update && apk add --no-cache \
    openssl-dev \
    gcc \
    musl-dev

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o kv-api .

EXPOSE 8080

CMD ["./kv-api"]