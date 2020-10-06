.PHONY: client echo help version

# VARS
WORKDIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

# ENV
DOCKER_BUILDKIT=1

# BASIC TARGETS
## client: builds the triton client
client: echo web/*
	@echo "Building triton client"
	@cd $(WORKDIR)/web/triton; npm run build

## echo: prints out the project root
echo:
	@echo Project directory: $(WORKDIR)

## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

kraken:
	@echo Building kraken release container
	@docker build $(WORKDIR) -f cmd/kraken/Dockerfile --target=base

kraken-debug: kraken
	@echo Building kraken debug container
	@docker build $(WORKDIR) -f cmd/kraken/Dockerfile --target=debug --tag kraken:latest
	@bash -c "trap 'docker kill kraken-debug && docker rm kraken-debug > /dev/null' 0;(docker run -p 8080:8080 -p 40000:40000 --name kraken-debug kraken:latest &) && sleep infinity"

kraken-release:
	@docker build $(WORKDIR) -f cmd/kraken/Dockerfile --target=release --tag kraken:latest

## version: FIXME
version:
	@echo Version 0.0.1
