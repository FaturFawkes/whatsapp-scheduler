FROM golang:1.21-alpine AS builder
WORKDIR /app

COPY . .

RUN go build -o scheduler .

CMD ["./scheduler"]