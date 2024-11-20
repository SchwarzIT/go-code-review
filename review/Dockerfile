FROM golang:1.18-alpine AS builder

RUN apk --no-cache add  git gcc libc-dev

WORKDIR /go/src/coupon-service

COPY go.mod go.sum ./

RUN go mod tidy

COPY . .

WORKDIR /go/src/coupon-service/cmd/coupon_service/

RUN go build -o coupon-service .

FROM alpine:3.20

RUN apk --no-cache add ca-certificates

WORKDIR /root/

COPY --from=builder /go/src/coupon-service/cmd/coupon_service .

ENTRYPOINT ["./coupon-service"]

EXPOSE 8080

CMD ["./coupon-service"]

