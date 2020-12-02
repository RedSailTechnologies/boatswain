# Installation

## Helm
### Basics
`helm repo add redsail https://redsailtechnologies.github.io/helm`

`helm install <name> redsail/boatswain`

### Reasonable Helm Defaults
The services come with reasonable out of the box defaults, but can be customized with the normal helm parameters. If you cannot run a pvc when installing with helm set `mongodb.persistence.enabled` to false. No ingress is configured by default but port 80 for the mate service will expose everything else. See the OIDC section below for how to configure auth.

## Leviathan
### Basics
Download and extract leviathan-0.1.0.zip. The leviathan binary and static web content is located there.
A mongodb connection string (and working instance) is required either with
the environment variable `MONGO_CONNECTION_STRING` or the `--mongo-conn` flag.

## OIDC
To configure OIDC for hte server the following environment variables are used: `OIDC_URL`, `USER_SCOPE`, `USER_ADMIN_ROLE`, `USER_EDITOR_ROLE`, `USER_READER_ROLE` (or cli params --oidc-url etc.).
The client serves a json file found at assets/config/config.prod.json relative to the html root, and values here can be overwritten for the client.
For helm installations these values can be customized in values.yaml
*IMPORTANT:* The roles used are not part of the default oidc profile and must be setup and added to both the id token and the access token.
