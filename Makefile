test:
	go test -v -cover -short ./...

run:
	go run ./cmd/coupon_service/

.PHONY: swag-init
swag:
	swag init --parseDependency --parseInternal -g ./internal/api/api.go -o ./docs

