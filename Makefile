UNIT_TEST_ONLY_PKGS=$(shell go list ./... | grep -v "/vendor/")

COVERALLS_IGNORE := cmd/server/main.go

.PHONY: test
test:
	@echo "run unit tests with coverage"
	go test -v -cover $(UNIT_TEST_ONLY_PKGS)

.PHONY: deps
deps:
	govendor sync

.PHONY: build
build:
	# same as on heroku (in vendor/vendor.json["heroku"]).
	go build ./cmd/server/

.PHONY: run
run:
	./server

.PHONY: all
all: deps test build run
