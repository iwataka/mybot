GO := go
TEST_PACKAGES := . ./lib ./models ./worker ./utils ./oauth ./tmpl ./data ./runner
GOLINT := gometalinter

DOCKER_COMPOSE := docker-compose
DOCKER_COMPOSE_SCRIPT := scripts/docker-compose.yml

DOCKER := docker
DOCKER_REPO := iwataka/mybot

build:
	$(GO) generate
	$(GO) build

test:
	$(GO) test $(TEST_PACKAGES) -race $(args)

lint:
	$(GOLINT) --deadline=30s

deploy_app:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_SCRIPT) up -d

deploy_app_single:
	$(DOCKER) run -d --name mybot -p 8080:8080 $(DOCKER_REPO)

clean_app:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_SCRIPT) stop
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_SCRIPT) rm -f

update_images:
	$(DOCKER_COMPOSE) -f $(DOCKER_COMPOSE_SCRIPT) pull

create_image:
	$(DOCKER) build -t $(DOCKER_REPO) .
