UNIT_TEST_ONLY_PKGS=$(shell go list ./... | grep -v "/vendor/")

.PHONY: test
test:
	go test -v $(UNIT_TEST_ONLY_PKGS)

.PHONY: deps
deps:
	dep ensure

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo ./cmd/server/
