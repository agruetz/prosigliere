# Makefile for building Protocol Buffer definitions and server binary

# Variables
GO := go
BUF := $(shell which buf || echo "$(GO) run github.com/bufbuild/buf/cmd/buf")
FLYWAY := $(shell which flyway || echo "docker run --rm -v $(PWD)/db:/flyway/conf -v $(PWD)/db/migrations:/flyway/sql flyway/flyway")
MOCKERY := $(shell which mockery || echo "$(GO) run github.com/vektra/mockery/v2")

# Server build variables
SERVER_BINARY_NAME := server
SERVER_OUTPUT_DIR := cmd/server
SERVER_MAIN_FILE := cmd/server/server.go

# Go build flags for static compilation
# -s -w: strip debugging information
# -extldflags "-static": use static linking
# CGO_ENABLED=0: disable CGO for pure Go compilation
GO_BUILD_FLAGS := -ldflags="-s -w -extldflags '-static'"

# Default target
.PHONY: all
all: proto mocks build-server

# Install tools
.PHONY: tools
tools:
	$(GO) install github.com/bufbuild/buf/cmd/buf@latest
	$(GO) install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	$(GO) install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
	$(GO) install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@latest
	$(GO) install github.com/bufbuild/protovalidate-go/cmd/protoc-gen-validate-go@latest
	$(GO) install github.com/vektra/mockery/v2@latest

# Update protocol buffer dependencies
.PHONY: update-proto-deps
update-proto-deps:
	$(BUF) dep update

# Generate code from proto files
.PHONY: generate-protos
generate-protos:
	$(BUF) generate

# Generate and update protocol buffer dependencies
.PHONY: proto
proto: update-proto-deps generate-protos

# Generate mocks
.PHONY: mocks
mocks:
	$(MOCKERY) --dir=internal/datastore --name=Store --output=internal/datastore/mocks --outpkg=mocks --filename=store.go

# Clean mock files
.PHONY: clean-mocks
clean-mocks:
	rm -rf internal/datastore/mocks/*

# Lint proto files
.PHONY: lint
lint:
	$(BUF) lint

# Check for breaking changes
.PHONY: breaking
breaking:
	$(BUF) breaking --against '.git#branch=main'

# Clean proto generated files
.PHONY: clean-protos
clean-protos:
	find . -name "*.pb.go" -delete
	find . -name "*.pb.validate.go" -delete
	find . -name "*.pb.gw.go" -delete
	find ./docs -name "*.json" -delete
	find ./docs -name "*.swagger.json" -delete

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

# Build the server binary
.PHONY: build-server
build-server:
	CGO_ENABLED=0 $(GO) build -o $(SERVER_OUTPUT_DIR)/$(SERVER_BINARY_NAME) $(GO_BUILD_FLAGS) $(SERVER_MAIN_FILE)

# Clean server build artifacts
.PHONY: clean-server
clean-server:
	rm -f $(SERVER_OUTPUT_DIR)/$(SERVER_BINARY_NAME)

# Clean generated files
.PHONY: clean
clean: clean-server clean-mocks clean-protos

# Help target
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  all        - Generate code from proto files, mocks, and build server (default)"
	@echo "  deps       - Install required dependencies"
	@echo "  proto      - Update proto dependencies and generate code from proto files"
	@echo "  update-proto-deps - Update protocol buffer dependencies"
	@echo "  generate-protos - Generate code from proto files"
	@echo "  mocks      - Generate mocks for interfaces"
	@echo "  build-server - Build the server binary"
	@echo "  lint       - Lint proto files"
	@echo "  breaking   - Check for breaking changes against main branch"
	@echo "  clean      - Remove generated files"
	@echo "  clean-server - Remove server build artifacts"
	@echo "  clean-mocks  - Remove generated mock files"
	@echo "  clean-protos - Remove generated protocol buffer files"
	@echo "  db-migrate - Run database migrations"
	@echo "  db-clean   - Clean the database"
	@echo "  db-info    - Show information about migrations"
	@echo "  db-validate - Validate applied migrations"
	@echo "  db-repair  - Repair the schema history table"
	@echo "  help       - Show this help message"
