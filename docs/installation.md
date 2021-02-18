# Installation

## Helm
### Basics
`helm repo add redsail https://redsailtechnologies.github.io/helm`

`helm install <name> redsail/boatswain`

### Reasonable Helm Defaults
The services come with reasonable out of the box defaults, but can be customized with the normal helm parameters. If you cannot run a pvc when installing with helm set `mongodb.persistence.enabled` to false. No ingress is configured by default but port 80 for the mate service will expose everything else. See the OIDC section below for how to configure auth.

## Leviathan
### Basics
Leviathan is now a docker container! Pulling `docker.io/redsailtechnologies/bosn-leviathan` will get you the image for any version from 0.6 onward. To run the container the following env variables are required (and fairly self-explanatory): MONGO_CONNECTION_STRING, OIDC_URL (this is the oidc endpoint with /.wellknown/openid-configuration on the end), OIDC_CLIENT (this is just the oidc endpoint), USER_SCOPE, CLIENT_SCOPE, CLIEND_ID. Optionally oidc can be further configured with USER_ADMIN_ROLE (duplicated for EDITOR and READER).

## OIDC
To configure OIDC for the server the following environment variables are used: `OIDC_URL`, `USER_SCOPE`, `USER_ADMIN_ROLE`, `USER_EDITOR_ROLE`, `USER_READER_ROLE`
(or cli params --oidc-url etc.). The client serves a json file found at assets/config/config.prod.json relative to the html root, and values here can be 
overwritten for the client. For helm installations these values can be customized in values.yaml

*IMPORTANT:* The roles used are not part of the default oidc profile and must be setup and added to both the id token and the access token for your oidc provider. They are however optionally changed but are by default set to what is seen below.

Helm Example (using Azure AD and an App Registration):
```yaml
global:
  oidc:
    url: https://login.microsoftonline.com/<azure-directory>/v2.0/
    clientId: <app-registration-client-id>
    clientScopes: "openid profile api://<app-registration-client-id>/boatswain"
    apiScope: "boatswain"
    roles:
      admin: Boatswain.Admin
      editor: Boatswain.Editor
      reader: Boatswain.Reader
```
