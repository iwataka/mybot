GO := go
GOLINT := golangci-lint
FMT_OPTS = -s -w
LINT_OPTS =
TEST_OPTS = -race

all: fmt lint test

fmt:
	gofmt $(FMT_OPTS) .

lint:
	$(GOLINT) run $(LINT_OPTS) ./...

test:
	$(GO) test $(TEST_OPTS) ./...
