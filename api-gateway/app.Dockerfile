# Use the official Golang image as the base image
FROM golang:1.24 AS builder

# Set environment variables
ENV GO111MODULE=on \
    CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum to download dependencies
COPY go.mod go.sum ./
RUN go mod tidy && go mod download

# Copy the rest of the application code
COPY . .

# Build the Go application
RUN go build -o main .

# Use a lightweight image for the final container
FROM alpine:latest

# Set working directory
WORKDIR /root/

# Install necessary dependencies
RUN apk add --no-cache ca-certificates

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Copy .env file (Optional: Use only if you are not using docker-compose env_file)
# COPY .env .env

# Expose the port that your app will run on
EXPOSE 8000

# Command to run the application
CMD ["./main"]
