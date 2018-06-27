GO := go
TEST_PACKAGES := . ./lib ./models ./worker ./utils ./oauth ./tmpl ./data ./runner

build:
	$(GO) generate
	$(GO) build

test:
	$(GO) test $(TEST_PACKAGES) -race $(ARGS)
