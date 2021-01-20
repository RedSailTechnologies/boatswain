#!/bin/bash
set -e
# helm uninstall -n bosn bosn
make package
make all DOCKER_REPO=localhost:32000/bosn/
make push DOCKER_REPO=localhost:32000/bosn/
helm install bosn -n bosn deploy/boatswain/ -f temp.yaml 