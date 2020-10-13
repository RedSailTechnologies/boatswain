package kraken

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// ClusterConfig is the configuration struct for a single cluster
type ClusterConfig struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Token    string `yaml:"token"`
	Cert     string `yaml:"cert"`
}

// Config is a list of configurations
type Config struct {
	Clusters []ClusterConfig `yaml:"clusters"`
}

// ToHelmClient gets a helm action configuration given the cluster's name
func (c *Config) ToHelmClient(clusterName string) (*action.Configuration, error) {
	cluster, err := c.getClusterConfig(clusterName)
	if err != nil {
		return nil, err
	}

	flags := &genericclioptions.ConfigFlags{
		APIServer:   &cluster.Endpoint,
		BearerToken: &cluster.Token,
		KeyFile:     &cluster.Cert,
	}
	config := new(action.Configuration)
	if err := config.Init(flags, "", "secrets", nil); err != nil {
		return nil, err
	}
	return config, nil
}

// ToClientset gets the client-go Clientset for this cluster given the cluster's name
func (c *Config) ToClientset(clusterName string) (*kubernetes.Clientset, error) {
	cluster, err := c.getClusterConfig(clusterName)
	if err != nil {
		return nil, err
	}

	restConfig := &rest.Config{
		Host:        cluster.Endpoint,
		BearerToken: cluster.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(cluster.Cert),
		},
	}
	return kubernetes.NewForConfig(restConfig)
}

// YAML takes a relative filename and returns the config found in it
func (c *Config) YAML(file string) error {
	y, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err := yaml.UnmarshalStrict(y, c); err != nil {
		return err
	}
	return nil
}

func (c *Config) getClusterConfig(clusterName string) (*ClusterConfig, error) {
	for _, config := range c.Clusters {
		if config.Name == clusterName {
			return &config, nil
		}
	}
	return nil, errors.New("cluster not found")
}
