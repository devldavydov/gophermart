.PHONY: all
all: clean build test

.PHONY: prepare_env
prepare_env:
	@echo "\n### $@"
	@wget --no-check-certificate https://github.com/Yandex-Practicum/go-autotests/releases/download/v0.7.13/gophermarttest -O ./gophermarttest
	@chmod u+x ./gophermarttest
	@wget --no-check-certificate https://github.com/Yandex-Practicum/go-autotests/releases/download/v0.7.13/statictest -O ./statictest
	@chmod u+x ./statictest
	@wget --no-check-certificate https://github.com/Yandex-Practicum/go-autotests/releases/download/v0.7.11/random -O ./random
	@chmod u+x ./random

.PHONY: build
build:
	@echo "\n### $@"
	@cd cmd/gophermart && go build .

.PHONY: test
test: test_units test_integration test_static test_gophermart

.PHONY: test_units
test_units: 
	@echo "\n### $@"
	@go test ./... -v --count 1

.PHONY: test_integration
test_integration: 
	@echo "\n### $@"
	@echo "DON'T FORGET TO START postgres.sh AND cmd/gophermart/gophermart\n"
	@export TEST_GOPHERMART_SRV_ADDR=http://127.0.0.1:8080 && \
	 export TEST_ACCRUAL_SRV_LISTEN_ADDR=127.0.0.1:9090 && \
	 go test ./... -v -run=^TestGophermartIntegration$$ --count 1

.PHONY: test_static
test_static:
	@echo "\n### $@"
	@go vet -vettool=./statictest ./...

.PHONY: test_gophermart
test_gophermart: build
	@echo "\n### $@"
	@echo "DON'T FORGET TO START postgres.sh\n"
	@./gophermarttest -test.v -test.run=^TestGophermart$$ \
     -gophermart-binary-path=cmd/gophermart/gophermart \
     -gophermart-host=localhost \
     -gophermart-port=8080 \
     -gophermart-database-uri="postgresql://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable" \
     -accrual-binary-path=cmd/accrual/accrual_linux_amd64 \
     -accrual-host=localhost \
     -accrual-port=$$(./random unused-port) \
     -accrual-database-uri="postgresql://postgres:postgres@127.0.0.1:5432/praktikum?sslmode=disable"

.PHONY: clean
clean:
	@echo "\n### $@"
	@rm -rf cmd/gophermart/gophermart
