GO := go
GOLINT := golangci-lint

all: fmt lint test

fmt:
	$(GO) fmt ./...

lint:
	$(GOLINT) run ./...

test:
	$(GO) test -race ./...
