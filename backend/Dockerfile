# Build stage
FROM golang:1.24.1-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Set working directory
WORKDIR /build

# Copy go mod and sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 GOOS=linux go build -o /build/syncer ./api/main.go

# Final stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create non-root user
RUN adduser -D -g '' appuser

# Set working directory
WORKDIR /app

# Create source directory for volume mount
RUN mkdir /app/src && chown -R appuser:appuser /app

# Copy the binary from builder
COPY --from=builder /build/syncer /app/syncer

# Copy migrations
COPY --from=builder /build/db/migrations /app/db/migrations

# Set ownership to non-root user
RUN chown -R appuser:appuser /app

# Use non-root user
USER appuser

# Expose the application port
EXPOSE 8080

# Set environment variables
ENV GIN_MODE=debug

# Run the application
CMD ["/app/syncer"]