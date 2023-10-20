BENCHCOUNT ?= 10

generate:
	- mkdir -p vtproto && protoc --go_out=vtproto --go-vtproto_out=vtproto bench.proto
	- mkdir -p polyglot && protoc --go-polyglot_out=polyglot bench.proto

benchmark:
	- go test -bench=. -timeout=24h -count=$(BENCHCOUNT) ./... | tee bench.txt && go run -mod=mod golang.org/x/perf/cmd/benchstat bench.txt && go mod tidy && rm -rf bench.txt

leaks:
	- go test -bench=. -gcflags="-m=2" ./...