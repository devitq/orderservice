# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOTEST=$(GOCMD) test
GODOWNLOAD=$(GOCMD) mod download
BINARY_NAME=orderservice
BINARY_DIR=bin

# Protobuf parameters
PROTOC=protoc
PROTO_DIR=api/proto
PROTO_FILE=$(PROTO_DIR)/order.proto
PROTO_OUT=.

.PHONY: install i generate gen protoc test build run lint fmt format clean help

install:
	$(GODOWNLOAD)

i: install

generate:
	$(PROTOC) --version || (echo "protoc not found, install protoc"; exit 1)
	$(PROTOC) --go_out=$(PROTO_OUT) --go-grpc_out=$(PROTO_OUT) \
		$(PROTO_FILE)

gen: generate

protoc: generate

test:
	$(GOTEST) ./...

build:
	$(GOBUILD) -o ./$(BINARY_DIR)/$(BINARY_NAME) ./cmd/server
	chmod +x ./$(BINARY_DIR)/$(BINARY_NAME)

run: build
	./$(BINARY_DIR)/$(BINARY_NAME)

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
