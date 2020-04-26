GO := go
GOLINT := golangci-lint

all: fmt lint test

fmt:
	gofmt -s -w .

lint:
	$(GOLINT) run ./...

test:
	$(GO) test -race ./...
