FROM golang:alpine AS builder

WORKDIR /app

# Set env variables for Go
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the application
RUN go build -o main ./cmd/api

# Final minimal image
FROM alpine:latest

WORKDIR /app

# Copy the binary from the builder
COPY --from=builder /app/main .

# Copy configuration folder
COPY --from=builder /app/config ./config

EXPOSE 8080

CMD ["./main"]
