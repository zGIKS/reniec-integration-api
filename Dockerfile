# Dockerfile for ACME Backend (Go) - Render Compatible
FROM golang:1.21-alpine AS builder

# Install ca-certificates for SSL connections
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Copy Go modules
COPY src/main/acme/go.mod src/main/acme/go.sum ./
RUN go mod download

# Copy source code
COPY src/main/acme/ ./

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o acme-server main.go

# Final image
FROM alpine:latest

# Install ca-certificates for SSL connections to Supabase
RUN apk --no-cache add ca-certificates tzdata
RUN mkdir /app

WORKDIR /app

# Copy binary from builder
COPY --from=builder /app/acme-server ./

# Create resources directory and copy properties
RUN mkdir -p resources
COPY src/main/resources/app.properties ./resources/app.properties

# Copy .env file for environment configuration
COPY .env ./

# Expose port (will be overridden by Render's PORT env var)
EXPOSE 8080

CMD ["./acme-server"]
