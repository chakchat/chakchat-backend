# Stage 1: Build
FROM golang:1.23.1 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files for dependency installation
COPY go.mod go.sum ./

# Download and cache dependencies
RUN go mod download

# Copy the rest of the application code
COPY . .

# Build the application binary
RUN go build -o main .

# Stage 2: Run
FROM ubuntu:latest

# https://chatgpt.com/share/67851f69-d6e8-8011-b3fb-839ce5bf275c
RUN apt-get update && apt-get install -y ca-certificates && rm -rf /var/lib/apt/lists/*

# Set the working directory inside the container
WORKDIR /app

# Copy the compiled binary from the builder stage
COPY --from=builder /app/main .

# Expose the port your application runs on (optional)
EXPOSE 5004

# Command to run the application
CMD ["./main"]
