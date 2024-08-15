# Start from the Go base image
FROM golang:1.21.5 AS builder

# Adding metadata labels
LABEL maintainer="femsworld@gmail.com"
LABEL version="1.0"
LABEL description="Docker image for Literary Lions backend application"

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum files and download dependencies
COPY backend/go.mod backend/go.sum ./
RUN go mod download

# Copy the entire backend directory
COPY backend/ /app/backend

# Set the working directory to backend/src
WORKDIR /app/backend/src

# Build the Go application
RUN go build -o /main cmd/main.go

# Final stage: Create a small image for runtime
FROM golang:1.21.5

# Copy the built application from the builder stage
COPY --from=builder /main /main

# Copy the database file (You might need to include migration scripts if any)
COPY backend/src/cmd/literary_lions.db /app/literary_lions.db

# Expose the port that your application will run on
EXPOSE 8080

# Set the entry point command to run your Go app
CMD ["/main"]
