BINARY=gitlabmcp
VERSION?=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")

.PHONY: build test lint clean

build:
	go build -ldflags "-X main.version=$(VERSION)" -o $(BINARY) ./cmd/gitlabmcp

test:
	go test ./...

lint:
	go vet ./...

clean:
	rm -f $(BINARY)
