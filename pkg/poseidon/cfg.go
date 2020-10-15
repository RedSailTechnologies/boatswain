package poseidon

import (
	"errors"
	"io/ioutil"

	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/getter"
	"helm.sh/helm/v3/pkg/repo"

	"gopkg.in/yaml.v2"
)

// RepoConfig is the configuration struct for a helm repo
type RepoConfig struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
}

// ToChartPathOptions returns the chart path options for helm
func (c *RepoConfig) ToChartPathOptions() *action.ChartPathOptions {
	return &action.ChartPathOptions{
		InsecureSkipTLSverify: true,
		RepoURL:               c.Endpoint,
	}
}

// ToChartRepo takes the configuration and makes it into a working repo
func (c *RepoConfig) ToChartRepo() (*repo.ChartRepository, error) {
	providers := []getter.Provider{
		getter.Provider{
			Schemes: []string{"http", "https"},
			New:     getter.NewHTTPGetter,
		},
	}

	entry := &repo.Entry{
		Name: c.Name,
		URL:  c.Endpoint,
		// TODO AdamP - we definitely want to support this soon!
		InsecureSkipTLSverify: true,
	}

	return repo.NewChartRepository(entry, providers)
}

// Config is a list of configurations
type Config struct {
	Repos    []RepoConfig `yaml:"repos"`
	CacheDir string       `yaml:"cacheDir"`
}

// YAML takes a relative filename and returns the config found in it
func (c *Config) YAML(file string) error {
	y, err := ioutil.ReadFile(file)
	if err != nil {
		return err
	}

	// TODO AdamP - need to get keyed values from a file to pass here so we can revert to strict
	if err := yaml.Unmarshal(y, c); err != nil {
		return err
	}
	return nil
}

func (c *Config) getRepoConfig(repoName string) (*RepoConfig, error) {
	for _, config := range c.Repos {
		if config.Name == repoName {
			return &config, nil
		}
	}
	return nil, errors.New("repo not found")
}
