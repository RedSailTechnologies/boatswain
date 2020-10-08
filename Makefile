.PHONY: client echo help kraken target version

# ENV
DEBUG=false
DOCKER_BUILDKIT=1
GEN_GO=pkg/
GEN_TS=web/triton/src/app/services/
PROJECT_NAME=null
SERVICE_LIST=kraken
WORKDIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# BASIC TARGETS
## client: builds the triton client
client: echo web/*
	@echo "Building triton client"
	@cd $(WORKDIR)/web/triton; npm start

## echo: prints out the project root
echo:
	@echo Project directory: $(WORKDIR)

## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

## kraken: builds the kraken base image
kraken: echo
	@$(MAKE) -f $(WORKDIR)/Makefile PROJECT_NAME=kraken template

## proto: generates the services from their proto definitions
proto:
	@echo Generating service code
	@protoc -I services/ --go_out=$(GEN_GO) --go_opt=paths=source_relative --go-grpc_out=$(GEN_GO) --go-grpc_opt=paths=source_relative \
	 --plugin="protoc-gen-ts=web/triton/node_modules/.bin/protoc-gen-ts" --js_out="import_style=commonjs,binary:$(GEN_TS)" --ts_out="service=grpc-web:$(GEN_TS)" \
	 $$(find services/ -iname "*.proto");

template:
ifeq ($(DEBUG),true)
	@echo Building $(PROJECT_NAME) debug container
	@docker build $(WORKDIR) -f cmd/$(PROJECT_NAME)/Dockerfile --target=debug --tag $(PROJECT_NAME):latest
	@bash -c "trap 'docker kill $(PROJECT_NAME)-debug; docker rm $(PROJECT_NAME)-debug > /dev/null' 0;(docker run -p 8080:8080 -p 8081:8081 -p 40000:40000 --name $(PROJECT_NAME)-debug $(PROJECT_NAME):latest &) && sleep infinity"
else
	@echo Building $(PROJECT_NAME) release container
	@docker build $(WORKDIR) -f cmd/$(PROJECT_NAME)/Dockerfile --target=release --tag $(PROJECT_NAME):latest
endif

## version: FIXME
version:
	@echo Version 0.0.1
