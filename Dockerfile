# Build stage
FROM golang:1.21-alpine AS builder

# Install build dependencies
RUN apk add --no-cache gcc musl-dev pkgconfig gtk+3.0-dev webkit2gtk-dev

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=1 go build -ldflags="-s -w" -o singbox-xboard-client cmd/main.go

# Runtime stage
FROM alpine:latest

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Create app user
RUN addgroup -g 1001 -S app && \
    adduser -u 1001 -S app -G app

# Set working directory
WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/singbox-xboard-client .

# Create config directory
RUN mkdir -p /app/config && chown -R app:app /app

# Switch to app user
USER app

# Expose port
EXPOSE 7890 9090

# Set default command
CMD ["./singbox-xboard-client"]