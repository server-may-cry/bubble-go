.PHONY: test
test:
	# vgo use 'clang' as C compiler by some reason (need 'gcc'). sqlite need CGO
	CC=gcc vgo test ./... -coverprofile coverage.out

.PHONY: bench
bench:
	CC=gcc vgo test -cpuprofile=cpu.pprof -memprofile=mem.pprof -benchmem -bench=. ./notification

.PHONY: build
build:
	GOOS=linux CGO_ENABLED=0 vgo build
