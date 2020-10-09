package kraken

import (
	"errors"
	"io/ioutil"

	"gopkg.in/yaml.v2"
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

// GetClusterConfig gets the cluster config given the cluster name
func (c *Config) GetClusterConfig(cluster string) (*ClusterConfig, error) {
	for _, config := range c.Clusters {
		if config.Name == cluster {
			return &config, nil
		}
	}
	return nil, errors.New("cluster not found")
}

// ToClientset gets the client-go Clientset for this cluster
func (c *ClusterConfig) ToClientset() (*kubernetes.Clientset, error) {
	config := &rest.Config{
		Host:        c.Endpoint,
		BearerToken: c.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(c.Cert),
		},
	}
	return kubernetes.NewForConfig(config)
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
