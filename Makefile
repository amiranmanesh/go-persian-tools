GO ?= go
GOLANGCI_LINT_VERSION ?= v2.12.2
FUZZTIME ?= 30s

.DEFAULT_GOAL := check

.PHONY: check
check: fmt vet lint test ## Run everything CI runs

.PHONY: build
build: ## Build the persian-tools CLI into ./persian-tools
	$(GO) build -trimpath -ldflags "-s -w" -o persian-tools ./cmd/persian-tools

.PHONY: install
install: ## Install the persian-tools CLI into $$GOPATH/bin
	$(GO) install ./cmd/persian-tools

.PHONY: image
image: ## Build the container image locally as persian-tools:dev
	docker build -t persian-tools:dev --build-arg VERSION=dev .

.PHONY: test
test: ## Run the tests with the race detector
	$(GO) test -race -shuffle=on ./...

.PHONY: cover
cover: ## Run the tests and open the coverage report
	$(GO) test -coverprofile=coverage.out -covermode=atomic ./...
	$(GO) tool cover -func=coverage.out
	$(GO) tool cover -html=coverage.out

.PHONY: bench
bench: ## Run the benchmarks
	$(GO) test -run='^$$' -bench=. -benchmem ./...

.PHONY: fuzz
fuzz: ## Run every fuzz target for $(FUZZTIME)
	@# -fuzz takes one package at a time, hence the nested loop.
	@for pkg in $$($(GO) list ./...); do \
		for target in $$($(GO) test -list 'Fuzz.*' $$pkg | grep '^Fuzz' || true); do \
			echo "==> $$pkg $$target"; \
			$(GO) test -run='^$$' -fuzz="^$$target$$" -fuzztime=$(FUZZTIME) $$pkg || exit 1; \
		done; \
	done

.PHONY: fmt
fmt: ## Format the source
	$(GO) fmt ./...

.PHONY: vet
vet: ## Run go vet
	$(GO) vet ./...

.PHONY: lint
lint: ## Run golangci-lint
	$(GO) run github.com/golangci/golangci-lint/v2/cmd/golangci-lint@$(GOLANGCI_LINT_VERSION) run ./...

.PHONY: vulncheck
vulncheck: ## Check the dependency graph for known vulnerabilities
	$(GO) run golang.org/x/vuln/cmd/govulncheck@latest ./...

.PHONY: docs
docs: ## Serve the package documentation at http://localhost:6060
	$(GO) run golang.org/x/pkgsite/cmd/pkgsite@latest -open .

.PHONY: clean
clean: ## Remove build and coverage artifacts
	rm -f coverage.out coverage.txt
	$(GO) clean -testcache

.PHONY: help
help: ## List the available targets
	@grep -hE '^[a-zA-Z_-]+:.*?## ' $(MAKEFILE_LIST) \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'
