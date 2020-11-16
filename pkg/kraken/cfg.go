package kraken

// ClusterConfig is the configuration struct for a single cluster
type ClusterConfig struct {
	Name     string `yaml:"name"`
	Endpoint string `yaml:"endpoint"`
	Token    string `yaml:"token"`
	Cert     string `yaml:"cert"`
}

// Config is the service's configuration
type Config struct {
	Clusters []ClusterConfig `yaml:"clusters"`
}
