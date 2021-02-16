package helm

import (
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart"
	"helm.sh/helm/v3/pkg/release"
	"k8s.io/cli-runtime/pkg/genericclioptions"
	"k8s.io/client-go/rest"
)

// Agent is the interface we use to talk to helm packages
type Agent interface {
	Install(kube rest.Config, args Args) (*release.Release, error)
	Rollback(kube rest.Config, version int, args Args) error
	Test(kube rest.Config, args Args) (*release.Release, error)
	Uninstall(kube rest.Config, args Args) (*release.UninstallReleaseResponse, error)
	Upgrade(kube rest.Config, args Args) (*release.Release, error)
}

// DefaultAgent is the default implementation of the Agent interface
type DefaultAgent struct {
	kube *rest.Config
}

// Args are arguments common to all commands
type Args struct {
	Name      string
	Namespace string
	Chart     *chart.Chart
	Values    map[string]interface{}
	Wait      bool
	Logger    func(string, ...interface{})
}

// NewDefaultAgent inits the default agent with the specified kube interface
func NewDefaultAgent(kube *rest.Config) *DefaultAgent {
	return &DefaultAgent{
		kube: kube,
	}
}

// Install is the equivalent of `helm install`
func (a DefaultAgent) Install(kube rest.Config, args Args) (*release.Release, error) {
	cfg, err := helmClient(kube, args.Namespace, args.Logger)
	if err != nil {
		return nil, err
	}

	install := action.NewInstall(cfg)
	install.ReleaseName = args.Name
	install.Namespace = args.Namespace
	install.Wait = args.Wait
	return install.Run(args.Chart, args.Values)
}

// Rollback is the equivalent of `helm rollback`
func (a DefaultAgent) Rollback(kube rest.Config, version int, args Args) error {
	cfg, err := helmClient(kube, args.Namespace, args.Logger)
	if err != nil {
		return err
	}

	rollback := action.NewRollback(cfg)
	rollback.Version = version
	rollback.Wait = args.Wait
	return rollback.Run(args.Name)
}

// Test is the equivalent of `helm test`
func (a DefaultAgent) Test(kube rest.Config, args Args) (*release.Release, error) {
	cfg, err := helmClient(kube, args.Namespace, args.Logger)
	if err != nil {
		return nil, err
	}

	test := action.NewReleaseTesting(cfg)
	test.Namespace = args.Namespace
	return test.Run(args.Name)
}

// Uninstall is the equivalent of `helm uninstall`
func (a DefaultAgent) Uninstall(kube rest.Config, args Args) (*release.UninstallReleaseResponse, error) {
	cfg, err := helmClient(kube, args.Namespace, args.Logger)
	if err != nil {
		return nil, err
	}

	uninstall := action.NewUninstall(cfg)
	return uninstall.Run(args.Name)
}

// Upgrade is the equivalent of `helm upgrade`
func (a DefaultAgent) Upgrade(kube rest.Config, args Args) (*release.Release, error) {
	cfg, err := helmClient(kube, args.Namespace, args.Logger)
	if err != nil {
		return nil, err
	}

	upgrade := action.NewUpgrade(cfg)
	upgrade.Namespace = args.Namespace
	upgrade.Wait = args.Wait
	return upgrade.Run(args.Name, args.Chart, args.Values)
}

func helmClient(kube rest.Config, namespace string, logger func(t string, a ...interface{})) (*action.Configuration, error) {
	url, _, err := rest.DefaultServerURL(kube.Host, kube.APIPath, *kube.GroupVersion, kube.Insecure)
	if err != nil {
		return nil, err
	}
	urlString := url.String()

	flags := &genericclioptions.ConfigFlags{
		APIServer:   &urlString,
		BearerToken: &kube.BearerToken,
		Namespace:   &namespace,
		CAFile:      &kube.CAFile,
	}

	actionConfig := new(action.Configuration)
	if err := actionConfig.Init(flags, namespace, "secrets", logger); err != nil {
		return nil, err
	}

	return actionConfig, nil
}
