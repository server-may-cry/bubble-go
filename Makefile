UNIT_TEST_ONLY_PKGS=$(shell go list ./... | grep -v "/vendor/")

test:
	@echo "run unit tests with coverage"
	@go test -v -cover $(UNIT_TEST_ONLY_PKGS)
