# Stage 1: Build Stage
FROM golang:1.21.5 AS builder

# Adding metadata labels
LABEL maintainer="femsworld@gmail.com"
LABEL version="1.0"
LABEL description="Docker image for Literary Lions frontend application"

# Set the working directory inside the container
WORKDIR /frontend

# Copy go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application code
COPY . /frontend

# Set the working directory for runtime
WORKDIR /frontend/src

# Build the Go application
RUN go build -o /main main.go

# Stage 2: Runtime Stage
FROM golang:1.21.5

# Create necessary directories
RUN mkdir -p /src

# Copy the built application from the builder stage
COPY --from=builder /main /main

# Copy templates and static files to the final image
COPY --from=builder /frontend/src/templates ./templates
COPY --from=builder /frontend/src/static ./static

# Expose the port your application will run on
EXPOSE 8000

# Command to run your Go application
CMD ["/main"]
