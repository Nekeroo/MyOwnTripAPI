# syntax=docker/dockerfile:1

FROM golang:1.25.8-trixie AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o /my-own-trip ./cmd

FROM debian:trixie-slim

COPY --from=builder /my-own-trip /my-own-trip

EXPOSE 8080

CMD ["/my-own-trip"]