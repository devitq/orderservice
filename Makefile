# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build -trimpath -ldflags="-s -w" 
GOTEST=$(GOCMD) test
GODOWNLOAD=$(GOCMD) mod download
BINARY_NAME=orderservice
MIGRATE_BINARY_NAME=orderservice-migrate
BINARY_DIR=bin

# Protobuf parameters
PROTOC=protoc
PROTO_DIR=api/proto
PROTO_FILE=$(PROTO_DIR)/order.proto
PROTO_OUT=.

.PHONY: install i generate gen generate-gw test build run migrate lint fmt format clean help

install:
	$(GODOWNLOAD)

i: install

generate:
	$(PROTOC) --version || (echo "protoc not found, install protoc"; exit 1)
	$(PROTOC) --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) \
		$(PROTO_FILE)

gen: generate

generate-gw:
	$(PROTOC) --version || (echo "protoc not found, install protoc"; exit 1)
	$(PROTOC) --grpc-gateway_out=$(PROTO_OUT) --grpc-gateway_opt generate_unbound_methods=true \
		$(PROTO_FILE)

test:
	$(GOTEST) ./...

build:
	$(GOBUILD) -o ./$(BINARY_DIR)/$(BINARY_NAME) ./cmd/server
	chmod +x ./$(BINARY_DIR)/$(BINARY_NAME)

build-migrate:
	$(GOBUILD) -o ./$(BINARY_DIR)/$(MIGRATE_BINARY_NAME) ./cmd/migrate
	chmod +x ./$(BINARY_DIR)/$(MIGRATE_BINARY_NAME)

run: build
	./$(BINARY_DIR)/$(BINARY_NAME)

migrate: build-migrate
	@cmd=$(word 2,$(MAKECMDGOALS)); \
	if [ -z "$$cmd" ]; then \
		echo "Usage: make migrate <command>"; \
		exit 1; \
	fi; \
	./$(BINARY_DIR)/$(MIGRATE_BINARY_NAME) -cmd $$cmd

lint:
	golangci-lint run -c .golangci.yaml ./...

fmt:
	golangci-lint fmt -c .golangci.yaml ./...

format: fmt

clean:
	rm -rf bin/*

help:
	@echo "Available commands:"
	@echo "  install   - Install all deps using go mod download"
	@echo "        i"
	@echo "  generate  - Generate gRPC code"
	@echo "       gen"
	@echo "    protoc"
	@echo "  test      - Run tests"
	@echo "  build     - Build the binary"
	@echo "  run       - Run the application"
	@echo "  lint      - Run golangci-lint linter"
	@echo "  format    - Run golangci-lint formatter"
	@echo "     fmt"
	@echo "  clean     - Clean build artifacts"

.DEFAULT_GOAL := help

%:
	@:
