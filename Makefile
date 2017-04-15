UNIT_TEST_ONLY_PKGS=$(shell go list ./... | grep -v "/vendor/")

test:
	@echo "run unit tests with coverage"
	@go test -v -cover $(UNIT_TEST_ONLY_PKGS)

build:
	# same as on heroku (in vendor/vendor.json["heroku"]).
	@go build ./cmd/server/
