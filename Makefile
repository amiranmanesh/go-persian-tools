GO      ?= go
LINTER  ?= golangci-lint

.PHONY: all test cover race lint fmt vet tidy example clean

## all: run vet + tests
all: vet test

## test: run all tests
test:
	$(GO) test ./...

## race: run tests with the race detector
race:
	$(GO) test -race ./...

## cover: run tests and generate a coverage profile
cover:
	$(GO) test -race -covermode=atomic -coverprofile=coverage.txt ./...
	$(GO) tool cover -func=coverage.txt | tail -1

## lint: run golangci-lint (install: https://golangci-lint.run)
lint:
	$(LINTER) run

## fmt: format the code
fmt:
	$(GO) fmt ./...

## vet: run go vet
vet:
	$(GO) vet ./...

## tidy: tidy go.mod / go.sum
tidy:
	$(GO) mod tidy

## example: run the example program
example:
	$(GO) run ./examples

## clean: remove build/coverage artifacts
clean:
	$(GO) clean
	rm -f coverage.txt
