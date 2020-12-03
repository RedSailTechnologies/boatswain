# Boatswain (bow·sn)
![Test](https://github.com/RedSailTechnologies/boatswain/workflows/Test/badge.svg?branch=main)
![Develop](https://github.com/RedSailTechnologies/boatswain/workflows/Develop/badge.svg)

![GoReport](https://goreportcard.com/badge/github.com/redsailtechnologies/boatswain)

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
  * Observability/Run logs
  * Dev debug setup
* Mate
  * Envoy proxy/routing
* Leviathan
  * Single binary version
* Kraken
  * Kube management
  * Kube monitoring
* Poseidon
  * Repo plugin
    * Helm
    * Docker
* Gyrados
  * Canary/CD flow management
  * Deployment running
  * Test running
* Cthulu
  * Azure DevOps SCM plugin
  * Github SCM plugin

## Roadmap
* 0.4
  * Service health checks
  * Dev/debug canaries
  * mongodb security (for helm)
* 0.5
  * Cthulu first cut
    * azdo plugin
    * triggers
    * PR canaries
* 0.6
  * Change logs/history
