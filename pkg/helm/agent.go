package helm

import (
	"bytes"
	"fmt"

	"github.com/redsailtechnologies/boatswain/pkg/logger"
	"helm.sh/helm/v3/pkg/action"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
)

// Agent is the interface we use to talk to helm packages
type Agent interface {
	Install(args *Args) (*Result, error)
	Rollback(args *Args) (*Result, error)
	Test(args *Args) (*Result, error)
	Uninstall(args *Args) (*Result, error)
	Upgrade(args *Args) (*Result, error)
}

// AgentAction is an enum/alias to make calling methods typed
type AgentAction string

const (
	// Install represents the Install method
	Install AgentAction = "Install"

	// Rollback represents the Rollback method
	Rollback AgentAction = "Rollback"

	// Test represents the Test method
	Test AgentAction = "Test"

	// Uninstall represents the Uninstall method
	Uninstall AgentAction = "Uninstall"

	// Upgrade represents the Upgrade method
	Upgrade AgentAction = "Upgrade"
)

// DefaultAgent is the default implementation of the Agent interface
type DefaultAgent struct {
	kube *rest.Config
}

// NewDefaultAgent inits the default agent with the specified kube interface
func NewDefaultAgent(kube *rest.Config) *DefaultAgent {
	return &DefaultAgent{
		kube: kube,
	}
}

// Install is the equivalent of `helm install`
func (a DefaultAgent) Install(args *Args) (*Result, error) {
	cfg, logs, err := helmClient(a.kube, args.Namespace)
	if err != nil {
		return nil, err
	}

	install := action.NewInstall(cfg)
	install.ReleaseName = args.Name
	install.Namespace = args.Namespace
	install.Wait = args.Wait

	logger.Info("chart dependencies", "chart", args.Chart.Dependencies())
	logger.Info("chart vals", "vals", args.Chart.Values)
	result, err := install.Run(args.Chart, args.Values)
	return &Result{
		Data: result,
		Logs: logs.String(),
		Type: ReleaseResult,
	}, err
}

// Rollback is the equivalent of `helm rollback`
func (a DefaultAgent) Rollback(args *Args) (*Result, error) {
	cfg, logs, err := helmClient(a.kube, args.Namespace)
	if err != nil {
		return nil, err
	}

	rollback := action.NewRollback(cfg)
	rollback.Version = args.Version
	rollback.Wait = args.Wait

	err = rollback.Run(args.Name)
	return &Result{
		Logs: logs.String(),
		Type: NoneResult,
	}, err
}

// Test is the equivalent of `helm test`
func (a DefaultAgent) Test(args *Args) (*Result, error) {
	cfg, logs, err := helmClient(a.kube, args.Namespace)
	if err != nil {
		return nil, err
	}

	test := action.NewReleaseTesting(cfg)
	test.Namespace = args.Namespace

	result, err := test.Run(args.Name)
	return &Result{
		Data: result,
		Logs: logs.String(),
		Type: ReleaseResult,
	}, err
}

// Uninstall is the equivalent of `helm uninstall`
func (a DefaultAgent) Uninstall(args *Args) (*Result, error) {
	cfg, logs, err := helmClient(a.kube, args.Namespace)
	if err != nil {
		return nil, err
	}

	uninstall := action.NewUninstall(cfg)
	result, err := uninstall.Run(args.Name)
	return &Result{
		Data: result,
		Logs: logs.String(),
		Type: UninstallReleaseResponseResult,
	}, err
}

// Upgrade is the equivalent of `helm upgrade`
func (a DefaultAgent) Upgrade(args *Args) (*Result, error) {
	cfg, logs, err := helmClient(a.kube, args.Namespace)
	if err != nil {
		return nil, err
	}

	upgrade := action.NewUpgrade(cfg)
	upgrade.Namespace = args.Namespace
	upgrade.Wait = args.Wait

	result, err := upgrade.Run(args.Name, args.Chart, args.Values)
	return &Result{
		Data: result,
		Logs: logs.String(),
		Type: ReleaseResult,
	}, err
}

func helmClient(kube *rest.Config, namespace string) (*action.Configuration, *bytes.Buffer, error) {
	logs := &bytes.Buffer{}
	logger := func(t string, a ...interface{}) {
		str := fmt.Sprintf(t, a...)
		if str[len(str)-1] != '\n' {
			str = str + "\n"
		}
		logs.Write([]byte(str))
	}

	flags := &genericclioptions.ConfigFlags{
		APIServer:   &kube.Host,
		BearerToken: &kube.BearerToken,
		Namespace:   &namespace,
		CAFile:      &kube.CAFile,
	}

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(flags, namespace, "secrets", logger); err != nil {
		return nil, nil, err
	}

	return actionConfig, logs, nil
}
