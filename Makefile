UNIT_TEST_ONLY_PKGS=$(shell go list ./... | grep -v "/vendor/")

.PHONY: test
test:
	go test -v $(UNIT_TEST_ONLY_PKGS) -covermode=count -coverprofile=coverage.out

.PHONY: deps
deps:
	govendor sync

.PHONY: build
build:
	# same as on heroku (in vendor/vendor.json["heroku"]).
	go build ./cmd/server/

.PHONY: run
run:
	exec ./server
