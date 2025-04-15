SHELL := /bin/bash
BINARY_NAME := tagstore
BUILD_DIR := ./bin
GO := go
GOFLAGS := -v
LDFLAGS := -ldflags "-s -w"

.PHONY: build test run clean fmt lint help

help:
	@echo "Usage: make <target>"
	@echo ""
	@echo "Targets:"
	@echo "  build    Build the binary"
	@echo "  test     Run tests"
	@echo "  run      Run the application"
	@echo "  clean    Remove build artifacts"
	@echo "  fmt      Format source code"
	@echo "  lint     Run linters"

build:
	@mkdir -p $(BUILD_DIR)
	$(GO) build $(GOFLAGS) $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

run:
	$(GO) run $(GOFLAGS) main.go

test:
	$(GO) test $(GOFLAGS) ./...

clean:
	@rm -rf $(BUILD_DIR)

fmt:
	$(GO) fmt ./...

lint:
	golangci-lint run
