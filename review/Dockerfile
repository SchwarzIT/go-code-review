# build stage
FROM golang:latest AS builder

RUN apt-get update && apt-get install -y gcc libc-dev

WORKDIR /go/src/coupon-service

COPY . .

WORKDIR /go/src/coupon-service/cmd/coupon_service
RUN go build -o main .

ENTRYPOINT ./main
EXPOSE 8080
