.PHONY: all clean changes docs echo help init gyarados kraken leviathan package poseidon proto push template tentacle test triton version versionprev

# Build env
DEBUG=false
CHART_LIST=boatswain mate triton $(SERVICE_LIST)
DOCKER_BUILDKIT=1
DOCKER_OPTS=
DOCKER_REPO=
DOCKER_TAG=$$(git describe --tags --abbrev=0 | cut -d'v' -f2)
GEN_DOC=docs/api/
GEN_GO=rpc/
GEN_TS=$(TRITON_PATH)src/app/services/
HELM_OUT=bin/
LEVI_CMD=cmd/leviathan/
LEVI_OUT=bin/
PROJECT_NAME=null
SERVICE_LIST=gyarados kraken poseidon tentacle
SHELL=/bin/bash
TRITON_PATH=web/triton/
TEST_OUT=
WORKDIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

# BASIC TARGETS
## all: builds the client and all services
all: docs echo proto gyarados leviathan kraken poseidon tentacle triton

## clean: removes binaries, images, etc
clean:
	@echo Running clean
	@rm -f main
	@rm -f coverage.out
	@rm -rf $(TRITON_PATH)dist
	@rm -rf $(LEVI_OUT)
	@rm -rf $(GEN_DOC)
	@rm -rf $(GEN_GO)
	@mkdir $(GEN_GO)
	@rm -rf $(GEN_TS)
	@mkdir $(GEN_TS)
	@docker rmi -f $(DOCKER_REPO)triton:latest
	@docker rmi -f $(DOCKER_REPO)triton:$(DOCKER_TAG)
	@for service in $(SERVICE_LIST); do \
	  docker rmi -f $(DOCKER_REPO)$$service:latest; \
	  docker rmi -f $(DOCKER_REPO)$$service:$(DOCKER_TAG); \
	done

## changes: regenerates the CHANGELOG from tags and commit entries
changes:
	@echo "# CHANGELOG" > CHANGELOG.md
	@echo >> CHANGELOG.md
	@for i in $$(git tag | sort -Vr); do \
		[[ "$$i" == "v0.0.1" ]] && continue; \
		echo "## $$i" >> CHANGELOG.md; \
		git log $$(git describe --abbrev=0 --tags $$i^)..$$i --oneline --first-parent | xargs -i echo "*" "{}" >> CHANGELOG.md; \
		echo >> CHANGELOG.md; \
	done
	@tail -n 1 "CHANGELOG.md" | wc -c | xargs -I {} truncate "CHANGELOG.md" -s -{}

## docs: builds the documentation into docs/
docs: proto
	@echo "# Boatswain Api" > docs/api.md
	@echo "Click below for each service's documentation." >> docs/api.md
	@for doc in $$(ls docs/api); do \
	  echo "* [$$(echo $$doc | cut -d '.' -f1)](https://redsailtechnologies.github.io/boatswain/api/$$(echo $$doc | cut -d '.' -f1).html)" >> docs/api.md; \
	done

echo:
	@echo Project directory: $(WORKDIR)

## help: prints this help message
help:
	@echo "Usage:"
	@sed -n 's/^##//p' $(MAKEFILE_LIST) | column -t -s ':' |  sed -e 's/^/ /'

## init: downloads support packages (mostly for proto3)
init: echo
	@go get github.com/golang/protobuf/protoc-gen-go
	@go get github.com/twitchtv/twirp/protoc-gen-twirp
	@go get -u go.larrymyers.com/protoc-gen-twirp_typescript
	@go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc

## gyarados: builds the kraken image
gyarados: echo
	@$(MAKE) -f $(WORKDIR)/Makefile PROJECT_NAME=gyarados template

## kraken: builds the kraken image
kraken: echo
	@$(MAKE) -f $(WORKDIR)/Makefile PROJECT_NAME=kraken template

## leviathan: builds the leviathan binary
leviathan: echo proto
	@-$(MAKE) -f $(WORKDIR)/Makefile PROJECT_NAME=leviathan template

## package: generates helm packages
package: echo
	@echo Packaging charts to $(HELM_OUT)
	@mkdir -p $(HELM_OUT)
	@helm dependency update --skip-refresh deploy/boatswain
	@for chart in $(CHART_LIST); do \
		helm package deploy/$$chart --version $(shell make version) --app-version $(shell make version) --destination $(HELM_OUT); \
	done

## poseidon: builds the poseidon image
poseidon: echo
	@$(MAKE) -f $(WORKDIR)/Makefile PROJECT_NAME=poseidon template

## proto: generates the services from their proto definitions
proto: echo
	@echo Generating service code and docs
	@rm -rf $(GEN_DOC); mkdir $(GEN_DOC); 
	@for service in $$(find services/ -mindepth 1 -maxdepth 1 -printf '%f\n' | cut -d '.' -f1); do \
	  rm -rf $(GEN_GO)$$service; rm -rf $(GEN_TS)$$service; \
	  mkdir $(GEN_GO)$$service; mkdir $(GEN_TS)$$service; \
	  protoc -I services/ --go_out=$(GEN_GO)$$service --go_opt=paths=source_relative --twirp_out=$(GEN_GO)$$service --twirp_opt=paths=source_relative \
	    --twirp_typescript_out=version=v6:$(GEN_TS)/$$service \
		--doc_out=$(GEN_DOC) --doc_opt=markdown,$$service.md \
	    $$service.proto; \
	done

## push: pushes local images
push:
	@for service in $(SERVICE_LIST); do \
	  docker push $(DOCKER_REPO)$$service:$(DOCKER_TAG); \
	done
	docker push $(DOCKER_REPO)leviathan:$(DOCKER_TAG) $(DOCKER_OPTS)
	docker push $(DOCKER_REPO)triton:$(DOCKER_TAG) $(DOCKER_OPTS)

template:
ifeq ($(DEBUG),true)
	@echo Building $(PROJECT_NAME) debug container
	@docker build $(WORKDIR) -f cmd/$(PROJECT_NAME)/Dockerfile --target=debug --tag $(PROJECT_NAME):$(DOCKER_TAG) $(DOCKER_OPTS)
	@bash -c "trap 'docker kill $(PROJECT_NAME)-debug; docker rm $(PROJECT_NAME)-debug > /dev/null' 0;(docker run -p 8080:8080 -p 40000:40000 --name $(PROJECT_NAME)-debug $(PROJECT_NAME):$(DOCKER_TAG) &) && sleep infinity"
else
	@echo Building $(PROJECT_NAME) release container
	@docker build $(WORKDIR) -f cmd/$(PROJECT_NAME)/Dockerfile --target=release --tag $(DOCKER_REPO)$(PROJECT_NAME):$(DOCKER_TAG) $(DOCKER_OPTS)
endif

## tentacle: builds the tentacle image
tentacle: echo
	@$(MAKE) -f $(WORKDIR)/Makefile PROJECT_NAME=tentacle template

## test: runs all unit tests, set TEST_OUT=html for html coverage report
test: echo proto
	@go test ./pkg/... -cover -coverprofile coverage.out
ifeq ($(TEST_OUT),html)
	@go tool cover -html=coverage.out
endif
ifeq ($(TEST_OUT),func)
	@go tool cover -func=coverage.out
endif

## triton: builds the triton client
triton: echo
ifeq ($(DEBUG),true)
	@echo Serving debug triton client
	@cd $(WORKDIR)/$(TRITON_PATH); npm start
else
	@echo Building triton client release container
	@docker build web/triton --target=release --tag $(DOCKER_REPO)triton:$(DOCKER_TAG) $(DOCKER_OPTS)
endif

## version: gets the current version of the repo
version:
	@git describe --tags --abbrev=0

versionprev:
	@version=$(shell make version);git describe --tags --abbrev=0 --tags $$version^
