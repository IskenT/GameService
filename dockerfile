# Stage 1: Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

WORKDIR /app/cmd/
RUN go build -o game


# Stage 2:
FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app/cmd/ .

COPY cmd/.env /app/cmd/.env

CMD ["./game"]
