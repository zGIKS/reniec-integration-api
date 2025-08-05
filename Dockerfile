# Dockerfile for ACME Backend (Go)
FROM golang:1.21-alpine AS builder
WORKDIR /app
COPY src/main/acme/go.mod src/main/acme/go.sum ./
RUN go mod download
COPY src/main/acme/ ./
RUN go build -o acme-server main.go

# Final image
FROM alpine:latest
WORKDIR /app
COPY --from=builder /app/acme-server ./
RUN mkdir -p resources
COPY src/main/resources/app.properties ./resources/app.properties
EXPOSE 8080
ENV GIN_MODE=release
CMD ["./acme-server"]
