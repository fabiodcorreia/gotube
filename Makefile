COV_REPORT	:= coverage.out
PKG			:= . ./internal/...
SRC = $(shell find . -type f -name '*.go' -not -path "./vendor/*")

test:
	go test $(PKG) -covermode=atomic

cov: format
	go test -coverprofile=$(COV_REPORT) $(PKG)  -v -covermode=atomic $(PKG)
	go tool cover -func=$(COV_REPORT)

format:
	gofmt -s -l -w $(SRC)

bench:
	go test -benchtime=10s -count=4 -benchmem -bench='Benchmark*' ./internal/serial 2>&1 | tee benchmarks/serial.txt
	go test -benchtime=10s -count=4 -benchmem -bench='Benchmark*' ./internal/parse 2>&1 | tee benchmarks/parse.txt
	go test -benchtime=10s -count=4 -benchmem -bench='Benchmark*' ./internal/info 2>&1 | tee benchmarks/info.txt

.PHONY: bench format cov test