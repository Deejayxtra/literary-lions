# Stage 1: Build Stage
FROM golang:1.21.5 AS builder

# Adding metadata labels
LABEL maintainer="femsworld@gmail.com"
LABEL version="1.0"
LABEL description="Docker image for Literary Lions frontend application"

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire frontend directory
COPY . /app

# Set the working directory to frontend/src
WORKDIR /app/src

# Build the Go application
RUN go build -o /main main.go

# Stage 2: Runtime Stage
FROM golang:1.21.5

# Create necessary directories
RUN mkdir -p /app/src

# Copy the built application from the builder stage
COPY --from=builder /main /main

# Ensure the templates and static directories are copied
COPY --from=builder /app/src/templates /app/src/templates
COPY --from=builder /app/src/static /app/src/static

# Set the working directory for runtime
WORKDIR /app

# Expose the port your application will run on
EXPOSE 8000

# Command to run your Go application
CMD ["/main"]
