.PHONY: docker-build
docker-build: 
	@docker build -t coupon-service:dev .

.PHONY: docker-run
docker-run: docker-build
	@docker run -it --rm -p 8080:8080 coupon-service:dev

.PHONY: clean-up-mocks
clean-up-mocks:
	@rm -rf internal/mocks
	
.PHONY: gen-mocks
gen-mocks: clean-up-mocks
	@docker run -v `pwd`:/src -w /src vektra/mockery:v2.51
