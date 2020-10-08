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
* Gyarados
  * Blue/green and Canary plugin
* Siren
  * Github (actions) management plugin
* Cthulhu
  * Azdo management plugin
* Leviathan
  * Single binary version

## TODO
- [x] Create client
- [x] Create basic api for versions
- [ ] Dockerfile setup
- [ ] Envoy proxy
- [ ] Helm charts
- [ ] All-in-one binary
- [ ] Logging
- [ ] Pprof
- [ ] Docs/API Docs
