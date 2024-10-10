# Use an official Golang image as a builder
FROM golang:1.19-alpine AS builder

# Set the working directory inside the container
WORKDIR /app

# Copy the go.mod and go.sum files and download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code into the container
COPY . .

# Build the Go binary
# RUN go build -o multiplayer-service ./cmd/main.go
# Use a lightweight image for the final container
FROM alpine:latest

# Set the working directory inside the container
WORKDIR /app

# Copy the built binary from the builder stage
COPY --from=builder /app/multiplayer-service .

# Expose the gRPC port
EXPOSE 50051

# Set environment variables if required (optional)
ENV MONGO_URI="mongodb+srv://yash:yashMongodb@ecommerce-backend.eyunzek.mongodb.net/?retryWrites=true&w=majority&appName=Ecommerce-Backend"

# Run the Go service binary
CMD ["./multiplayer-service"]