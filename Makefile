BINARY  := hc
PKGS    := ./...
GOFLAGS :=

.DEFAULT_GOAL := help

.PHONY: help build run test race cover fmt fmt-check vet lint tidy clean

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

build: ## Build the hc binary
	go build $(GOFLAGS) -o $(BINARY) .

run: build ## Build and run the CLI
	./$(BINARY)

test: ## Run all tests
	go test $(GOFLAGS) $(PKGS)

race: ## Run tests with race detector and coverage (same as CI)
	go test $(GOFLAGS) -race -coverprofile=coverage.out -covermode=atomic $(PKGS)

cover: ## Run tests and print coverage percentage
	go test $(GOFLAGS) -cover $(PKGS)

fmt: ## Format all Go files in place
	gofmt -w .

fmt-check: ## Fail if any file is not gofmt-formatted (same as CI)
	@unformatted=$$(gofmt -l .); \
	if [ -n "$$unformatted" ]; then \
		echo "Files not gofmt-formatted:"; \
		echo "$$unformatted"; \
		exit 1; \
	fi
	@echo "gofmt: OK"

vet: ## Run go vet
	go vet $(GOFLAGS) $(PKGS)

lint: ## Run golangci-lint (requires golangci-lint installed)
	golangci-lint run $(PKGS)

tidy: ## Tidy go.mod / go.sum
	go mod tidy

ci: fmt-check vet lint race ## Run everything CI runs (lint + tests)
	@echo "CI checks passed"

clean: ## Remove build artifacts
	rm -f $(BINARY) coverage.out
