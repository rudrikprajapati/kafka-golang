# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY ../go.mod ../go.sum ./
RUN go mod download
COPY ../producer ./producer/
RUN go build -o producer ./producer/producer.go

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/producer .
EXPOSE 8080
CMD ["./producer"]