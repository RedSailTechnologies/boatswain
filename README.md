# Boatswain (bow·sn)

## About
A kubernetes-native CD tool.

## Design
* Kubernetes native technologies
  * Angular frontend
  * Go backend
* First Class Helm Support
  * Maintain an easy to use helm chart with high configurability when desired.
  * Out of the box experience should be simple.
* Multiple Run Configurations
  * Run locally/as a single docker container.
  * Run in Kubernetes as a collection of microservices with a single entrypoint.
  * Agnostic about other technologies used in conjuction with Kubernetes.

## Components
* Triton
  * Angular client
* DavyJones/Mate
  * Gateway/reverse proxy (Envoy?)
* Kraken
  * Cluster management/scraper
* Poseidon
  * Helm management plugin
* Leviathan
  * Single binary version
* Gyarados (planned)
  * Blue/green and Canary plugin
* Siren (planned)
  * Github actions management plugin
* Cthulhu (planned)
  * Azdo management plugin

## TODO for 0.1
- [x] Create client
- [x] Create basic api for versions
- [x] Dockerfile setup
- [x] Helm charts
- [x] Envoy proxy
- [x] Logging
- [x] API Docs
- [x] All-in-one binary
- [ ] Helm integration
- [ ] Update docs
- [ ] Plan Future Features
  - [ ] Configuration for leviathan/development
  - [ ] Filter/organize by NS
  - [ ] Refresh button for a cluster
  - [ ] Save user presets in a cookie
  - [ ] Users/Auditing
  - [ ] Pipelines/tasks/triggers
  - [ ] Acceptance critera
