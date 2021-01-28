#!/bin/bash
set -e
make all DOCKER_REPO=docker.io/redsailtechnologies/bosn- DOCKER_TAG=demo
make push DOCKER_REPO=docker.io/redsailtechnologies/bosn- DOCKER_TAG=demo