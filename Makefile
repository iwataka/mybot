GO := go
GOLINT := golangci-lint
FMT_OPTS = -s -w
LINT_OPTS =
TEST_OPTS = -race

all: fmt lint build test

fmt:
	gofmt $(FMT_OPTS) .
	yarn --cwd ./web format

lint:
	$(GOLINT) run $(LINT_OPTS) ./...
	yarn --cwd ./web lint

build:
	go build
	yarn --cwd ./web build

test:
	$(GO) test $(TEST_OPTS) ./...
	yarn --cwd ./web test --watchAll=false
