package repo

// Type represents the kind of repo
type Type int64

const (
	// HELM repo
	HELM Type = 0

	// GIT repo
	GIT Type = 1
)
