# Installation

## Helm
### Basics
`helm repo add redsail https://redsailtechnologies.github.io/helm`

`helm install <name> redsail/boatswain`

### Reasonable Defaults
Entities can be added upon installation, so no default configuration is needed. If you do want to have things pre-populated, below is a small set of reasonable configurations. You can also use `--set-file kraken.config.clusters.kubeConfig=<full kubeconfig path>` in the helm install in conjunction with the contexts setting below to automatically add clusters to the configuration.
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
        # add the contexts you want configured from the kube config file above
        - <your kubeconfig context>

poseidon:
  enabled: true
  config:
    repos:
    # add more helm repos here
    - name: redsail
      endpoint: https://redsailtechnologies.github.io/helm
  cacheDir: ./temp
```

## Leviathan
Download and extract leviathan-0.1.0.zip. The leviathan binary and static web content is located there. Running leviathan requires no special configuration unless you want to specify configuration as above. Note for leviathan everything under config for the services above is all at the top level of the yaml.