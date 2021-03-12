package template

// A Substitution is the yaml syntax used to substitute yaml from another file
type Substitution struct {
	Name      string                  `yaml:"template"`
	Branch    string                  `yaml:"branch"`
	Repo      string                  `yaml:"repo"`
	Arguments *map[string]interface{} `yaml:"arguments,omitempty"`
}

// A Template is the go type representing yaml for a cd object
type Template struct {
	Name        string       `yaml:"name"`
	Version     string       `yaml:"version"`
	Application *Application `yaml:"apps,omitempty"`
	Links       []struct {
		Name string `yaml:"name"`
		URL  string `yaml:"url"`
	} `yaml:"links"`
	Triggers *[]Trigger `yaml:"triggers,omitempty"`
	Strategy *[]Step    `yaml:"strategy,omitempty"`
}

// An Application is a way to link deployments to applications
type Application struct {
	Name    string `yaml:"name"`
	PartOf  string `yaml:"partOf"`
	Version string `yaml:"version"`
}

// A Trigger is a way of validating how this deployment is run
type Trigger struct {
	Deployment *struct {
		Name string `yaml:"name"`
	} `yaml:"deployment,omitempty"`
	Web *struct {
		Name string `yaml:"name"`
	} `yaml:"web,omitempty"`
	Manual *struct {
		Name  string   `yaml:"name"`
		Role  string   `yaml:"role"`
		Users []string `yaml:"users"`
	} `yaml:"manual,omitempty"`
}

// A Step is an individual step within a Strategy
type Step struct {
	// metadata/execution information
	Name      string `yaml:"name"`
	Hold      string `yaml:"hold"`
	Condition string `yaml:"condition"`

	// step information
	Helm *struct {
		Name      string             `yaml:"name"`
		Selector  *map[string]string `yaml:"selector,omitempty"`
		Cluster   string             `yaml:"cluster"`
		Namespace string             `yaml:"namespace"`
		Command   string             `yaml:"command"`
		Chart     *struct {
			Name    string `yaml:"name"`
			Repo    string `yaml:"repo"`
			Version string `yaml:"version"`
		} `yaml:"chart,omitempty"`
		Options *struct {
			RollbackVersion int    `yaml:"rollbackVersion"`
			Wait            bool   `yaml:"wait"`
			Install         bool   `yaml:"install"`
			Timeout         string `yaml:"timeout"`
			ReuseValues     bool   `yaml:"reuseValues"`
		} `yaml:"options"`
		Values *struct {
			Library *struct {
				Chart   string `yaml:"chart"`
				Repo    string `yaml:"repo"`
				Version string `yaml:"version"`
				File    string `yaml:"file"`
			} `yaml:"library,omitempty"`
			Raw *map[string]interface{} `yaml:"raw,omitempty"`
		} `yaml:"values,omitempty"`
	} `yaml:"helm,omitempty"`

	Approval *struct {
		Name   string   `yaml:"name"`
		Action string   `yaml:"action"`
		Roles  []string `yaml:"roles"`
		Users  []string `yaml:"users"`
	} `yaml:"approval,omitempty"`

	Trigger *struct {
		Name       string                 `yaml:"name"`
		Deployment string                 `yaml:"deployment"`
		Arguments  map[string]interface{} `yaml:"arguments,omitempty"`
	} `yaml:"trigger,omitempty"`
}
