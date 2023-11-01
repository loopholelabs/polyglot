BENCHCOUNT ?= 10

install:
	- go install ../protoc-gen-go-polyglot
	- go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest
	- go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28

generate:
	- mkdir -p polyglot
	- protoc --go-polyglot_out=polyglot bench.proto
	- mkdir -p vtproto
	- protoc --go_out=vtproto --go-vtproto_out=vtproto bench.proto

benchmark-polyglot: benchmark-polyglot-cmp
	- go run -mod=mod golang.org/x/perf/cmd/benchstat bench.txt
	- go mod tidy
	- rm -rf bench.txt

benchmark-polyglot-cmp:
	- go test -bench=. -timeout=24h -count=$(BENCHCOUNT) ./... -test.short | tee bench.txt

benchmark-polyglot-long:
	- go test -bench=. -timeout=24h -count=$(BENCHCOUNT) ./... | tee bench.txt
	- go run -mod=mod golang.org/x/perf/cmd/benchstat bench.txt
	- go mod tidy
	- rm -rf bench.txt

benchmark-vtproto:
	- go test -bench=. -timeout=24h -count=$(BENCHCOUNT) ./... -tags=vtproto -test.short | tee bench.txt
	- go run -mod=mod golang.org/x/perf/cmd/benchstat bench.txt
	- go mod tidy
	- rm -rf bench.txt

benchmark-vtproto-long:
	- go test -bench=. -timeout=24h -count=$(BENCHCOUNT) ./... -tags=vtproto | tee bench.txt
	- go run -mod=mod golang.org/x/perf/cmd/benchstat bench.txt
	- go mod tidy
	- rm -rf bench.txt

leaks:
	- go test -bench=. -gcflags="-m=2" ./...