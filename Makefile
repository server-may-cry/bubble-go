VERSION := $(shell git describe --always --long --dirty)

.PHONY: all
all: lint test bench build

.PHONY: bench
bench:
	go test -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem -bench=. ./notification

.PHONY: build
build:
	GOOS=linux CGO_ENABLED=0 go build -ldflags "-X main.version=${VERSION}"

.PHONY: clean
clean:
	rm -rf bubble-go* *.out *.test *.pprof

.PHONY: lint
lint:
	golangci-lint run --enable-all --disable=gas,gochecknoglobals

.PHONY: start
start:
	docker-compose up --build

.PHONY: test
test:
	go test ./... -coverprofile coverage.out
