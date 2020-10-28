# Boatswain (bow·sn)
![Test](https://github.com/RedSailTechnologies/boatswain/workflows/Test/badge.svg?branch=main)
![Develop](https://github.com/RedSailTechnologies/boatswain/workflows/Develop/badge.svg)

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
  * Kube/cluster management
  * Helm commands
* Poseidon
  * Repo plugin
    * Helm
    * Docker
    * Git
* Gyrados
  * Canary/CD flow management
  * Test running
* Cthulu
  * scm plugins
    * azdo
    * github
  * triggers

## Roadmap
* 0.2
  * Config/client cleanup and performance optimization
  * Add cluster/repo functionality (through client)
  * Gyrados initial cut and design
    * Run/deployment objects
    * Canary strategies/settings
* 0.3
  * Persistent storage
  * Client side filtering/part-of views
* 0.4
  * Users/auth (LDAP)
  * RBAC
  * Client login/cookies
* 0.5
  * Cthulu first cut
    * azdo plugin
    * triggers
    * PR canaries
* 0.6
  * Dev/debug canaries
