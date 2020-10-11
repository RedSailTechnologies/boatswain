.PHONY: client echo help kraken target version

# ENV
DEBUG=false
DOCKER_BUILDKIT=1
DOCKER_REPO=
DOCKER_TAG=latest
GEN_DOC=docs/api/
GEN_GO=rpc/
GEN_TS=$(TRITON_PATH)src/app/services/
LEVI_CMD=cmd/leviathan/
LEVI_OUT=bin/
PROJECT_NAME=null
SERVICE_LIST=kraken
TRITON_PATH=web/triton/
WORKDIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# BASIC TARGETS
## all: builds the client and all services
all: echo proto kraken client push

## clean: removes binaries, images, etc
clean:
	@echo Running clean
.SILENT:
	@rm -f main
	@rm -rf $(TRITON_PATH)dist
	@rm -rf $(LEVI_OUT)
	@docker rmi -f $(DOCKER_REPO)triton:latest
	@docker rmi -f $(DOCKER_REPO)triton:$(DOCKER_TAG)
	@for service in $(SERVICE_LIST); do \
	  docker rmi -f $(DOCKER_REPO)$$service:latest; \
	  docker rmi -f $(DOCKER_REPO)$$service:$(DOCKER_TAG); \
	done

## client: builds the triton client
client: echo
ifeq ($(DEBUG),true)
	@echo Serving debug triton client
	@cd $(WORKDIR)/$(TRITON_PATH); npm start
else
	@echo Building triton client release container
	@docker build web/triton --target=release --tag $(DOCKER_REPO)triton:$(DOCKER_TAG)
endif

echo:
	@echo Project directory: $(WORKDIR)

## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

## kraken: builds the kraken base image
kraken: echo
	@$(MAKE) -f $(WORKDIR)/Makefile PROJECT_NAME=kraken template

## leviathan: builds the leviathan binary
leviathan: echo
	@echo Building leviathan server to $(LEVI_OUT)
	@ rm -rf $(LEVI_OUT)
	@go build -o $(LEVI_OUT)leviathan $(LEVI_CMD)main.go
	@cd $(WORKDIR)/$(TRITON_PATH); npm run build
	@cp -r $(TRITON_PATH)dist/triton $(LEVI_OUT)
ifeq ($(DEBUG),true)
	@cp $(LEVI_CMD)leviathan-debug-config.yaml $(LEVI_OUT)leviathan-debug-config.yaml
	./bin/leviathan --config $(LEVI_OUT)leviathan-debug-config.yaml
endif

## proto: generates the services from their proto definitions
proto:
	@echo Generating service code
	@for service in $$(find services/ -mindepth 1 -maxdepth 1 -printf '%f\n'); do \
	  rm -rf $(GEN_GO)$$service; rm -rf $(GEN_TS)$$service; rm -rf $(GEN_DOC); \
	  mkdir $(GEN_GO)$$service; mkdir $(GEN_TS)$$service; mkdir $(GEN_DOC); \
	  protoc -I services/ --go_out=$(GEN_GO) --go_opt=paths=source_relative --twirp_out=$(GEN_GO) --twirp_opt=paths=source_relative \
	    --twirp_typescript_out=version=v6:$(GEN_TS)/$$service \
		--doc_out=$(GEN_DOC) --doc_opt=markdown,$$service.md \
	    $$(find services/ -iname "*.proto"); \
	done

## push: pushes docker images
push:
	@docker push $(DOCKER_REPO)triton:$(DOCKER_TAG)
	@for service in $(SERVICE_LIST); do \
	  docker push $(DOCKER_REPO)$$service:$(DOCKER_TAG); \
	done

template:
ifeq ($(DEBUG),true)
	@echo Building $(PROJECT_NAME) debug container
	@docker build $(WORKDIR) -f cmd/$(PROJECT_NAME)/Dockerfile --target=debug --tag $(PROJECT_NAME):$(DOCKER_TAG)
	@bash -c "trap 'docker kill $(PROJECT_NAME)-debug; docker rm $(PROJECT_NAME)-debug > /dev/null' 0;(docker run -p 8080:8080 -p 8081:8081 -p 40000:40000 --name $(PROJECT_NAME)-debug $(PROJECT_NAME):$(DOCKER_TAG) &) && sleep infinity"
else
	@echo Building $(PROJECT_NAME) release container
	@docker build $(WORKDIR) -f cmd/$(PROJECT_NAME)/Dockerfile --target=release --tag $(DOCKER_REPO)$(PROJECT_NAME):$(DOCKER_TAG)
endif

## version: FIXME
version:
	@echo Version 0.0.1
