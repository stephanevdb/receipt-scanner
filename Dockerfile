FROM golang:1.21-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache git

# Copy go mod and sum files
COPY go.mod ./

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -o pocketbase

# Create final image
FROM alpine:latest

# Install required dependencies
RUN apk add --no-cache ca-certificates

# Create a non-root user
RUN adduser -D pocketbase

# Set working directory
WORKDIR /pb

# Copy the binary from builder
COPY --from=builder /app/pocketbase /pb/pocketbase
COPY --from=builder /app/pb_migrations /pb/pb_migrations

# Set proper permissions
RUN chown -R pocketbase:pocketbase /pb

# Switch to non-root user
USER pocketbase

# Expose the default PocketBase port
EXPOSE 8080

# Start PocketBase
CMD ["/pb/pocketbase", "serve", "--http=0.0.0.0:8080"] 