package kraken

import (
	"errors"
	"fmt"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
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
func (c *Config) ToHelmClient(clusterName string, namespace string) (*action.Configuration, error) {
	cluster, err := c.getClusterConfig(clusterName)
	if err != nil {
		return nil, err
	}

	flags := &genericclioptions.ConfigFlags{
		APIServer:   &cluster.Endpoint,
		BearerToken: &cluster.Token,
		// TODO AdamP - flags only supports cert files, how do we want to handle?
		// CertFile:    &cluster.Cert,
		Insecure: &[]bool{true}[0],
	}
	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(flags, namespace, "secrets", helmLogger); err != nil {
		return nil, err
	}

	return actionConfig, nil
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

func (c *Config) getClusterConfig(clusterName string) (*ClusterConfig, error) {
	for _, config := range c.Clusters {
		if config.Name == clusterName {
			return &config, nil
		}
	}
	return nil, errors.New("cluster not found")
}

func helmLogger(template string, args ...interface{}) {
	logger.Info(fmt.Sprintf(template, args...))
}
