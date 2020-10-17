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
  * Agnostic about other technologies used in conjuction with Kubernetes.

## Components
* Triton
  * Angular client
* Mate
  * Envoy proxy
* Kraken
  * Cluster management/scraper
* Poseidon
  * Helm management plugin
* Leviathan
  * Single binary version
  
## Roadmap
* 0.2
  * Multiple ways to add values to an upgrade
  * Install button/functionality
  * Ways to add a cluster and repo (in client)
  * Better configuration management
* 0.3
  * Persistent storage/service cache management
  * Smart default settings for upgrades
  * Client side filtering/dynamic page updates
* 0.4
  * Users and authentication
  * Role based authorization
  * Client state
* 0.5
  * Flow management
  * Strategies
  * Triggers
