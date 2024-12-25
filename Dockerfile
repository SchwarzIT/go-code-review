# Build stage
FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/coupon_service

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .


EXPOSE 8080 9090
CMD [ "/app/main" ]

