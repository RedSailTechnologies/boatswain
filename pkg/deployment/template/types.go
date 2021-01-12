package template

// A Template is the yaml syntax used to substitute yaml from another file
type Template struct {
	Name      string                  `yaml:"template"`
	Branch    string                  `yaml:"branch"`
	Repo      string                  `yaml:"repo"`
	Arguments *map[string]interface{} `yaml:"arguments,omitempty"`
}

// A Deployment is the go type representing yaml for a cd object
type Deployment struct {
	Version  string     `yaml:"version"`
	Clusters []string   `yaml:"clusters"`
	Apps     *[]App     `yaml:"apps,omitempty"`
	Tests    *[]Test    `yaml:"tests,omitempty"`
	Triggers *[]Trigger `yaml:"triggers,omitempty"`
	Strategy *[]Step    `yaml:"strategy,omitempty"`
}

// Validate is how a deployment can verify structural correctness before execution
func (d Deployment) Validate() error {
	return nil
}

// An App is a yaml representation of a deployable asset
type App struct {
	Name string  `yaml:"name"`
	Helm HelmApp `yaml:"helm"`
}

// A HelmApp is the helm specification for an application
type HelmApp struct {
	Chart   string `yaml:"chart"`
	Repo    string `yaml:"repo"`
	Version string `yaml:"version"`
}

// A Test represents a test to run against a deployment
type Test struct {
	Name string    `yaml:"name"`
	Helm *HelmTest `yaml:"helm,omitempty"`
}

// A HelmTest leverages the helm test for a
type HelmTest struct {
	Name    string `yaml:"name"`
	Timeout string `yaml:"timeout"`
}

// A Trigger is the thing which starts a deployment
type Trigger struct {
	Deployment *DeploymentTrigger `yaml:"deployment,omitempty"`
	Web        *WebTrigger        `yaml:"web,omitempty"`
	Manual     *ManualTrigger     `yaml:"manual,omitempty"`
}

// A DeploymentTrigger is when this deployment can be triggered by another
type DeploymentTrigger struct {
	Name string `yaml:"name"`
}

// A WebTrigger is when an endpoint is called triggering this deployment
type WebTrigger struct {
	Name string `yaml:"name"`
}

// A ManualTrigger is when a person (through the ui, generally) triggers a deployment
type ManualTrigger struct {
	Groups []string `yaml:"groups"`
	Users  []string `yaml:"users"`
}

// A Step is a particular part of a deployment cycle
type Step struct {
	Name    string    `yaml:"name"`
	Hold    string    `yaml:"hold"`
	Success *[]Action `yaml:"success,omitempty"`
	Failure *[]Action `yaml:"failure,omitempty"`
	Any     *[]Action `yaml:"any,omitempty"`
	Always  *[]Action `yaml:"always,omitempty"`
}

// An Action is an action taken within a step
type Action struct {
	App      *AppAction      `yaml:"app,omitempty"`
	Test     *TestAction     `yaml:"test,omitempty"`
	Approval *ApprovalAction `yaml:"approval,omitempty"`
	Trigger  *TriggerAction  `yaml:"trigger,omitempty"`
}

// An AppAction is the action of deploying an application
type AppAction struct {
	Name    string         `yaml:"name"`
	Cluster string         `yaml:"cluster"`
	Helm    *HelmAppAction `yaml:"helm,omitempty"`
}

// A HelmAppAction is the helm specifics for deploying an app
type HelmAppAction struct {
	Command string               `yaml:"command"`
	Wait    bool                 `yaml:"wait"`
	Version int                  `yaml:"version"`
	Values  *HelmAppActionValues `yaml:"values,omitempty"`
}

// HelmAppActionValues are the values for a helm command
type HelmAppActionValues struct {
	Library *HelmValuesLibrary     `yaml:"library,omitempty"`
	Raw     map[string]interface{} `yaml:"raw,omitempty"`
}

// A HelmValuesLibrary are when we get the values from a helm library
type HelmValuesLibrary struct {
	Chart   string   `yaml:"chart"`
	Repo    string   `yaml:"repo"`
	Version string   `yaml:"version"`
	Files   []string `yaml:"files"`
}

// A TestAction is the action which runs tests
type TestAction struct {
	Name    string `yaml:"name"`
	Cluster string `yaml:"cluster"`
}

// An ApprovalAction is when a person must manually approve to unblock further actions
type ApprovalAction struct {
	Groups []string `yaml:"groups"`
	Users  []string `yaml:"users"`
}

// A TriggerAction is where this deployment triggers another
type TriggerAction struct {
	Name       string                 `yaml:"name"`
	Deployment string                 `yaml:"deployment"`
	Arguments  map[string]interface{} `yaml:"arguments,omitempty"`
}
