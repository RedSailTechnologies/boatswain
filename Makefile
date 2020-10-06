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
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

## version: FIXME
version:
	@echo Version 0.0.1
