# Makefile for building Protocol Buffer definitions

# Variables
GO := go
BUF := $(shell which buf || echo "$(GO) run github.com/bufbuild/buf/cmd/buf")

# Default target
.PHONY: all
all: generate

# Install tools
.PHONY: tools
tools:
	$(GO) install github.com/bufbuild/buf/cmd/buf@latest
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	$(GO) install github.com/bufbuild/protovalidate-go/cmd/protoc-gen-validate-go@latest

# Generate code from proto files
.PHONY: generate
generate:
	$(BUF) dep update
	$(BUF) generate

# Lint proto files
.PHONY: lint
lint:
	$(BUF) lint

# Check for breaking changes
.PHONY: breaking
breaking:
	$(BUF) breaking --against '.git#branch=main'

# Clean generated files
.PHONY: clean
clean:
	find . -name "*.pb.go" -delete
	find . -name "*.pb.validate.go" -delete
	find . -name "*.pb.gw.go" -delete

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        - Generate code from proto files (default)"
	@echo "  deps       - Install required dependencies"
	@echo "  generate   - Generate code from proto files"
	@echo "  lint       - Lint proto files"
	@echo "  breaking   - Check for breaking changes against main branch"
	@echo "  clean      - Remove generated files"
	@echo "  help       - Show this help message"
