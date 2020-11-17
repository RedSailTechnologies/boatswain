package poseidon

// RepoConfig is the configuration struct for a helm repo
type RepoConfig struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
}

// Config is a list of configurations
type Config struct {
	Repos    []RepoConfig `yaml:"repos"`
	CacheDir string       `yaml:"cacheDir"`
}
