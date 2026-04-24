FROM golang:1.25-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o relayops-api ./cmd/api
RUN go build -o relayops-worker ./cmd/worker

FROM alpine:latest

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/relayops-api .
COPY --from=builder /app/relayops-worker .
COPY --from=builder /app/migrations ./migrations

EXPOSE 8080