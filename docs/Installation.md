# Installation

## Helm
### Basics
`helm repo add redsail https://redsailtechnologies.github.io/helm`

`helm install <name> redsail/boatswain`

### Reasonable Helm Defaults
The services come with reasonable out of the box defaults, but can be customized with the normal helm parameters. If you cannot run a pvc when installing with helm set `mongodb.persistence.enabled` to false. No ingress is configured but port 80 for the mate service will expose everything else.

## Leviathan
### Basics
Download and extract leviathan-0.1.0.zip. The leviathan binary and static web content is located there. A mongodb connection string (and working instance) is required either with
the environment variable `MONGO_CONNECTION_STRING` or the `--mongo-conn` flag.
