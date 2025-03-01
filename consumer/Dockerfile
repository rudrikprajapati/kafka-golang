# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY ../go.mod ../go.sum ./
RUN go mod download
COPY ../consumer ./consumer/
RUN go build -o consumer ./consumer/consumer.go

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/consumer .
CMD ["./consumer"]