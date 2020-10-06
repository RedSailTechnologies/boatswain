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
* Poseidon
  * Gateway service
* Kraken
  * Cluster management/scraper
* DavyJones
  * Helm management plugin

## TODO
- [x] Create client
- [ ] Create basic api for versions
- [ ] Dockerfile setup
- [ ] Helm charts
- [ ] All-in-one binary
- [ ] Logging
- [ ] Swagger
