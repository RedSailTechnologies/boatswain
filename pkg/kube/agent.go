package kube

import (
	"context"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/kubernetes"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
)

const (
	helmAnnotation string = "meta.helm.sh/release-name"
)

// Agent is used for communication with Kubernetes, made into an interface for testability
type Agent interface {
	GetDeployments(args *Args) (*Result, error)
	GetReleaseName(args *Args) (*Result, error)
	GetStatefulSets(args *Args) (*Result, error)
	GetStatus(args *Args) (*Result, error)
}

// AgentAction is an enum/alias to make calling methods typed
type AgentAction string

const (
	// GetDeployments represents the GetDeployments method
	GetDeployments AgentAction = "GetDeployments"

	// GetReleaseName represents the GetReleaseName method
	GetReleaseName AgentAction = "GetReleaseName"

	// GetStatefulSets represents the GetStatefulSets method
	GetStatefulSets AgentAction = "GetStatefulSets"

	// GetStatus represents the GetStatus method
	GetStatus AgentAction = "GetStatus"
)

// DefaultAgent is the default implementation of the KubeAgent
type DefaultAgent struct {
	kubeFunc func() (kubernetes.Interface, error)
}

// NewDefaultAgent inits the default agent with the specified kube interface
func NewDefaultAgent(kubeFunc func() (kubernetes.Interface, error)) *DefaultAgent {
	return &DefaultAgent{
		kubeFunc: kubeFunc,
	}
}

// GetDeployments gets all the deployments for a particular cluster
func (k DefaultAgent) GetDeployments(args *Args) (*Result, error) {
	k8s, err := k.kubeFunc()
	if err != nil {
		return nil, err
	}

	d, err := k8s.AppsV1().Deployments(args.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromValidatedSet(*args.Labels).String(),
	})
	if err != nil {
		logger.Error("could not get deployments from cluster", "error", err)
		return nil, err
	}

	return &Result{
		Data: d.Items,
		Type: DeploymentsResult,
	}, nil
}

// GetReleaseName name gets a helm release name from the annotation meta.helm.sh/release-name based on label selectors
func (k DefaultAgent) GetReleaseName(args *Args) (*Result, error) {
	k8s, err := k.kubeFunc()
	if err != nil {
		return nil, err
	}

	deps, err := k8s.AppsV1().Deployments(args.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromValidatedSet(*args.Labels).String(),
	})
	if err != nil {
		return nil, err
	}
	for _, d := range deps.Items {
		if val, ok := d.Annotations[helmAnnotation]; ok {
			return &Result{
				Data: val,
			}, nil
		}
	}

	sets, err := k8s.AppsV1().StatefulSets(args.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromValidatedSet(*args.Labels).String(),
	})
	if err != nil {
		return nil, err
	}
	for _, s := range sets.Items {
		if val, ok := s.Annotations[helmAnnotation]; ok {
			return &Result{
				Data: val,
			}, nil
		}
	}

	return nil, ReleaseNotFoundError{}
}

// GetStatefulSets gets all the statefulsets for a particular cluster
func (k DefaultAgent) GetStatefulSets(args *Args) (*Result, error) {
	k8s, err := k.kubeFunc()
	if err != nil {
		return nil, err
	}

	ss, err := k8s.AppsV1().StatefulSets(args.Namespace).List(context.TODO(), metav1.ListOptions{
		LabelSelector: labels.SelectorFromValidatedSet(*args.Labels).String(),
	})
	if err != nil {
		logger.Error("could not get statefulsets from cluster", "error", err)
		return nil, err
	}

	return &Result{
		Data: ss.Items,
		Type: StatefulSetsResult,
	}, nil
}

// GetStatus returns the status of a cluster by ensuring each node is in a ready state
func (k DefaultAgent) GetStatus(args *Args) (*Result, error) {
	k8s, err := k.kubeFunc()
	if err != nil {
		return nil, err
	}

	nodes, err := k8s.CoreV1().Nodes().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		logger.Error("could not get nodes from cluster", "error", err)
		return nil, err
	}

	for _, node := range nodes.Items {
		for _, condition := range node.Status.Conditions {
			if condition.Type == "Ready" {
				if condition.Status != "True" {
					return &Result{
						Data: false,
						Type: StatusResult,
					}, nil
				}
			}
		}
	}

	return &Result{
		Data: true,
		Type: StatusResult,
	}, nil
}
