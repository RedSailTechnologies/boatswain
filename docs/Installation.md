# Installation

## Helm
### Basics
Add the repo.

`helm repo add redsail https://redsailtechnologies.github.io/helm`

Install

`helm install <name> redsail/boatswain`

### Reasonable Defaults
Below is a small set of reasonable configuration. You can also use `--set-file kraken.config.clusters.kubeConfig=<full kubeconfig path>` in the helm install in conjunction with the contexts setting below to automatically add clusters to the configuration.
```
mate:
  ingress:
    enabled: true
    hosts:
      - host: <your host
        paths: ["/"]

kraken:
  enabled: true
  config:
    clusters:
      contexts:
        # add the contexts you want configured from the kube config file below
        - <your kubeconfig context>

poseidon:
  enabled: true
  config:
    repos:
    # add more helm repos here
    - name: redsail
      endpoint: https://redsailtechnologies.github.io/helm
```
