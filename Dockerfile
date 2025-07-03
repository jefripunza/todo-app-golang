FROM golang:alpine AS builder
LABEL version="1.0"
LABEL description="Todo List API"
LABEL author="Jefri Herdi Triyanto"
LABEL email="hi@jefripunza.com"

WORKDIR /app
COPY go.* ./
RUN go mod download
COPY . .
RUN go build -o app main.go




FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .
CMD ["./app"]
