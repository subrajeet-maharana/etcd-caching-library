# File: Dockerfile
FROM golang:1.23-alpine as builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go build -o etcd-caching-library ./cmd/main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/etcd-caching-library .

# Create a directory for logs
RUN mkdir -p /app/logs

EXPOSE 8080

CMD ["./etcd-caching-library"]