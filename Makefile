##! Make
.DEFAULT_GOAL=help## Sets the default target for when calling just "make".
SHELL=/bin/bash## Sets the default shell.

##! Common
DEBUG=false## Debug mode or not for build targets.
PROJECT_name=boatswain## Sets the project name.
SERVICE_LIST=gyarados kraken poseidon tentacle## List of all services in the project (excluding the web client).
SERVICE_NAME=## Sets the service to build.
TEST_OUT=## Output type for testing
TRITON_PATH=web/triton/## Triton web client directory.
WORKDIR=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))## The working directory to use, best left as default.

##! Docker/Helm
CHART_LIST=boatswain mate triton $(SERVICE_LIST)## List of charts to build.
DOCKER_BUILDKIT=1## Docker buildkit var.
DOCKER_OPTS=## Extra docker options.
DOCKER_REPO=## The docker repo.
DOCKER_TAG=$(shell git describe --tags --abbrev=0 | cut -d'v' -f2)## The docker tag.
HELM_OUT=bin/## Helm packaging output directory.

##! Code Gen
GEN_DOC=docs/api/## Destination for generated documentation.
GEN_GO=rpc/## Destination for generated go code.
GEN_TS=$(TRITON_PATH)src/app/services/## Destination for generated typescript code.

##! Leviathan
LEVI_CMD=cmd/leviathan/## Leviathan main.go directory.
LEVI_OUT=bin/## Destination for leviathan build.

ifneq (,$(wildcard ./.env))
    include .env
    export
endif

##@ Common

.PHONY: defaults
defaults: ## Display default environment values defined in this makefile.
	@$(info $(shell printf "\033[1mEnvironment Defaults:\033[0m\n"))
	@$(foreach v, $(.VARIABLES), $(if $(filter file,$(origin $(v))), $(info $(v)=$($(v)))))
	@echo

.PHONY: echo
echo: ## Prints the projects working directory.
	@echo Project directory: $(WORKDIR)

.PHONY: env
env: ## Display environment variable documentation.
	@awk 'BEGIN {FS = "=.*##" ;printf "\033[1mEnvironment Variables:\033[0m\n"} \
	/^[a-zA-Z_.]+=.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } \
	/^##!/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

.PHONY: help
help: ## Display makefile target documentation.
	@awk 'BEGIN {FS = ":.*##"; printf "\033[1mUsage:\033[0m\n  make \033[36m<target>\033[0m\n"} \
	/^[a-zA-Z_0-9-]+%?:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } \
	/^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

##@ Build

.PHONY: all
all: docs echo proto gyarados leviathan kraken poseidon tentacle triton ## Build the client and all services.

.PHONY: clean
clean: ## Remove binaries, images, etc.
	@echo Cleaning binaries, coverage, web dist, generated code, docker images, and packaged charts...
	@rm -f main
	@rm -f coverage.out
	@rm -rf $(HELM_OUT)
	@rm -rf $(TRITON_PATH)dist
	@rm -rf $(LEVI_OUT)
	@rm -rf $(GEN_DOC)
	@rm -rf $(GEN_GO)
	@mkdir $(GEN_GO)
	@rm -rf $(GEN_TS)
	@mkdir $(GEN_TS)
	@docker rmi -f $(DOCKER_REPO)triton:latest &> /dev/null
	@docker rmi -f $(DOCKER_REPO)triton:$(DOCKER_TAG) &> /dev/null
	@for service in $(SERVICE_LIST); do \
	  docker rmi -f $$service:latest &> /dev/null; \
	  docker rmi -f $$service:$(DOCKER_TAG) &> /dev/null; \
	  docker rmi -f $(DOCKER_REPO)$$service:latest &> /dev/null; \
	  docker rmi -f $(DOCKER_REPO)$$service:$(DOCKER_TAG) &> /dev/null; \
	done

.PHONY: init
init: echo ## Download support packages (mostly for proto3).
	@go get -u github.com/golang/protobuf/protoc-gen-go
	@go get -u github.com/twitchtv/twirp/protoc-gen-twirp
	@go get -u go.larrymyers.com/protoc-gen-twirp_typescript
	@go get -u github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc
	@go mod tidy

.PHONY: proto
proto: echo ## Generate code based on proto3 definitions.
	@echo Generating service code and docs
	@rm -rf $(GEN_DOC); mkdir $(GEN_DOC); 
	@for service in $$(find services/ -mindepth 1 -maxdepth 1 -printf '%f\n' | cut -d '.' -f1); do \
	  rm -rf $(GEN_GO)$$service; rm -rf $(GEN_TS)$$service; \
	  mkdir $(GEN_GO)$$service; mkdir $(GEN_TS)$$service; \
	  protoc -I services/ --go_out=$(GEN_GO)$$service --go_opt=paths=source_relative --twirp_out=$(GEN_GO)$$service --twirp_opt=paths=source_relative \
	    --twirp_typescript_out=version=v6:$(GEN_TS)$$service \
		--doc_out=$(GEN_DOC) --doc_opt=markdown,$$service.md \
	    $$service.proto; \
	done

##@ Publish

.PHONY: changes
changes: ## Regenerate the CHANGELOG from tags and commit entries.
	@echo "# CHANGELOG" > CHANGELOG.md
	@echo >> CHANGELOG.md
	@for i in $$(git tag | sort -Vr); do \
		[[ "$$i" == "v0.0.1" ]] && continue; \
		echo "## $$i" >> CHANGELOG.md; \
		git log $$(git describe --abbrev=0 --tags $$i^)..$$i --oneline --first-parent | xargs -i echo "*" "{}" >> CHANGELOG.md; \
		echo >> CHANGELOG.md; \
	done
	@tail -n 1 "CHANGELOG.md" | wc -c | xargs -I {} truncate "CHANGELOG.md" -s -{}

.PHONY: docs
docs: proto ## Build all documentation.
	@echo "# Boatswain Api" > docs/api.md
	@echo "Click below for each service's documentation." >> docs/api.md
	@for doc in $$(ls docs/api); do \
	  echo "* [$$(echo $$doc | cut -d '.' -f1)](https://redsailtechnologies.github.io/boatswain/api/$$(echo $$doc | cut -d '.' -f1).html)" >> docs/api.md; \
	done

.PHONY: package
package: echo ## Generate helm packages.
	@echo Packaging charts to $(HELM_OUT)
	@rm -rf $(HELM_OUT)
	@mkdir -p $(HELM_OUT)
	@helm dependency update --skip-refresh deploy/boatswain
	@for chart in $(CHART_LIST); do \
		helm package deploy/$$chart --version $(shell make versionplain) --app-version $(shell make versionplain) --destination $(HELM_OUT); \
	done

.PHONY: push
push: ## Push local images.
	@for service in $(SERVICE_LIST); do \
	  docker push $(DOCKER_REPO)$$service:$(DOCKER_TAG); \
	done
	docker push $(DOCKER_REPO)leviathan:$(DOCKER_TAG) $(DOCKER_OPTS)
	docker push $(DOCKER_REPO)triton:$(DOCKER_TAG) $(DOCKER_OPTS)

.PHONY: version
version: ## Get the current version.
	@git describe --tags --abbrev=0

.PHONY: versionplain
versionplain: ## Get the current version without the leading v.
	@git describe --tags --abbrev=0 | cut -d'v' -f2

.PHONY: versionprev
versionprev: ## Get the previous version.
	@version=$(shell make version);git describe --tags --abbrev=0 --tags $$version^

##@ Services

.PHONY: gyarados
gyarados: echo ## Build the gyarados service.
	@$(MAKE) -f $(WORKDIR)/Makefile SERVICE_NAME=gyarados template

.PHONY: kraken
kraken: echo ## Build the kraken service.
	@$(MAKE) -f $(WORKDIR)/Makefile SERVICE_NAME=kraken template

.PHONY: leviathan
leviathan: echo proto ## Build the leviathan monolith.
	@-$(MAKE) -f $(WORKDIR)/Makefile SERVICE_NAME=leviathan template

.PHONY: poseidon
poseidon: echo ## Build the poseidon service.
	@$(MAKE) -f $(WORKDIR)/Makefile SERVICE_NAME=poseidon template

.PHONY: template
template: ## The build template for all services. If this target is used SERVICE_NAME must be set!
ifeq ($(DEBUG),true)
	@echo Building $(SERVICE_NAME) debug container
	@docker build $(WORKDIR) -f cmd/$(SERVICE_NAME)/Dockerfile --target=debug --tag $(SERVICE_NAME):$(DOCKER_TAG) $(DOCKER_OPTS)
	@bash -c "trap 'docker kill $(SERVICE_NAME)-debug; docker rm $(SERVICE_NAME)-debug > /dev/null' 0;(docker run -p 8080:8080 -p 40000:40000 --name $(SERVICE_NAME)-debug $(SERVICE_NAME):$(DOCKER_TAG) &) && sleep infinity"
else
	@echo Building $(SERVICE_NAME) release container
	@docker build $(WORKDIR) -f cmd/$(SERVICE_NAME)/Dockerfile --target=release --tag $(DOCKER_REPO)$(SERVICE_NAME):$(DOCKER_TAG) $(DOCKER_OPTS)
endif

.PHONY: tentacle
tentacle: echo ## Build the tentacle service.
	@$(MAKE) -f $(WORKDIR)/Makefile SERVICE_NAME=tentacle template

.PHONY: triton
triton: echo ## Build the triton client.
ifeq ($(DEBUG),true)
	@echo Serving debug triton client
	@cd $(WORKDIR)/$(TRITON_PATH); npm start
else
	@echo Building triton client release container
	@docker build web/triton --target=release --tag $(DOCKER_REPO)triton:$(DOCKER_TAG) $(DOCKER_OPTS)
endif

##@ Testing

.PHONY: test
test: echo proto ## Run all unit tests. Set TEST_OUT=html for html coverage report.
	@go test ./pkg/... -cover -coverprofile coverage.out
ifeq ($(TEST_OUT),html)
	@go tool cover -html=coverage.out
endif
ifeq ($(TEST_OUT),func)
	@go tool cover -func=coverage.out
endif
