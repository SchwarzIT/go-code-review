# =========================
# Build Stage
# =========================
FROM golang:1.22-alpine3.19 AS builder

# Set the working directory
WORKDIR /app

# Copy go.mod and go.sum to leverage caching
COPY go.mod go.sum ./
RUN go mod download

# Copy the source code
COPY . .

# Build the Go binary statically
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o main ./cmd/coupon_service

# compress the binary to reduce size
RUN apk add --no-cache upx && upx --best --lzma main

# =========================
# Run Stage
# =========================
FROM alpine:3.19

# Create a non-root user and group
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Set the working directory
WORKDIR /app

# Copy the binary from the builder stage
COPY --from=builder /app/main .

# Create and set permissions for the data directory in a single RUN statement
RUN mkdir -p /app/data && \
    echo '[]' > /app/data/coupons.data.json && \
    chown -R appuser:appgroup /app/data /app/main

# Switch to the non-root user
USER appuser

# Environment variables
ENV API_PORT=80

# Expose the application port
EXPOSE 80


# Command to run the application
ENTRYPOINT ["/app/main"]
