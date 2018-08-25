.PHONY: test
test:
	# golangci-lint run --enable-all --disable=gas
	go test ./... -coverprofile coverage.out

.PHONY: bench
bench:
	go test -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem -bench=. ./notification

.PHONY: build
build:
	GOOS=linux CGO_ENABLED=0 go build -ldflags "-X main.version=`git describe --always`"
