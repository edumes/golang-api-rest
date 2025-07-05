# Build stage
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o golang-api-rest ./cmd/api

# Run stage
FROM alpine:latest
WORKDIR /root/
COPY --from=builder /app/golang-api-rest .
COPY .env.example .env
EXPOSE 8080
CMD ["./golang-api-rest"] 