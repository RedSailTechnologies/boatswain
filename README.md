# Boatswain (bow·sn)
![Develop](https://github.com/RedSailTechnologies/boatswain/workflows/Develop/badge.svg)
![Release](https://github.com/RedSailTechnologies/boatswain/workflows/Release/badge.svg)

## About
A kubernetes-native CD tool with first-class support for Helm v3.

## Design
* Kubernetes native technologies
  * Angular frontend
  * Go backend
* First Class Helm Support
  * Easy to use helm chart with high configurability when desired.
  * Simple out of the box experience.
* Multiple Run Configurations
  * Run locally/as a single binary.
  * Run in Kubernetes as a collection of microservices with a single entrypoint.
* Canary strategy as first class object.
  * First class istio support
  * Debug canaries
  * Pre-merge canaries
  * Upgrades done with canaries
* SCM Plugins
  * AZDO integration for PRs
  * Github actions
  * Webhook calls/repository triggers
* CD tools/flows
  * Canary upgrades/testing
  * Automatic environment promotion based on conditions
  * Manual promotion when desired

## Components
* Triton
  * Angular client
  * Monitoring
  * Dev debug setup
* Mate
  * Envoy proxy/routing
* Leviathan
  * Single binary version
* Kraken
  * Kube/cluster management
  * Helm commands
  * Split off/cqrs?         ##########
* Poseidon
  * Repo plugin
  * Helm
  * Docker
* Gyrados
  * Canary/CD flow management
  * Istio management?       ##########
  * Test running?           ##########
* Cthulu
  * scm plugins
    * azdo
    * github
  * triggers
    * web calls
    * ???
  
## Questions/Tradeoffs
* Do we manage the istio/virtual service side of things?
* Or do we let helm do it and we just allow for the steps/flow?
  * What do we consider steps/flow?
    * Templates?
    * Steps/stages etc?
    * Inputs from existing deployment/cluster-how to work in?
  * Install a helm chart and upgrade another? Upgrade a single chart?
    * Do we want this to be configurable, or just one chart?
  * What options do we have and what opinions do we want?
    * We want values from a library file.
    * What about manually specified?
    * File upload?
* How do we want to do test running?
  * Helm tests?
  * Another chart?
  * Run a container? (which we probably make into a pod?)
  * What to do with results?

## Roadmap
* 0.2
  * Config/client cleanup and performance optimization
  * Add cluster/repo functionality (through client)
  * Gyrados initial cut and design
    * Dev canaries
    * Deployment canaries
* 0.3
  * Persistent storage
  * Client side filtering/part-of views
* 0.4
  * Users/auth (LBAC)
  * RBAC
  * Client login/cookies
* 0.5
  * Cthulu first cut
    * azdo plugin
    * triggers
    * PR canaries
