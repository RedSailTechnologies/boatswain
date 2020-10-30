package kraken

import (
	"fmt"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
	pb "github.com/redsailtechnologies/boatswain/rpc/kraken"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

// Cluster is a type def to include the protobuf definition so we can extend it
type Cluster struct {
	*pb.Cluster
}

// ToHelmClient gets a helm action configuration given the cluster's name
func (c *Cluster) ToHelmClient(namespace string) (*action.Configuration, error) {
	flags := &genericclioptions.ConfigFlags{
		APIServer:   &c.Endpoint,
		BearerToken: &c.Token,
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
func (c *Cluster) ToClientset() (*kubernetes.Clientset, error) {
	restConfig := &rest.Config{
		Host:        c.Endpoint,
		BearerToken: c.Token,
		TLSClientConfig: rest.TLSClientConfig{
			CAData: []byte(c.Cert),
		},
	}
	return kubernetes.NewForConfig(restConfig)
}

func helmLogger(template string, args ...interface{}) {
	logger.Info(fmt.Sprintf(template, args...))
}
