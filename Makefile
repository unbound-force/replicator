.PHONY: build test lint vet clean

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
LDFLAGS := -ldflags "-X main.Version=$(VERSION)"

build:
	go build $(LDFLAGS) -o bin/replicator ./cmd/replicator

test:
	go test ./... -count=1

vet:
	go vet ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/

# Run the MCP server (for local testing with AI agents).
serve: build
	./bin/replicator serve

# Quick check: vet + test.
check: vet test
