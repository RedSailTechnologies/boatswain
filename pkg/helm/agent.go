package helm

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"

	"k8s.io/cli-runtime/pkg/genericclioptions"

	"github.com/redsailtechnologies/boatswain/pkg/logger"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/chart/loader"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/release"
	"helm.sh/helm/v3/pkg/repo"
)

// Agent is the interface we use to talk to helm packages
type Agent interface {
	CheckIndex(*repo.ChartRepository) bool
	GetChart(name, version, endpoint string) ([]byte, error)
	GetCharts(*repo.ChartRepository) (map[string]repo.ChartVersions, error)
	Install(args Args) (*release.Release, error)
	Rollback(version int, args Args) error
	Test(args Args) (*release.Release, error)
	Uninstall(args Args) (*release.UninstallReleaseResponse, error)
	Upgrade(args Args) (*release.Release, error)
}

// DefaultAgent is the default implementation of the Agent interface
type DefaultAgent struct{}

// Args are arguments common to the install, upgrade, etc. commands
type Args struct {
	Name      string
	Namespace string
	Endpoint  string
	Token     string
	Chart     []byte
	Values    map[string]interface{}
	Wait      bool
	Logger    io.Writer
}

// CheckIndex checks the index.yaml file at the repo's endpoint
func (a DefaultAgent) CheckIndex(r *repo.ChartRepository) bool {
	str, err := r.DownloadIndexFile()
	return str != "" && err == nil
}

// GetChart downloads a single chart from a particular chart repo
func (a DefaultAgent) GetChart(name, version, endpoint string) ([]byte, error) {
	out := os.TempDir()

	pull := action.NewPull()
	pull.ChartPathOptions = action.ChartPathOptions{
		RepoURL: endpoint,
		// InsecureSkipTLSverify: true // TODO AdamP do we need this?
	}
	pull.Settings = cli.New()
	pull.RepoURL = endpoint
	pull.Version = version

	pull.DestDir = out

	_, err := pull.Run(name)
	if err != nil {
		return nil, err
	}

	path := path.Join(out, fmt.Sprintf("%s-%s.tgz", name, version))
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}

	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	// TODO - some better dir management here could be good as if this fails we just eventually overload the filesystem
	if err = os.Remove(path); err != nil {
		logger.Warn("could not remove temporary remove temporary file which could lead to disk overusage", "error", err)
	}

	return bytes, nil
}

// GetCharts gets all charts from a particular chart repo
func (a DefaultAgent) GetCharts(r *repo.ChartRepository) (map[string]repo.ChartVersions, error) {
	str, err := r.DownloadIndexFile()
	if err != nil {
		return nil, err
	}

	idx, err := repo.LoadIndexFile(str)
	if err != nil {
		return nil, err
	}

	return idx.Entries, nil
}

// Install is the equivalent of `helm install`
func (a DefaultAgent) Install(args Args) (*release.Release, error) {
	chart, err := loader.LoadArchive(bytes.NewReader(args.Chart))
	if err != nil {
		return nil, err
	}

	cfg, err := helmClient(args.Endpoint, args.Token, args.Namespace, args.Logger)
	if err != nil {
		return nil, err
	}

	install := action.NewInstall(cfg)
	install.ReleaseName = args.Name
	install.Namespace = args.Namespace
	return install.Run(chart, args.Values)
}

// Rollback is the equivalent of `helm rollback`
func (a DefaultAgent) Rollback(version int, args Args) error {
	cfg, err := helmClient(args.Endpoint, args.Token, args.Namespace, args.Logger)
	if err != nil {
		return err
	}

	rollback := action.NewRollback(cfg)
	rollback.Version = version
	return rollback.Run(args.Name)
}

// Test is the equivalent of `helm test`
func (a DefaultAgent) Test(args Args) (*release.Release, error) {
	cfg, err := helmClient(args.Endpoint, args.Token, args.Namespace, args.Logger)
	if err != nil {
		return nil, err
	}

	test := action.NewReleaseTesting(cfg)
	test.Namespace = args.Namespace
	return test.Run(args.Name)
}

// Uninstall is the equivalent of `helm uninstall`
func (a DefaultAgent) Uninstall(args Args) (*release.UninstallReleaseResponse, error) {
	cfg, err := helmClient(args.Endpoint, args.Token, args.Namespace, args.Logger)
	if err != nil {
		return nil, err
	}

	uninstall := action.NewUninstall(cfg)
	return uninstall.Run(args.Name)
}

// Upgrade is the equivalent of `helm upgrade`
func (a DefaultAgent) Upgrade(args Args) (*release.Release, error) {
	chart, err := loader.LoadArchive(bytes.NewReader(args.Chart))
	if err != nil {
		return nil, err
	}

	cfg, err := helmClient(args.Endpoint, args.Token, args.Namespace, args.Logger)
	if err != nil {
		return nil, err
	}

	upgrade := action.NewUpgrade(cfg)
	upgrade.Namespace = args.Namespace
	return upgrade.Run(args.Namespace, chart, args.Values)
}

func helmClient(endpoint, token, namespace string, logger io.Writer) (*action.Configuration, error) {
	flags := &genericclioptions.ConfigFlags{
		APIServer:   &endpoint,
		BearerToken: &token,
		// TODO AdamP - flags only supports cert files, how do we want to handle?
		// CertFile:    &cluster.Cert,
		Insecure: &[]bool{true}[0],
	}
	actionConfig := new(action.Configuration)
	logs := func(t string, a ...interface{}) {
		str := fmt.Sprintf(t, a...)
		if str[len(str)-1] != '\n' {
			str = str + "\n"
		}
		logger.Write([]byte(str))
	}

	if err := actionConfig.Init(flags, namespace, "secrets", logs); err != nil {
		return nil, err
	}

	return actionConfig, nil
}
