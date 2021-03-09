package main

import (
	"github.com/redsailtechnologies/boatswain/pkg/kube"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

func k8sConfig() (kubernetes.Interface, error) {
	config, err := kube.GetInClusterConfig()
	if err != nil {
		return nil, err
	}
	k8s, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}
	return k8s, nil
}

func restConfig() (*rest.Config, error) {
	config, err := kube.GetInClusterConfig()
	if err != nil {
		return nil, err
	}
	return config, nil
}
