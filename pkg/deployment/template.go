package deployment

// A Substitution is the yaml syntax used to substitute yaml from another file
type Substitution struct {
	Name      string                  `yaml:"template"`
	Branch    string                  `yaml:"branch"`
	Repo      string                  `yaml:"repo"`
	Arguments *map[string]interface{} `yaml:"arguments,omitempty"`
}

// A Template is the go type representing yaml for a cd object
type Template struct {
	Version string `yaml:"version"`

	Clusters []string `yaml:"clusters"`

	Apps *[]App `yaml:"apps,omitempty"`

	Tests []struct {
		Name string `yaml:"name"`
		Helm *struct {
			Name    string `yaml:"name"`
			Timeout string `yaml:"timeout"`
		} `yaml:"helm,omitempty"`
	} `yaml:"tests,omitempty"`

	Triggers *[]struct {
		Deployment *struct {
			Name string `yaml:"name"`
		} `yaml:"deployment,omitempty"`
		Web *struct {
			Name string `yaml:"name"`
		} `yaml:"web,omitempty"`
		Manual *struct {
			Groups []string `yaml:"groups"`
			Users  []string `yaml:"users"`
		} `yaml:"manual,omitempty"`
	} `yaml:"triggers,omitempty"`

	Strategy *[]Step `yaml:"strategy,omitempty"`
}

// An App is a deployable application within a template
type App struct {
	Name string `yaml:"name"`
	Helm *struct {
		Chart   string `yaml:"chart"`
		Repo    string `yaml:"repo"`
		Version string `yaml:"version"`
	} `yaml:"helm,omitempty"`
}

// A Step is an individual step within a Strategy
type Step struct {
	// metadata/execution information
	Name      string `yaml:"name"`
	Hold      string `yaml:"hold"`
	Condition string `yaml:"condition"`

	// internal fields
	Status Status `yaml:"-"`
	Logs   []log  `yaml:"-"`
	Start  int64  `yaml:"-"`

	// step information
	App *struct {
		Name      string `yaml:"name"`
		Cluster   string `yaml:"cluster"`
		Namespace string `yaml:"namespace"`
		Helm      *struct {
			Command string `yaml:"command"`
			Wait    bool   `yaml:"wait"`
			Version int    `yaml:"version"`
			Values  *struct {
				Library *struct {
					Chart   string   `yaml:"chart"`
					Repo    string   `yaml:"repo"`
					Version string   `yaml:"version"`
					Files   []string `yaml:"files"`
				} `yaml:"library,omitempty"`
				Raw *map[string]interface{} `yaml:"raw,omitempty"`
			} `yaml:"values,omitempty"`
		} `yaml:"helm,omitempty"`
	} `yaml:"app,omitempty"`

	Test *struct {
		Name      string `yaml:"name"`
		Cluster   string `yaml:"cluster"`
		Namespace string `yaml:"namespace"`
	} `yaml:"test,omitempty"`

	Approval *struct {
		Groups []string `yaml:"groups"`
		Users  []string `yaml:"users"`
	} `yaml:"approval,omitempty"`

	Trigger *struct {
		Name       string                 `yaml:"name"`
		Deployment string                 `yaml:"deployment"`
		Arguments  map[string]interface{} `yaml:"arguments,omitempty"`
	} `yaml:"trigger,omitempty"`
}

// Validate is how a deployment can verify structural correctness before execution
func (t Template) Validate() error {
	return nil
}
