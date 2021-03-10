package kube

import (
	"encoding/json"

	appsv1 "k8s.io/api/apps/v1"
)

// Args are arguments common to all commands
type Args struct {
	Labels    *map[string]string
	Namespace string
}

// ConvertArgs converts byte data back to args
func ConvertArgs(data []byte) (*Args, error) {
	args := Args{}
	err := json.Unmarshal(data, &args)
	if err != nil {
		return nil, err
	}
	return &args, nil
}

// Result is the type this agent returns
type Result struct {
	Data interface{}
	Type ResultType
}

// ResultType represents types to convert for this response
type ResultType string

const (
	// DeploymentsResult represents a result of Deployments
	DeploymentsResult ResultType = "DeploymentsResult"

	// ReleaseNameResult represents a helm release name lookup
	ReleaseNameResult ResultType = "ReleaseName"

	// StatefulSetsResult represents a result of StatefulSets
	StatefulSetsResult ResultType = "StatefulSetsResult"

	// StatusResult represents a Status result
	StatusResult ResultType = "StatusResult"
)

// ResultTypeError is an eror when the wrong type conversion is made with a result
type ResultTypeError struct{}

func (err ResultTypeError) Error() string {
	return "incorrect kube result type"
}

// ReleaseNotFoundError represents an error when searching for a result name but none is found
type ReleaseNotFoundError struct{}

func (err ReleaseNotFoundError) Error() string {
	return "release not found with given label(s)"
}

// ConvertDeployments returns the Deployments from this result
func ConvertDeployments(data []byte) ([]appsv1.Deployment, error) {
	result, err := unmarshalResult(data)
	if err != nil {
		return nil, err
	}

	if result.Type != DeploymentsResult {
		return nil, ResultTypeError{}
	}

	deployments := make([]appsv1.Deployment, 0)
	for _, dep := range result.Data.([]interface{}) {
		b, err := json.Marshal(dep)
		if err != nil {
			return nil, err
		}

		deployment := appsv1.Deployment{}
		err = json.Unmarshal(b, &deployment)
		if err != nil {
			return nil, err
		}

		deployments = append(deployments, deployment)
	}
	return deployments, nil
}

// ConvertReleaseName returns the ReleaseName from this result
func ConvertReleaseName(data []byte) (string, error) {
	result, err := unmarshalResult(data)
	if err != nil {
		return "", err
	}

	if result.Type != DeploymentsResult {
		return "", ResultTypeError{}
	}

	return result.Data.(string), nil
}

// ConvertStatefulSet returns the StatefulSets from this result
func ConvertStatefulSet(data []byte) ([]appsv1.StatefulSet, error) {
	result, err := unmarshalResult(data)
	if err != nil {
		return nil, err
	}

	if result.Type != StatefulSetsResult {
		return nil, ResultTypeError{}
	}

	sets := make([]appsv1.StatefulSet, 0)
	for _, set := range result.Data.([]interface{}) {
		b, err := json.Marshal(set)
		if err != nil {
			return nil, err
		}

		set := appsv1.StatefulSet{}
		err = json.Unmarshal(b, &set)
		if err != nil {
			return nil, err
		}

		sets = append(sets, set)
	}
	return sets, nil
}

// ConvertStatus returns the status from this result
func ConvertStatus(data []byte) (bool, error) {
	result, err := unmarshalResult(data)
	if err != nil {
		return false, err
	}

	if result.Type != StatusResult {
		return false, ResultTypeError{}
	}
	return result.Data.(bool), nil
}

func unmarshalResult(data []byte) (*Result, error) {
	result := Result{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
