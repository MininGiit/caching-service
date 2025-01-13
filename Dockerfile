# Указываем базовый образ
FROM golang:1.22.2 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download
COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -o app ./cmd/app

FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/app .

ENV SERVER_HOST_PORT=":8080"
ENV CACHE_SIZE="100"
ENV DEFAULT_CACHE_TTL="60"
ENV LOG_LEVEL="INFO"

EXPOSE 8080
CMD ["./app"]
