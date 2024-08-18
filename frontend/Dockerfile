# Start from the Go base image
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

# Final stage: Create a small image for runtime
FROM golang:1.21.5

# Copy the built application from the builder stage
COPY --from=builder /main /main

# Expose the port that your application will run on (if applicable)
EXPOSE 8000

# Set the entry point command to run your Go app
CMD ["/main"]
