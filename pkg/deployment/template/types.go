package template

type template struct {
	Ref       string                 `yaml:"template"`
	Repo      string                 `yaml:"repo"`
	Arguments map[string]interface{} `yaml:"arguments"`
	value     map[string]interface{}
}

type deployment struct {
	Version  string     `yaml:"version"`
	Clusters []string   `yaml:"clusters"`
	Apps     *[]app     `yaml:"apps,omitempty"`
	Tests    *[]test    `yaml:"tests,omitempty"`
	Triggers *[]trigger `yaml:"triggers,omitempty"`
	Strategy *[]step    `yaml:"strategy,omitempty"`
}

type app struct {
	Name string  `yaml:"name"`
	Helm helmApp `yaml:"helm"`
}

type helmApp struct {
	Chart   string `yaml:"chart"`
	Repo    string `yaml:"repo"`
	Version string `yaml:"version"`
}

type test struct {
	Name string    `yaml:"name"`
	Helm *helmTest `yaml:"helm,omitempty"`
}

type helmTest struct {
	Params map[string]interface{} `yaml:"params"`
}

type trigger struct {
	Deployment *deploymentTrigger `yaml:"deployment,omitempty"`
	Web        *webTrigger        `yaml:"web,omitempty"`
	Manual     *manualTrigger     `yaml:"manual,omitempty"`
}

type deploymentTrigger struct {
	Name string `yaml:"name"`
}

type webTrigger struct {
	Name string `yaml:"name"`
}

type manualTrigger struct {
	Groups []string `yaml:"groups"`
	Users  []string `yaml:"users"`
}

type step struct {
	Name    string    `yaml:"name"`
	Hold    string    `yaml:"hold"`
	Success *[]action `yaml:"success,omitempty"`
	Failure *[]action `yaml:"failure,omitempty"`
	Any     *[]action `yaml:"any,omitempty"`
	Always  *[]action `yaml:"always,omitempty"`
}

type action struct {
	App      *appAction      `yaml:"app,omitempty"`
	Test     *testAction     `yaml:"test,omitempty"`
	Approval *approvalAction `yaml:"approval,omitempty"`
	Trigger  *triggerAction  `yaml:"trigger,omitempty"`
}

type appAction struct {
	Name    string         `yaml:"name"`
	Cluster string         `yaml:"cluster"`
	Helm    *helmAppAction `yaml:"helm,omitempty"`
}

type helmAppAction struct {
	Command string               `yaml:"command"`
	Wait    bool                 `yaml:"wait"`
	Version int                  `yaml:"version"`
	Values  *helmAppActionValues `yaml:"values,omitempty"`
}

type helmAppActionValues struct {
	Library *helmValuesLibrary     `yaml:"library,omitempty"`
	Raw     map[string]interface{} `yaml:"raw,omitempty"`
}

type helmValuesLibrary struct {
	Chart   string   `yaml:"chart"`
	Repo    string   `yaml:"repo"`
	Version string   `yaml:"version"`
	Files   []string `yaml:"files"`
}

type testAction struct {
	Name    string `yaml:"name"`
	Cluster string `yaml:"cluster"`
}

type approvalAction struct {
	Groups []string `yaml:"groups"`
	Users  []string `yaml:"users"`
}

type triggerAction struct {
	Name       string                 `yaml:"name"`
	Deployment string                 `yaml:"deployment"`
	Arguments  map[string]interface{} `yaml:"arguments,omitempty"`
}
