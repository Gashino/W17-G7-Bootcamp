# Etapa de build
FROM golang:1.21-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app ./cmd

# Etapa final
FROM alpine:3.19

WORKDIR /app

RUN apk add --no-cache bash

COPY --from=builder /app/app .
COPY wait-for-it.sh /wait-for-it.sh
RUN chmod +x /wait-for-it.sh

EXPOSE 8080

ENV DB_HOST=mysql
ENV DB_PORT=3306
ENV DB_USER=root
ENV DB_PASSWORD=""
ENV DB_NAME=frescos

ENTRYPOINT ["/wait-for-it.sh", "mysql:3306", "--", "./app"]