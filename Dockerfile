# syntax=docker/dockerfile:1

FROM golang:1.25.8-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /my-own-trip ./cmd

FROM debian:trixie-slim

RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

COPY --from=builder /my-own-trip /my-own-trip

EXPOSE 3000

CMD ["/my-own-trip"]