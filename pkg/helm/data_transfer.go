package helm

import (
	"encoding/json"

	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
)

// Args are arguments common to all commands
type Args struct {
	Name      string
	Namespace string
	Chart     *chart.Chart
	Values    map[string]interface{}
	Version   int
	Wait      bool
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
	Logs string
	Type ResultType
}

// ResultType represents types to convert for this response
type ResultType string

const (
	// NoneResult represents helm commands that don't return anything, such as rollback
	NoneResult ResultType = "NoneResult"

	// ReleaseResult represents a Release result
	ReleaseResult ResultType = "ReleaseResult"

	// UninstallReleaseResponseResult represents an UninstallReleaseResponse result
	UninstallReleaseResponseResult ResultType = "UninstallReleaseResponse"
)

// ResultTypeError is an eror when the wrong type conversion is made with a result
type ResultTypeError struct{}

func (err ResultTypeError) Error() string {
	return "incorrect helm result type"
}

// ConvertNone converts the marshaled data into no type in order to return logs
func ConvertNone(data []byte) (string, error) {
	result, err := unmarshalResult(data)
	if err != nil {
		return "", err
	}

	if result.Type != NoneResult {
		return result.Logs, ResultTypeError{}
	}
	return "", nil
}

// ConvertRelease converts the marshaled data to a typed release
func ConvertRelease(data []byte) (*release.Release, string, error) {
	result, err := unmarshalResult(data)
	if err != nil {
		return nil, "", err
	}

	if result.Type != ReleaseResult {
		return nil, result.Logs, ResultTypeError{}
	}
	return result.Data.(*release.Release), result.Logs, nil
}

// ConvertUninstallReleaseResponse converts the marshaled data to a typed uninstall response
func ConvertUninstallReleaseResponse(data []byte) (*release.UninstallReleaseResponse, string, error) {
	result, err := unmarshalResult(data)
	if err != nil {
		return nil, "", err
	}

	if result.Type != UninstallReleaseResponseResult {
		return nil, result.Logs, ResultTypeError{}
	}
	return result.Data.(*release.UninstallReleaseResponse), result.Logs, nil
}

func unmarshalResult(data []byte) (*Result, error) {
	result := Result{}
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, err
	}
	return &result, nil
}
