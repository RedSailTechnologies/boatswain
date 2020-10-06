# VARS
WORKDIR=$(shell pwd)

# ENV
DOCKER_BUILDKIT=1

# BASIC TARGETS
client:
	cd $(WORKDIR)/web/triton && npm run build