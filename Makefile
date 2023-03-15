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
test: test_units test_static test_gophermart

.PHONY: test_units
test_units: 
	@echo "\n### $@"
	@go test ./... -v --count 1

.PHONY: test_static
test_static:
	@echo "\n### $@"
	@go vet -vettool=./statictest ./...

.PHONY: test_gophermart
test_devops: build
	@echo "\n### $@"
	@echo "DON'T FORGET TO START postgres.sh\n"

.PHONY: clean
clean:
	@echo "\n### $@"
	@rm -rf cmd/gophermart/gophermart
