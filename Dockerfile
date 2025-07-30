# Stage 1: Build the Go application
FROM golang:1.24-alpine AS builder

WORKDIR /app

# Copy go.mod and go.sum first to leverage Docker cache
COPY go.mod go.sum ./
RUN go mod download

# Copy the rest of the application source code
COPY . .

# Build the Go application statically
# CGO_ENABLED=0 disables Cgo, GOOS=linux sets the target OS, GOARCH=amd64 sets the target architecture
# -ldflags="-w -s" strips debug information to reduce binary size
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o discord-gemini-bot main.go

# Stage 2: Create a minimal image with the compiled binary
FROM alpine:latest

WORKDIR /root/

# Copy the compiled binary from the builder stage
COPY --from=builder /app/discord-gemini-bot .

# Expose any necessary ports (though Discord bots don't typically expose ports for incoming connections)
# EXPOSE 8080

# Command to run the executable
CMD ["./discord-gemini-bot"]
