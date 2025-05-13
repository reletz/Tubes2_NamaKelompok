# Build stage
FROM golang:1.24-alpine AS builder

# Set working directory
WORKDIR /app

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download all dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN go build -o main .

# Run stage
FROM alpine:latest

# Set working directory
WORKDIR /app

# Copy the binary from builder
COPY --from=builder /app/main .

# Copy data directory
COPY --from=builder /app/data ./data

# Expose port
EXPOSE 8080

# Command to run the executable
CMD ["./main"]