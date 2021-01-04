package template

// Type is a type alias for the enums for template types
type Type string

const (
	// DEPLOYMENT is the enum for a Deployment template
	DEPLOYMENT Type = "DEPLOYMENT"

	// STEP is the enum for a Step template
	STEP Type = "STEP"

	// TRIGGER is the enum for a Trigger template
	TRIGGER Type = "TRIGGER"
)
