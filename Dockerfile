# Use AWS ECR Public Go base image
FROM public.ecr.aws/docker/library/golang:1.24-alpine AS builder

# Install git for go mod download
RUN apk add --no-cache git

# Set working directory
WORKDIR /app

# Copy go mod files
COPY go.mod go.sum ./

# Set Go proxy to direct to bypass proxy issues in AWS
RUN go env -w GOPROXY=direct

# Download dependencies
RUN go mod download

# Copy source code
COPY . .

# Build the server binary
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server ./cmd/server

# Final stage - use minimal AWS ECR base image
FROM public.ecr.aws/docker/library/alpine:latest

# Install ca-certificates for HTTPS requests
RUN apk --no-cache add ca-certificates

# Create non-root user
RUN adduser -D -s /bin/sh tetris

# Set working directory
WORKDIR /app

# Copy binary from builder stage
COPY --from=builder /app/server .

# Change ownership to non-root user
RUN chown tetris:tetris /app/server

# Switch to non-root user
USER tetris

# Expose port
EXPOSE 8080

# Health check
HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
  CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

# Run the server (can be overridden with docker run command)
CMD ["./server"]
