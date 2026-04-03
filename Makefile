GO ?= go
GOCACHE ?= $(CURDIR)/.cache/go-build
GOLANGCI_LINT_CACHE ?= $(CURDIR)/.cache/golangci-lint
GOLANGCI_LINT ?= golangci-lint
PKG ?= ./...
GOFILES := $(shell find pkg -type f -name '*.go' | sort)
COVERPROFILE ?= .coverage/coverage.out

.PHONY: help fmt fmt-check vet test test-race integration-test cover lint ci clean

help:
	@printf "Targets:\n"
	@printf "  make fmt        Format Go sources\n"
	@printf "  make fmt-check  Fail if formatting is needed\n"
	@printf "  make vet        Run go vet\n"
	@printf "  make test       Run unit tests\n"
	@printf "  make test-race  Run tests with the race detector\n"
	@printf "  make integration-test  Run live MLB API integration tests (opt-in)\n"
	@printf "  make cover      Run tests with coverage output\n"
	@printf "  make lint       Run golangci-lint\n"
	@printf "  make ci         Run the local CI contract\n"
	@printf "  make clean      Remove local build artifacts\n"

fmt:
	@$(GO) fmt $(PKG)

fmt-check:
	@test -z "$$(gofmt -l $(GOFILES))" || (echo "gofmt reported unformatted files"; gofmt -l $(GOFILES); exit 1)

vet:
	@mkdir -p $(dir $(GOCACHE))
	@GOCACHE=$(GOCACHE) $(GO) vet $(PKG)

test:
	@mkdir -p $(dir $(GOCACHE))
	@GOCACHE=$(GOCACHE) $(GO) test $(PKG)

test-race:
	@mkdir -p $(dir $(GOCACHE))
	@GOCACHE=$(GOCACHE) $(GO) test -race $(PKG)

integration-test:
	@mkdir -p $(dir $(GOCACHE))
	@MLBAPI_RUN_INTEGRATION=1 GOCACHE=$(GOCACHE) $(GO) test -count=1 -run 'TestLive' $(PKG)

cover:
	@mkdir -p $(dir $(GOCACHE)) $(dir $(COVERPROFILE))
	@GOCACHE=$(GOCACHE) $(GO) test -coverprofile=$(COVERPROFILE) $(PKG)
	@GOCACHE=$(GOCACHE) $(GO) tool cover -func=$(COVERPROFILE)

lint:
	@command -v $(GOLANGCI_LINT) >/dev/null 2>&1 || (echo "golangci-lint not found; install it to use 'make lint'"; exit 1)
	@mkdir -p $(dir $(GOCACHE)) $(GOLANGCI_LINT_CACHE)
	@GOCACHE=$(GOCACHE) GOLANGCI_LINT_CACHE=$(GOLANGCI_LINT_CACHE) $(GOLANGCI_LINT) run $(PKG)

ci: fmt-check vet test

clean:
	@rm -rf .cache .coverage
