# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY consumer/consumer.go ./consumer/
# Debug: Verify consumer.go exists
RUN ls -l /app/consumer/
RUN go build -o /app/consumer /app/consumer/consumer.go
# Debug: Verify consumer exists
RUN ls -l /app/

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/consumer .
CMD ["./consumer"]