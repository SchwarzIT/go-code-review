# build stage
FROM golang:latest AS builder

RUN apk add git gcc libc-dev

WORKDIR /go/src/coupon-service

COPY . .

WORKDIR /go/src/coupon-service/cmd/export
RUN go build -o main .

ENTRYPOINT ./main
EXPOSE 8080
