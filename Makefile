GO := go
TEST_PACKAGES := . ./lib ./models ./worker ./utils ./oauth ./tmpl ./data ./runner

DOCKER_COMPOSE := docker-compose
DOCKER_COMPOSE_SCRIPT := scripts/docker-compose.yml

build:
	$(GO) generate
	$(GO) build

test:
	$(GO) test $(TEST_PACKAGES) -race $(args)

deploy_app:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_SCRIPT) up -d

clean_app:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_SCRIPT) stop
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_SCRIPT) rm -f
