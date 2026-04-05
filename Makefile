.PHONY: build test lint vet clean serve check release install

VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT  ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE    ?= $(shell date -u +%Y-%m-%dT%H:%M:%SZ)
LDFLAGS := -ldflags "-X main.Version=$(VERSION) -X main.commit=$(COMMIT) -X main.date=$(DATE)"

build:
	go build $(LDFLAGS) -o bin/replicator ./cmd/replicator

test:
	go test ./... -count=1

vet:
	go vet ./...

lint:
	golangci-lint run ./...

clean:
	rm -rf bin/ dist/

# Run the MCP server (for local testing with AI agents).
serve: build
	./bin/replicator serve

# Quick check: vet + test.
check: vet test

# Local release dry-run (no publish).
release:
	goreleaser release --snapshot --clean

# Install to GOPATH/bin.
install:
	go install $(LDFLAGS) ./cmd/replicator
