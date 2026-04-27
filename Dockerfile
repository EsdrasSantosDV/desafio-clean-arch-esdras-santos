FROM golang:1.22-alpine AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o ordersystem ./cmd/ordersystem

FROM alpine:3.20

WORKDIR /app

COPY --from=builder /app/ordersystem /app/ordersystem
COPY migrations /app/migrations

EXPOSE 8000 50051 8080

CMD ["/app/ordersystem"]
