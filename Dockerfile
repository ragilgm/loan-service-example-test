# Stage 1: Build the Go application
FROM golang:1.24 AS builder

# Set working directory
WORKDIR /app

# Copy go.mod and go.sum to cache dependencies
COPY go.mod go.sum ./

# Install dependencies
RUN go mod tidy

# Copy all source files into the container
COPY . ./

# Build the Go application binary statically
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o loan-service ./cmd

# Stage 2: Create the runtime image
FROM alpine:latest

# Install required dependencies (ca-certificates for SSL)
RUN apk --no-cache add ca-certificates

# Copy the statically linked binary from the build stage into the runtime container
COPY --from=builder /app/loan-service /loan-service

# Copy the en.json file from the host into the container's root directory
COPY en.json /en.json
COPY .env /.env

# Verify if the file is correctly copied into the container (optional)
RUN ls -l /en.json

# Set the command to run the application
CMD ["/loan-service"]
