# Boatswain (bow·sn)
![Test](https://github.com/RedSailTechnologies/boatswain/workflows/Test/badge.svg?branch=main)
![Develop](https://github.com/RedSailTechnologies/boatswain/workflows/Develop/badge.svg)

![GoReport](https://goreportcard.com/badge/github.com/redsailtechnologies/boatswain)

## About
A kubernetes-native CD tool with first-class support for Helm v3.

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
* 0.3
  * Users/auth (LDAP)
  * RBAC
  * Client login/cookies
  * Service health checks
* 0.4
  * Dev/debug canaries
  * mongodb security (for helm)
* 0.5
  * Cthulu first cut
    * azdo plugin
    * triggers
    * PR canaries
