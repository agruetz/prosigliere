# Makefile for building Protocol Buffer definitions

# Variables
GO := go
BUF := $(shell which buf || echo "$(GO) run github.com/bufbuild/buf/cmd/buf")
FLYWAY := $(shell which flyway || echo "docker run --rm -v $(PWD)/db:/flyway/conf -v $(PWD)/db/migrations:/flyway/sql flyway/flyway")
MOCKERY := $(shell which mockery || echo "$(GO) run github.com/vektra/mockery/v2")

# Default target
.PHONY: all
all: generate mocks

# Install tools
.PHONY: tools
tools:
	$(GO) install github.com/bufbuild/buf/cmd/buf@latest
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	$(GO) install github.com/bufbuild/protovalidate-go/cmd/protoc-gen-validate-go@latest
	$(GO) install github.com/vektra/mockery/v2@latest

# Generate code from proto files
.PHONY: generate
generate:
	$(BUF) dep update
	$(BUF) generate

# Generate mocks
.PHONY: mocks
mocks:
	$(MOCKERY) --dir=internal/datastore --name=Store --output=internal/datastore/mocks --outpkg=mocks --filename=store.go

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
	find ./docs -name "*.json" -delete
	find ./docs -name "*.swagger.json" -delete
	rm -rf internal/datastore/mocks/*

# Database migration targets
.PHONY: db-migrate db-clean db-info db-validate db-repair

db-migrate:
	$(FLYWAY) migrate

db-clean:
	$(FLYWAY) clean

db-info:
	$(FLYWAY) info

db-validate:
	$(FLYWAY) validate

db-repair:
	$(FLYWAY) repair

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        - Generate code from proto files and mocks (default)"
	@echo "  deps       - Install required dependencies"
	@echo "  generate   - Generate code from proto files"
	@echo "  mocks      - Generate mocks for interfaces"
	@echo "  lint       - Lint proto files"
	@echo "  breaking   - Check for breaking changes against main branch"
	@echo "  clean      - Remove generated files"
	@echo "  db-migrate - Run database migrations"
	@echo "  db-clean   - Clean the database"
	@echo "  db-info    - Show information about migrations"
	@echo "  db-validate - Validate applied migrations"
	@echo "  db-repair  - Repair the schema history table"
	@echo "  help       - Show this help message"
