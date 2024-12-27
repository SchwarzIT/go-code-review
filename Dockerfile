# Build stage
FROM golang:1.22-alpine3.19 AS builder
WORKDIR /app
COPY . .
RUN go build -o main ./cmd/coupon_service

# Run stage
FROM alpine:3.19
WORKDIR /app
COPY --from=builder /app/main .

ENV API_PORT=80
EXPOSE 80
CMD [ "/app/main" ]

