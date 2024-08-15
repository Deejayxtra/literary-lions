# Start from the Go base image
FROM golang:1.21.5 AS builder

# Adding metadata labels
LABEL maintainer="femsworld@gmail.com"
LABEL version="1.0"
LABEL description="Docker image for Literary Lions frontend application"

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY frontend/go.mod frontend/go.sum ./
RUN go mod download

# Copy the entire frontend directory
COPY frontend/ /app/frontend

# Set the working directory to frontend/src
WORKDIR /app/frontend/src

# Build the Go application
RUN go build -o /frontend-main main.go

# Final stage: Create a small image for runtime
FROM golang:1.21.5

# Copy the built application from the builder stage
COPY --from=builder /frontend-main /frontend-main

# Expose the port that your application will run on (if applicable)
EXPOSE 8000

# Set the entry point command to run your Go app
CMD ["/frontend-main"]
