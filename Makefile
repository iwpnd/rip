VERSION ?= edge
GITSHA    := $(shell git rev-parse --short HEAD)
BUILDTIME := $(shell date +%FT%T%z)
GOARCH ?= amd64
CGO_ENABLED ?= 0

.PHONY: build
build:
	@echo "building"
	@echo "VERSION=$(VERSION) GITSHA=$(GITSHA) BUILDTIME=$(BUILDTIME)"
	GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build .

.PHONY: install
install:
	@echo "Checking for version manager..."
	@if command -v mise >/dev/null 2>&1; then \
		echo "Using mise to install dependencies"; \
		mise install; \
	else \
		echo "Error: Neither mise not installed"; \
		echo "Please install from:"; \
		echo "  • mise (recommended): https://mise.jdx.dev"; \
		exit 1; \
	fi
	@echo "Installing pre-commit hooks..."
	pre-commit install
	@echo "Installation complete!"

.PHONY: test
test:
	@echo "running tests"
	@go test -v -json ./... | tparse -all

.PHONY: lint
lint:
	@echo "running linter"
	@golangci-lint run
