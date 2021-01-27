#!/bin/bash
make gyarados DOCKER_REPO=docker.io/redsailtechnologies/bosn- DOCKER_TAG=demo
make poseidon DOCKER_REPO=docker.io/redsailtechnologies/bosn- DOCKER_TAG=demo
make push DOCKER_REPO=docker.io/redsailtechnologies/bosn- DOCKER_TAG=demo