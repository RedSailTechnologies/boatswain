# Home

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