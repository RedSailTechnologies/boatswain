package gyarados

// DeliverySpec is functionality common to all Delivery Objects
type DeliverySpec interface {
	SubstituteTemplates(*[]*DeliverySpec) error
	SubsituteValues(*map[string]interface{}) error
	Validate() error
	YAML([]byte) error
}

// Template represents a template substitution
type Template struct {
	Template  *string     `yaml:"template,omitempty,flow"`
	Arguments interface{} `yaml:"arguments,omitempty,flow"`
}
