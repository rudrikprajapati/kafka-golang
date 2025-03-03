# Build stage
FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY producer/producer.go ./producer/
# Debug: Verify producer.go exists
RUN ls -l /app/producer/
RUN go build -v -o /app/producer /app/producer/producer.go
# Debug: Verify producer exists
RUN ls -l /app/

# Run stage
FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /app/producer .
EXPOSE 8080
CMD ["./producer"]