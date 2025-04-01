FROM golang:1.24-alpine AS builder

WORKDIR /app
COPY . .
RUN go build -o grpc-auth ./cmd

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/grpc-auth .
COPY --from=builder /app/migrations ./migrations

EXPOSE 50051
CMD ["./grpc-auth"]