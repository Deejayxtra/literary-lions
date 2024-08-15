# Start from the Go base image for building
FROM golang:1.21.5 AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the entire backend directory into the container
COPY . /app

# Set the working directory to backend/src
WORKDIR /app/src

# Build the Go application
RUN go build -o /main cmd/main.go

# Final stage: Create a small image for runtime
FROM golang:1.21.5

# Install sqlite3 if you need it for DB operations
RUN apt-get update && apt-get install -y sqlite3

# Create necessary directories
RUN mkdir -p /app/src

# Copy the built application from the builder stage
COPY --from=builder /main /main

# Copy the database file to the correct location
COPY --from=builder /app/src/literary_lions.db /app/src/literary_lions.db

# Expose the port that your application will run on
EXPOSE 8080

# Set the entry point command to run your Go app
CMD ["/main"]
