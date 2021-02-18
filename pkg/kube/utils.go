package kube

import "k8s.io/client-go/rest"

// GetInClusterConfig gets the rest config for this cluster
func GetInClusterConfig() (*rest.Config, error) {
	return rest.InClusterConfig()
}
