# Build stage
FROM golang:1.24 AS builder

WORKDIR /app

# Copy go.mod and go.sum files
COPY go.mod go.sum ./

# Download dependencies
RUN go mod download

# Copy the source code
COPY . .

# Build the application using the Makefile
RUN make build-server-debug && cp cmd/server/server /app/server

# Install ca-certificates in the builder stage
RUN apt-get update && apt-get install -y ca-certificates

# Install Delve debugger
RUN go install github.com/go-delve/delve/cmd/dlv@latest

# Final stage
FROM scratch

WORKDIR /app

# Copy Delve debugger
COPY --from=builder /go/bin/dlv /go/bin/dlv

# Copy SSL certificates from builder
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

# Copy the binary from the builder stage
COPY --from=builder /app/server /app/server

# Expose ports
EXPOSE 8080 9090 2345

# Set environment variables
ENV DB_HOST=postgres \
    DB_PORT=5432 \
    DB_USER=postgres \
    DB_PASSWORD=postgres \
    DB_NAME=blog_db \
    DB_SSLMODE=disable

# Run the server with Delve debugger
# Using dlv exec to run the server binary with debugging enabled
ENTRYPOINT ["/go/bin/dlv", "exec", "/app/server", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "--continue"]
CMD ["--", "--db-host", "postgres", \
     "--db-port", "5432", \
     "--db-user", "postgres", \
     "--db-password", "postgres", \
     "--db-name", "blog_db", \
     "--db-sslmode", "disable"]
