GO := go
TEST_PACKAGES := . ./lib ./models ./worker ./utils ./oauth ./tmpl ./data ./runner

build:
	$(GO) generate
	$(GO) build

test:
	$(GO) test $(TEST_PACKAGES)

test/full:
	$(GO) test -cover -race $(TEST_PACKAGES)
