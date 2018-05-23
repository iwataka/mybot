GO := go
TEST_PACKAGES := . ./lib ./models ./worker ./utils ./oauth ./tmpl ./data ./runner

generate:
	$(GO) generate

build: generate
	$(GO) build

test: generate
	$(GO) test $(TEST_PACKAGES)

cover: generate
	$(GO) test -cover $(TEST_PACKAGES)
