package kraken

import (
	"io/ioutil"

	"gopkg.in/yaml.v2"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Cluster is the configuration struct for a single cluster
type Cluster struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Token    string `yaml:"token"`
	Cert     string `yaml:"cert"`
}

// ClusterList is the list of all clusters in the configuration
type ClusterList struct {
	Clusters []Cluster `yaml:"clusters"`
}

// YAML takes a relative filename and returns the config found in it
func (c *ClusterList) YAML(file string) error {
	y, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}
	if err := yaml.UnmarshalStrict(y, c); err != nil {
		return err
	}
	return nil
}

// ToClientset gets the client-go Clientset for this cluster
func (c *Cluster) ToClientset() (*kubernetes.Clientset, error) {
	config := &rest.Config{
		Host:        c.Endpoint,
		BearerToken: c.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(c.Cert),
		},
	}
	return kubernetes.NewForConfig(config)
}
