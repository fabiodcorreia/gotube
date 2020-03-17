COV_REPORT	:= coverage.out
PKG			:= ./internal/...


cov:
	go test -coverprofile=$(COV_REPORT) $(PKG)  -v -covermode=atomic $(PKG)
	go tool cover -func=$(COV_REPORT)

test:
	go test $(PKG) -v -covermode=atomic

.PHONY: format
format:
	gofmt -s -w ./*/*.go
	gofmt -s -w *.go

bench:
	#go test -benchtime=10s -count=4 -benchmem -bench='Benchmark*' ./internal/serial 2>&1 | tee benchmarks/serial.txt
	#go test -benchtime=10s -count=4 -benchmem -bench='Benchmark*' ./internal/parse 2>&1 | tee benchmarks/parse.txt
	#go test -benchtime=10s -count=4 -benchmem -bench='Benchmark*' ./internal/info 2>&1 | tee benchmarks/info.txt
	go test -benchtime=10s -count=4 -benchmem -bench='BenchmarkGetVideoInfo' ./internal/info 2>&1

analyze:
	go build -gcflags=-m=2 $(PKG) 2>&1
	
analyze-inline:
	go build -gcflags=-m=2 $(PKG) 2>&1 | grep "cannot inline"

analyze-escapes:
	go build -gcflags=-m=2 $(PKG) 2>&1 | grep "escapes to heap"

analyze-leaking:
	go build -gcflags=-m=2 $(PKG) 2>&1 | grep "leaking"

.PHONY: bench